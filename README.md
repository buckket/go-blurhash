# go-blurhash

**go-blurhash** is a pure Go implementation of the Blurhash algorithm, which is used by
[Mastodon](https://github.com/tootsuite/mastodon) an other Fediverse software to implement a swift way of preloading images as well
as hiding sensitive media. Read more about it [here](https://blog.joinmastodon.org/2019/05/improving-support-for-adult-content-on-mastodon/).

This library allows generating the Blurhash of a given image, as well as
reconstructing a blurred version with specified dimensions from a given Blurhash.

This library is based entirely on the current reference implementations<sup>[1](#note1)</sup>:
- Encoder: https://github.com/Gargron/blurhash (C, Ruby)
- Deocder: https://github.com/Gargron/blurhash.js (TypeScript)

Blurhash is written by [Dag Ågren](https://github.com/DagAgren).

<a name="note">1</a>: Because there is no real spec as of yet.

## Installation

### From source

    go get -u git.buckket.org/buckket/go-blurhash



## Usage

go-blurhash exports three functions:
- blurhash.Encode()
- blurhash.Decodde()
- blurhash.Components()

Here’s a simple demonstration, till the documentation catches up:


```go
package main

import (
	"fmt"
	"github.com/buckket/go-blurhash/blurhash"
	"image/png"
	"os"
)

func main() {
	// Encode an image
	imageFile, _ := os.Open("test.png")
	loadedImage, err := png.Decode(imageFile)
	str, _ := blurhash.Encode(4, 3, &loadedImage)
	if err != nil {
		// Handling errors
	}
	fmt.Printf("Hash: %s", str)

	// Generating an image
	img, err := blurhash.Decode(str, 300, 500, 1)
	if err != nil {
		// Handling errors
	}
	f, _ := os.Create("test_blur.png")
	_ = png.Encode(f, img)
}

```

## Limitations

- Documentation lacking (here, as well as upstream)
- Automated tests WIP

## License

 GNU GPLv3+
 