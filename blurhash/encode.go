package blurhash

import (
	"fmt"
	"github.com/buckket/go-blurhash/base83"
	"image"
	"math"
)

func Encode(xComponents int, yComponents int, rgba *image.Image) (string, error) {
	if xComponents < 1 || xComponents > 9 {
		return "", fmt.Errorf("blurhash: xComponents out of valid range (1-9)")
	}
	if yComponents < 1 || yComponents > 9 {
		return "", fmt.Errorf("blurhash: yComponents out of valid range (1-9)")
	}

	pos := 0
	buffer := make([]byte, 2+4+(9*9-1)*2+1) // 167

	factors := make([] float64, yComponents*xComponents*3)
	for y := 0; y < yComponents; y++ {
		for x := 0; x < xComponents; x++ {
			factor := multiplyBasisFunction(x, y, rgba)
			factors[0+x*3+y*3*xComponents] = factor[0]
			factors[1+x*3+y*3*xComponents] = factor[1]
			factors[2+x*3+y*3*xComponents] = factor[2]
		}
	}

	acCount := xComponents*yComponents - 1
	sizeFlag := (xComponents - 1) + (yComponents-1)*9

	pos = base83.Encode(sizeFlag, 1, pos, &buffer)

	var maximumValue float64
	if acCount > 0 {
		var actualMaximumValue float64
		for i := 0; i < acCount*3; i++ {
			actualMaximumValue = math.Max(math.Abs(factors[i+3]), actualMaximumValue)
		}
		quantisedMaximumValue := int(math.Max(0, math.Min(82, math.Floor(actualMaximumValue*166-0.5))))
		maximumValue = (float64(quantisedMaximumValue) + 1) / 166
		pos = base83.Encode(quantisedMaximumValue, 1, pos, &buffer)

	} else {
		maximumValue = 1
		pos = base83.Encode(0, 1, pos, &buffer)
	}

	pos = base83.Encode(int(encodeDC(factors[0], factors[1], factors[2])), 4, pos, &buffer)

	for i := 0; i < acCount; i++ {
		pos = base83.Encode(encodeAC(factors[3+(i*3+0)], factors[3+(i*3+1)], factors[3+(i*3+2)], maximumValue), 2, pos, &buffer)
	}

	return string(buffer[:pos]), nil
}

func multiplyBasisFunction(xComponent int, yComponent int, rgba *image.Image) [3]float64 {
	var r, g, b float64

	height := (*rgba).Bounds().Max.Y
	width := (*rgba).Bounds().Max.X

	var normalisation float64
	if xComponent == 0 && yComponent == 0 {
		normalisation = 1
	} else {
		normalisation = 2
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			basis := math.Cos(math.Pi*float64(xComponent)*float64(x)/float64(width)) * math.Cos(math.Pi*float64(yComponent)*float64(y)/float64(height))
			rt, gt, bt, _ := (*rgba).At(x, y).RGBA()
			r += basis * sRGBToLinear(int(rt>>8))
			g += basis * sRGBToLinear(int(gt>>8))
			b += basis * sRGBToLinear(int(bt>>8))
		}
	}

	scale := normalisation / float64(width*height)
	return [3]float64{r * scale, g * scale, b * scale}
}

func encodeDC(r, g, b float64) int {
	return (linearTosRGB(r) << 16) + (linearTosRGB(g) << 8) + linearTosRGB(b)
}

func encodeAC(r, g, b, maximumValue float64) int {
	quant := func(f float64) float64 {
		return math.Max(0, math.Min(18, math.Floor(signPow(f/maximumValue, 0.5)*9+9.5)))
	}
	return int(quant(r)*19*19 + quant(g)*19 + quant(b))
}
