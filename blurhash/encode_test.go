package blurhash_test

import (
	"fmt"
	"github.com/buckket/go-blurhash/blurhash"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"
)

func TestEncodeFile(t *testing.T) {
	imageFile, err := os.Open("../test.png")
	if err != nil {
		t.Fatal("could not load test.png")
	}

	loadedImage, err := png.Decode(imageFile)
	if err != nil {
		t.Fatal("could not decode test.png")
	}

	str, err := blurhash.Encode(4, 3, &loadedImage)
	if err != nil {
		t.Fatal(err)
		return
	}

	const expectedHash = "LFE.@D9F01_2%L%MIVD*9Goe-;WB"
	if str != expectedHash {
		t.Fatalf("got %q, expected %q", str, expectedHash)
	}
}

func TestEncodeWrongParameters(t *testing.T) {
	testCases := []struct {
		xComp int
		yComp int
	}{
		{0, 1},
		{1, 0},
		{0, 0},
		{10, 1},
		{1, 10},
		{10, 10},
		{-1, -1},
	}

	var img image.Image
	img = image.NewNRGBA(image.Rect(0, 0, 100, 100))

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("xComp:%d yComp:%d", tc.xComp, tc.yComp), func(t *testing.T) {
			_, err := blurhash.Encode(tc.xComp, tc.yComp, &img)
			if err == nil {
				t.Fatal("should have failed")
			}
			err, ok := err.(blurhash.InvalidParameterError)
			if !ok {
				t.Fatal("wrong error type")
			}
			_ = err.Error()
		})
	}
}

func TestEncodeEmptyImage(t *testing.T) {
	var img image.Image
	img = image.NewNRGBA(image.Rect(0, 0, 100, 100))

	str, err := blurhash.Encode(4, 3, &img)
	if err != nil {
		t.Fatal(err)
	}
	const expectedHash = "L00000fQfQfQfQfQfQfQfQfQfQfQ"
	if str != expectedHash {
		t.Fatalf("got %q, expected %q", str, expectedHash)
	}
}

func TestEncodeSizeFlag(t *testing.T) {
	testCases := []struct {
		xComp int
		yComp int
	}{
		{1, 2},
		{9, 8},
		{5, 4},
		{2, 3},
		{4, 5},
		{7, 3},
	}

	var img image.Image
	img = image.NewNRGBA(image.Rect(0, 0, 100, 100))

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("xComp:%d yComp:%d", tc.xComp, tc.yComp), func(t *testing.T) {
			str, err := blurhash.Encode(tc.xComp, tc.yComp, &img)
			if err != nil {
				t.Fatal(err)
			}
			xComp, yComp, err := blurhash.Components(str)
			if err != nil {
				t.Fatal(err)
			}
			if xComp != tc.xComp || yComp != tc.yComp {
				t.Fatal("component mismatch")
			}
		})
	}
}

func TestEncodeSingleColor(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	tcolor := color.RGBA{213, 30, 120, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{tcolor}, image.ZP, draw.Src)
	var img2 image.Image = img

	str, err := blurhash.Encode(1, 1, &img2)
	if err != nil {
		t.Fatal(err)
	}
	const expectedHash = "00OZZy"
	if str != expectedHash {
		t.Fatalf("got %q, expected %q", str, expectedHash)
	}
}

func BenchmarkEncode(b *testing.B) {
	imageFile, err := os.Open("../test.png")
	if err != nil {
		b.Fatal("could not load test.png")
	}

	loadedImage, err := png.Decode(imageFile)
	if err != nil {
		b.Fatal("could not decode test.png")
	}

	for i := 0; i < b.N; i++ {
		_, _ = blurhash.Encode(4, 3, &loadedImage)
	}
}
