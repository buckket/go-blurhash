# go-blurhash [![Build Status](https://travis-ci.org/buckket/go-blurhash.svg)](https://travis-ci.org/buckket/go-blurhash) [![Go Report Card](https://goreportcard.com/badge/github.com/buckket/go-blurhash)](https://goreportcard.com/report/github.com/buckket/go-blurhash) [![codecov](https://codecov.io/gh/buckket/go-blurhash/branch/master/graph/badge.svg)](https://codecov.io/gh/buckket/go-blurhash) [![GoDoc](https://godoc.org/github.com/buckket/go-blurhash?status.svg)](https://godoc.org/github.com/buckket/go-blurhash)

**go-blurhash** is a pure Go implementation of the Blurhash algorithm, which is used by
[Mastodon](https://github.com/tootsuite/mastodon) an other Fediverse software to implement a swift way of preloading images as well
as hiding sensitive media. Read more about it [here](https://blog.joinmastodon.org/2019/05/improving-support-for-adult-content-on-mastodon/).

This library allows generating the Blurhash of a given image, as well as
reconstructing a blurred version with specified dimensions from a given Blurhash.

This library is based entirely on the current reference implementations:
- Encoder: https://github.com/Gargron/blurhash (C, Ruby)
- Deocder: https://github.com/Gargron/blurhash.js (TypeScript)

Blurhash is written by [Dag Ågren](https://github.com/DagAgren).

|        | Before                         | After                          |
| ------ |:------------------------------:| :-----------------------------:|
| Image  | ![alt text][test]              | "LFE.@D9F01_2%L%MIVD*9Goe-;WB" |
| Hash   | "LFE.@D9F01_2%L%MIVD*9Goe-;WB" | ![alt text][test_blur]

[test]: test.png "Blurhash example input."
[test_blur]: test_blur.png "Blurhash example output"

## Installation

### From source

    go get -u github.com/buckket/go-blurhash

## Usage

go-blurhash exports three functions:
- blurhash.Encode(xComponents int, yComponents int, rgba *image.Image) (string, error)
- blurhash.Decode(hash string, width, height, punch int) (image.Image, error)
- blurhash.Components(hash string) (xComp, yComp int, err error)

Here’s a simple demonstration. Check [GoDoc](https://godoc.org/github.com/buckket/go-blurhash) for the full documentation.

```go
package main

import (
	"fmt"
	"github.com/buckket/go-blurhash"
	"image/png"
	"os"
)

func main() {
	// Generate the Blurhash for a given image
	imageFile, _ := os.Open("test.png")
	loadedImage, err := png.Decode(imageFile)
	str, _ := blurhash.Encode(4, 3, &loadedImage)
	if err != nil {
		// Handle errors
	}
	fmt.Printf("Hash: %s\n", str)

	// Generate an image for a given Blurhash
	// Width will be 300px and Height will be 500px
	// Punch specifies the contrasts and defaults to 1
	img, err := blurhash.Decode(str, 300, 500, 1)
	if err != nil {
		// Handle errors
	}
	f, _ := os.Create("test_blur.png")
	_ = png.Encode(f, img)
	
	// Get the x and y components used for encoding a given Blurhash
	x, y, err := blurhash.Components("LFE.@D9F01_2%L%MIVD*9Goe-;WB")
	if err != nil {
		// Handle errors
	}
	fmt.Printf("xComp: %d, yComp: %d", x, y)
}

```

## Limitations

- Presumably a bit slower than the C implementation

## License

 GNU GPLv3+
 