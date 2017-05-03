# orientation [ ![Codeship Status for minodisk/orientation](https://img.shields.io/codeship/489273a0-105b-0135-b1ec-2e1c9a6cac85/master.svg?style=flat)](https://app.codeship.com/projects/216267) [![Go Report Card](https://goreportcard.com/badge/github.com/minodisk/orientation)](https://goreportcard.com/report/github.com/minodisk/orientation) [![codecov](https://codecov.io/gh/minodisk/orientation/branch/master/graph/badge.svg)](https://codecov.io/gh/minodisk/orientation) [![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat)](https://godoc.org/github.com/minodisk/orientation) [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

Apply EXIF Orientation tag to the pixels of a image.

## Installation

```sh
go get github.com/minodisk/orientation
```

## Usage

### When processing various formats (Recommended):

```go
import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/minodisk/orientation"
)

func main() {
	file, _ := os.Open("path/to/image")
	img, err := orientation.Apply(file)
	if err != nil {
		if err, ok := err.(*orientation.DecodeError); ok {
			panic(err)
		}
	}
	fmt.Println(visualize(img))
}
```

### When processing any JPEG:

```go
import (
	_ "image/jpeg"

	"github.com/minodisk/orientation"
)

func main() {
	file, _ := os.Open("path/to/jpeg")
	img, err := orientation.Apply(file)
	if err != nil {
		if err, ok := err.(*orientation.DecodeError); ok {
			panic(err)
		}
	}
	fmt.Println(visualize(img))
}
```

### When processing only JPEG with correct Orientation tag:

```go
import (
	_ "image/jpeg"

	"github.com/minodisk/orientation"
)

func main() {
	file, _ := os.Open("path/to/jpeg")
	img, err := orientation.Apply(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(visualize(img))
}
```
