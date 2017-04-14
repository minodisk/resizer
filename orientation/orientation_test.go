package orientation_test

import (
	"fmt"
	"image"
	"image/color"
	"reflect"
	"testing"

	"github.com/go-microservices/resizer/orientation"
)

type Pixels struct {
	size   image.Point
	matrix []int
}

var (
	black  = color.RGBA{0, 0, 0, 0xff}
	white  = color.RGBA{0xff, 0xff, 0xff, 0xff}
	colors = map[int]color.RGBA{
		0: black,
		1: white,
	}
)

func NewImage(pixels Pixels) image.Image {
	w := pixels.size.X
	b := image.Rect(0, 0, pixels.size.X, pixels.size.Y)
	i := image.NewRGBA(b)
	// log.Println("-------", b.Dx(), b.Dy())
	for y := b.Min.Y; y < b.Max.Y; y++ {
		dy := y - b.Min.Y
		for x := b.Min.X; x < b.Max.X; x++ {
			dx := x - b.Min.X
			pos := w*dy + dx
			i.SetRGBA(x, y, colors[pixels.matrix[pos]])
			// log.Println(x, y, pos, colors[pixels.matrix[pos]], i.At(x, y))
		}
	}
	return i
}

func Visualize(i image.Image) string {
	// log.Println("=======")
	b := i.Bounds()
	s := ""
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := i.At(x, y)
			// log.Println(x, y, c, isBlack(c), isWhite(c))
			if isBlack(c) {
				s += "□"
				continue
			}
			if isWhite(c) {
				s += "■"
				continue
			}
			s += "○"
		}
		s += "\n"
	}
	s += fmt.Sprintf("Origin: %v, Size: %s", b.Min, b.Size())
	return s
}

func isBlack(c color.Color) bool {
	r, g, b, a := c.RGBA()
	br, bg, bb, ba := black.RGBA()
	return r == br && g == bg && b == bb && a == ba
}

func isWhite(c color.Color) bool {
	r, g, b, a := c.RGBA()
	wr, wg, wb, wa := white.RGBA()
	return r == wr && g == wg && b == wb && a == wa
}

func TestOrient1to8(t *testing.T) {
	type Case struct {
		spec   string
		orient func(image.Image) image.Image
		input  Pixels
	}

	expected := NewImage(Pixels{
		size: image.Pt(3, 5),
		matrix: []int{
			1, 1, 1,
			1, 0, 0,
			1, 1, 0,
			1, 0, 0,
			1, 0, 0,
		},
	})

	cases := []Case{
		{
			spec:   "Can rotate correctly with orientation tag 1",
			orient: orientation.Orient1,
			input: Pixels{
				size: image.Pt(3, 5),
				matrix: []int{
					1, 1, 1,
					1, 0, 0,
					1, 1, 0,
					1, 0, 0,
					1, 0, 0,
				},
			},
		},
		{
			spec:   "Can rotate correctly with orientation tag 2",
			orient: orientation.Orient2,
			input: Pixels{
				size: image.Pt(3, 5),
				matrix: []int{
					1, 1, 1,
					0, 0, 1,
					0, 1, 1,
					0, 0, 1,
					0, 0, 1,
				},
			},
		},
		{
			spec:   "Can rotate correctly with orientation tag 3",
			orient: orientation.Orient3,
			input: Pixels{
				size: image.Pt(3, 5),
				matrix: []int{
					0, 0, 1,
					0, 0, 1,
					0, 1, 1,
					0, 0, 1,
					1, 1, 1,
				},
			},
		},
		{
			spec:   "Can rotate correctly with orientation tag 4",
			orient: orientation.Orient4,
			input: Pixels{
				size: image.Pt(3, 5),
				matrix: []int{
					1, 0, 0,
					1, 0, 0,
					1, 1, 0,
					1, 0, 0,
					1, 1, 1,
				},
			},
		},
		{
			spec:   "Can rotate correctly with orientation tag 5",
			orient: orientation.Orient5,
			input: Pixels{
				size: image.Pt(5, 3),
				matrix: []int{
					1, 1, 1, 1, 1,
					1, 0, 1, 0, 0,
					1, 0, 0, 0, 0,
				},
			},
		},
		{
			spec:   "Can rotate correctly with orientation tag 6",
			orient: orientation.Orient6,
			input: Pixels{
				size: image.Pt(5, 3),
				matrix: []int{
					1, 0, 0, 0, 0,
					1, 0, 1, 0, 0,
					1, 1, 1, 1, 1,
				},
			},
		},
		{
			spec:   "Can rotate correctly with orientation tag 7",
			orient: orientation.Orient7,
			input: Pixels{
				size: image.Pt(5, 3),
				matrix: []int{
					0, 0, 0, 0, 1,
					0, 0, 1, 0, 1,
					1, 1, 1, 1, 1,
				},
			},
		},
		{
			spec:   "Can rotate correctly with orientation tag 8",
			orient: orientation.Orient8,
			input: Pixels{
				size: image.Pt(5, 3),
				matrix: []int{
					1, 1, 1, 1, 1,
					0, 0, 1, 0, 1,
					0, 0, 0, 0, 1,
				},
			},
		},
	}

	for _, c := range cases {
		actual := c.orient(NewImage(c.input))
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("%s. Rotated image is expected:\n%s\nbut actual:\n%s\n", c.spec, Visualize(expected), Visualize(actual))
		}
	}
}
