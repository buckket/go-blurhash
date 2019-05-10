package blurhash

import (
	"fmt"
	"github.com/buckket/go-blurhash/base83"
	"image"
	"math"
	"strings"
)

// An InvalidParameterError occurs when an invalid argument is passed to either the Decode or Encode function.
type InvalidParameterError struct {
	value     int
	parameter string
}

func (e InvalidParameterError) Error() string {
	return fmt.Sprintf("blurhash: %sComponents (%d) must be element of [1-9]", e.parameter, e.value)
}

// An EncodingError represents an error that occurred during the encoding of the given value.
// This most likely means that your input image is invalid and can not be processed.
type EncodingError string

func (e EncodingError) Error() string {
	return fmt.Sprintf("blurhash: %s", string(e))
}

// Encode calculates the Blurhash for an image using the given x and y component counts.
// The x and y components have to be between 1 and 9 respectively.
// The image must be of image.Image type.
func Encode(xComponents int, yComponents int, rgba *image.Image) (string, error) {
	if xComponents < 1 || xComponents > 9 {
		return "", InvalidParameterError{xComponents, "x"}
	}
	if yComponents < 1 || yComponents > 9 {
		return "", InvalidParameterError{yComponents, "y"}
	}

	var blurhash strings.Builder
	blurhash.Grow(4 + 2*xComponents*yComponents)

	// Size Flag
	str, err := base83.Encode((xComponents-1)+(yComponents-1)*9, 1)
	if err != nil {
		return "", EncodingError("could not encode size flag")
	}
	blurhash.WriteString(str)

	factors := make([] float64, yComponents*xComponents*3)
	for y := 0; y < yComponents; y++ {
		for x := 0; x < xComponents; x++ {
			factor := multiplyBasisFunction(x, y, rgba)
			factors[0+x*3+y*3*xComponents] = factor[0]
			factors[1+x*3+y*3*xComponents] = factor[1]
			factors[2+x*3+y*3*xComponents] = factor[2]
		}
	}

	var maximumValue float64
	var quantisedMaximumValue int
	var acCount = xComponents*yComponents - 1
	if acCount > 0 {
		var actualMaximumValue float64
		for i := 0; i < acCount*3; i++ {
			actualMaximumValue = math.Max(math.Abs(factors[i+3]), actualMaximumValue)
		}
		quantisedMaximumValue = int(math.Max(0, math.Min(82, math.Floor(actualMaximumValue*166-0.5))))
		maximumValue = (float64(quantisedMaximumValue) + 1) / 166
	} else {
		maximumValue = 1
	}

	// Quantised max AC component
	str, err = base83.Encode(quantisedMaximumValue, 1)
	if err != nil {
		return "", EncodingError("could not encode quantised max AC component")
	}
	blurhash.WriteString(str)

	// DC value
	str, err = base83.Encode(encodeDC(factors[0], factors[1], factors[2]), 4)
	if err != nil {
		return "", EncodingError("could not encode DC value")
	}
	blurhash.WriteString(str)

	// AC values
	for i := 0; i < acCount; i++ {
		str, err = base83.Encode(encodeAC(factors[3+(i*3+0)], factors[3+(i*3+1)], factors[3+(i*3+2)], maximumValue), 2)
		if err != nil {
			return "", EncodingError("could not encode AC value")
		}
		blurhash.WriteString(str)
	}

	if blurhash.Len() != 4+2*xComponents*yComponents {
		return "", EncodingError("hash does not match expected size")
	}

	return blurhash.String(), nil
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
			basis := math.Cos(math.Pi*float64(xComponent)*float64(x)/float64(width)) *
				math.Cos(math.Pi*float64(yComponent)*float64(y)/float64(height))
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
	quant := func(f float64) int {
		return int(math.Max(0, math.Min(18, math.Floor(signPow(f/maximumValue, 0.5)*9+9.5))))
	}
	return quant(r)*19*19 + quant(g)*19 + quant(b)
}
