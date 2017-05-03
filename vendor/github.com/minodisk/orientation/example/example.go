package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/minodisk/orientation"
)

func main() {
	file, _ := os.Open("../fixtures/f-orientation-8.jpg")
	img, err := orientation.Apply(file)
	if err != nil {
		if err, ok := err.(*orientation.DecodeError); ok {
			panic(err)
		}
	}
	fmt.Println(visualize(img))
}

func visualize(img image.Image) string {
	r := img.Bounds()
	var buf bytes.Buffer
	for y := 0; y < r.Dy(); y++ {
		for x := 0; x < r.Dx(); x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			if ((r >> 15 & 0x1) == 1) && ((g >> 15 & 0x1) == 1) && ((b >> 15 & 0x1) == 1) && ((a >> 15 & 0x1) == 1) {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}
		}
		buf.WriteString("\n")
	}
	return buf.String()
}
