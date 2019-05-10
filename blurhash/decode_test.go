package blurhash_test

import (
	"github.com/buckket/go-blurhash/blurhash"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestComponents(t *testing.T) {
	const str = "LFE.@D9F01_2%L%MIVD*9Goe-;WB"

	x, y, err := blurhash.Components(str)
	if err != nil {
		t.Fatal(err)
	}
	if x != 4 || y != 3 {
		t.Fatalf("component missmatch")
	}

	_, _, err = blurhash.Components("12345")
	if err == nil {
		t.Fatal("should have failed")
	}
	_, ok := err.(blurhash.InvalidHash)
	if !ok {
		t.Fatal("wrong error type")
	}

	_, _, err = blurhash.Components(str[:9])
	if err == nil {
		t.Fatal("should have failed")
	}
	_, ok = err.(blurhash.InvalidHash)
	if !ok {
		t.Fatal("wrong error type")
	}
}

func TestDecodeFile(t *testing.T) {
	const str = "LFE.@D9F01_2%L%MIVD*9Goe-;WB"

	imageFile, err := os.Open("../test_blur.png")
	if err != nil {
		t.Fatal("could not load test_blur.png")
	}

	loadedImage, err := png.Decode(imageFile)
	if err != nil {
		t.Fatal("could not decode test_blur.png")
	}

	img, err := blurhash.Decode(str, 300, 300, 1)
	if err != nil {
		t.Fatal(err)
	}

	if loadedImage.Bounds() != img.Bounds() {
		t.Fatal("bounds mismatch")
	}

	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r1, g1, b1, a1 := img.At(x, y).RGBA()
			r2, g2, b2, a2 := loadedImage.At(x, y).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				t.Fatalf("pixel mismatch")
			}
		}
	}

}

func TestDecodeSingleColor(t *testing.T) {
	const str = "00OZZy"

	img, err := blurhash.Decode(str, 1, 1, 1)
	if err != nil {
		t.Fatal(err)
	}

	r, g, b, a := img.At(0, 0).RGBA()
	bcolor := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	tcolor := color.RGBA{213, 30, 120, 255}
	if bcolor != tcolor {
		t.Fatal("color mismatch")
	}
}

func BenchmarkDecode(b *testing.B) {
	const str = "LFE.@D9F01_2%L%MIVD*9Goe-;WB"

	for i := 0; i < b.N; i++ {
		_, _ = blurhash.Decode(str, 300, 300, 1)
	}
}
