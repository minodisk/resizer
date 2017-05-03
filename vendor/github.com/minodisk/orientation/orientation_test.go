package orientation_test

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/minodisk/orientation"
)

var (
	black  = color.RGBA{0, 0, 0, 0xff}
	white  = color.RGBA{0xff, 0xff, 0xff, 0xff}
	colors = map[int]color.RGBA{
		0: black,
		1: white,
	}
)

func NewImage(size image.Point, matrix []int) image.Image {
	w := size.X
	b := image.Rect(0, 0, size.X, size.Y)
	i := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		dy := y - b.Min.Y
		for x := b.Min.X; x < b.Max.X; x++ {
			dx := x - b.Min.X
			pos := w*dy + dx
			i.SetRGBA(x, y, colors[matrix[pos]])
		}
	}
	return i
}

func Visualize(i image.Image) string {
	b := i.Bounds()
	s := ""
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := i.At(x, y)
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

func TestApply(t *testing.T) {
	t.Parallel()
	want := []int{
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		0, 1, 0, 0, 0,
		0, 1, 1, 1, 0,
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 0, 0, 0,
	}
	for _, path := range []string{
		"fixtures/f-orientation-1.jpg",
		"fixtures/f-orientation-2.jpg",
		"fixtures/f-orientation-3.jpg",
		"fixtures/f-orientation-4.jpg",
		"fixtures/f-orientation-5.jpg",
		"fixtures/f-orientation-6.jpg",
		"fixtures/f-orientation-7.jpg",
		"fixtures/f-orientation-8.jpg",
	} {
		path := path

		t.Run(path, func(t *testing.T) {
			t.Parallel()

			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			img, err := orientation.Apply(file)
			if err != nil {
				t.Fatal(err)
			}

			i := 0
			for y := 2; y < 42; y += 6 {
				for x := 2; x < 30; x += 6 {
					x := x
					y := y
					w := want[i]
					i++

					t.Run(fmt.Sprintf("at %d %d", x, y), func(t *testing.T) {
						t.Parallel()

						r32, g32, b32, a32 := img.At(x, y).RGBA()
						if a32 != 0xffff {
							t.Errorf("should be opaque")
						}

						r := int(r32 >> 15 & 0x1)
						g := int(g32 >> 15 & 0x1)
						b := int(b32 >> 15 & 0x1)
						if r != w {
							t.Errorf("red got %d, want %d", r, w)
						}
						if g != w {
							t.Errorf("green got %d, want %d", g, w)
						}
						if b != w {
							t.Errorf("blue got %d, want %d", b, w)
						}
					})
				}
			}
		})
	}
}

func TestApplyError(t *testing.T) {
	for _, c := range []struct {
		name string
		file string
		want string
	}{
		// DecodeError
		{
			"non-image file",
			"fixtures/non-image.txt",
			"image: unknown format",
		},
		// FormatError
		{
			"non-jpeg file",
			"fixtures/f-png24.png",
			"format got png, want jpeg",
		},
		// TagError
		{
			"without orientation tag",
			"fixtures/f.jpg",
			"orientation tag does not exist in EXIF: exif: tag \"Orientation\" is not present",
		},
		// OrientError
		{
			"invalid orientation tag",
			"fixtures/f-outofrange.jpg",
			"orientation tag got 9, want 1 to 8",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			file, err := os.Open(c.file)
			if err != nil {
				t.Fatal(err)
			}
			_, err = orientation.Apply(file)
			if err == nil {
				t.Fatal("error does not occur")
			}
			if got := err.Error(); got != c.want {
				t.Errorf("Error() got '%s', want '%s'", got, c.want)
			}
		})
	}
}

func TestDecodeError(t *testing.T) {
	t.Run("DecodeError", func(t *testing.T) {
		for _, c := range []struct {
			name string
			file string
			want string
		}{
			{
				"non-image file",
				"fixtures/non-image.txt",
				"image: unknown format",
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				file, err := os.Open(c.file)
				if err != nil {
					t.Fatal(err)
				}
				_, err = orientation.Decode(file)
				if err == nil {
					t.Fatal("error does not occur")
				}
				if e, ok := err.(*orientation.DecodeError); !ok {
					t.Errorf("type got %T, want *orientation.DecodeError", e)
				}
				if got := err.Error(); got != c.want {
					t.Errorf("Error() got '%s', want '%s'", got, c.want)
				}
			})
		}
	})
	t.Run("FormatError", func(t *testing.T) {
		for _, c := range []struct {
			name string
			file string
			want string
		}{
			{
				"non-jpeg file",
				"fixtures/f-png24.png",
				"format got png, want jpeg",
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				file, err := os.Open(c.file)
				if err != nil {
					t.Fatal(err)
				}
				_, err = orientation.Decode(file)
				if err == nil {
					t.Fatal("error does not occur")
				}
				if e, ok := err.(*orientation.FormatError); !ok {
					t.Errorf("type got %T, want *orientation.FormatError", e)
				}
				if got := err.Error(); got != c.want {
					t.Errorf("Error() got '%s', want '%s'", got, c.want)
				}
			})
		}
	})
}

func TestTagError(t *testing.T) {
	for _, c := range []struct {
		name string
		file string
		want string
	}{
		{
			"without EXIF",
			"fixtures/f-png24.png",
			"fail to decode EXIF",
		},
		{
			"without orientation tag",
			"fixtures/f.jpg",
			"orientation tag does not exist in EXIF",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			file, err := os.Open(c.file)
			if err != nil {
				t.Fatal(err)
			}
			_, err = orientation.Tag(file)
			if err == nil {
				t.Fatal("error does not occur")
			}
			if e, ok := err.(*orientation.TagError); !ok {
				t.Errorf("type got %T, want *orientation.TagError", e)
			}
			if strings.Index(err.Error(), c.want) != 0 {
				t.Errorf("Error() got %s, want %s", err.Error(), c.want)
			}
		})
	}
}

func TestOrient(t *testing.T) {
	t.Parallel()

	want := NewImage(
		image.Pt(3, 5),
		[]int{
			1, 1, 1,
			1, 0, 0,
			1, 1, 0,
			1, 0, 0,
			1, 0, 0,
		},
	)

	for _, c := range []struct {
		tag int
		img image.Image
	}{
		{
			1,
			NewImage(
				image.Pt(3, 5),
				[]int{
					1, 1, 1,
					1, 0, 0,
					1, 1, 0,
					1, 0, 0,
					1, 0, 0,
				},
			),
		},
		{
			2,
			NewImage(
				image.Pt(3, 5),
				[]int{
					1, 1, 1,
					0, 0, 1,
					0, 1, 1,
					0, 0, 1,
					0, 0, 1,
				},
			),
		},
		{
			3,
			NewImage(
				image.Pt(3, 5),
				[]int{
					0, 0, 1,
					0, 0, 1,
					0, 1, 1,
					0, 0, 1,
					1, 1, 1,
				},
			),
		},
		{
			4,
			NewImage(
				image.Pt(3, 5),
				[]int{
					1, 0, 0,
					1, 0, 0,
					1, 1, 0,
					1, 0, 0,
					1, 1, 1,
				},
			),
		},
		{
			5,
			NewImage(
				image.Pt(5, 3),
				[]int{
					1, 1, 1, 1, 1,
					1, 0, 1, 0, 0,
					1, 0, 0, 0, 0,
				},
			),
		},
		{
			6,
			NewImage(
				image.Pt(5, 3),
				[]int{
					1, 0, 0, 0, 0,
					1, 0, 1, 0, 0,
					1, 1, 1, 1, 1,
				},
			),
		},
		{
			7,
			NewImage(
				image.Pt(5, 3),
				[]int{
					0, 0, 0, 0, 1,
					0, 0, 1, 0, 1,
					1, 1, 1, 1, 1,
				},
			),
		},
		{
			8,
			NewImage(
				image.Pt(5, 3),
				[]int{
					1, 1, 1, 1, 1,
					0, 0, 1, 0, 1,
					0, 0, 0, 0, 1,
				},
			),
		},
	} {
		c := c
		t.Run(fmt.Sprintf("tag %d", c.tag), func(t *testing.T) {
			t.Parallel()
			got, err := orientation.Orient(c.img, c.tag)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("\ngot :\n%s\nwant:\n%s", Visualize(got), Visualize(want))
			}
		})
	}
}

func TestOrientError(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		tag int
	}{
		{
			-1,
		},
		{
			0,
		},
		{
			9,
		},
	} {
		c := c
		t.Run(fmt.Sprintf("tag %d", c.tag), func(t *testing.T) {
			t.Parallel()
			_, err := orientation.Orient(nil, c.tag)
			if err == nil {
				t.Fatal(err)
			}
			if _, ok := err.(*orientation.OrientError); !ok {
				t.Errorf("got %T, want orientation.OrientError", err)
			}
		})
	}
}

func TestOrient1to8(t *testing.T) {
	t.Parallel()

	want := NewImage(
		image.Pt(3, 5),
		[]int{
			1, 1, 1,
			1, 0, 0,
			1, 1, 0,
			1, 0, 0,
			1, 0, 0,
		},
	)

	for _, c := range []struct {
		name string
		fn   func(image.Image) image.Image
		img  image.Image
	}{
		{
			"Orient1",
			orientation.Orient1,
			NewImage(
				image.Pt(3, 5),
				[]int{
					1, 1, 1,
					1, 0, 0,
					1, 1, 0,
					1, 0, 0,
					1, 0, 0,
				},
			),
		},
		{
			"Orient2",
			orientation.Orient2,
			NewImage(
				image.Pt(3, 5),
				[]int{
					1, 1, 1,
					0, 0, 1,
					0, 1, 1,
					0, 0, 1,
					0, 0, 1,
				},
			),
		},
		{
			"Orient3",
			orientation.Orient3,
			NewImage(
				image.Pt(3, 5),
				[]int{
					0, 0, 1,
					0, 0, 1,
					0, 1, 1,
					0, 0, 1,
					1, 1, 1,
				},
			),
		},
		{
			"Orient4",
			orientation.Orient4,
			NewImage(
				image.Pt(3, 5),
				[]int{
					1, 0, 0,
					1, 0, 0,
					1, 1, 0,
					1, 0, 0,
					1, 1, 1,
				},
			),
		},
		{
			"Orient5",
			orientation.Orient5,
			NewImage(
				image.Pt(5, 3),
				[]int{
					1, 1, 1, 1, 1,
					1, 0, 1, 0, 0,
					1, 0, 0, 0, 0,
				},
			),
		},
		{
			"Orient6",
			orientation.Orient6,
			NewImage(
				image.Pt(5, 3),
				[]int{
					1, 0, 0, 0, 0,
					1, 0, 1, 0, 0,
					1, 1, 1, 1, 1,
				},
			),
		},
		{
			"Orient7",
			orientation.Orient7,
			NewImage(
				image.Pt(5, 3),
				[]int{
					0, 0, 0, 0, 1,
					0, 0, 1, 0, 1,
					1, 1, 1, 1, 1,
				},
			),
		},
		{
			"Orient8",
			orientation.Orient8,
			NewImage(
				image.Pt(5, 3),
				[]int{
					1, 1, 1, 1, 1,
					0, 0, 1, 0, 1,
					0, 0, 0, 0, 1,
				},
			),
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got := c.fn(c.img)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("\ngot :\n%s\nwant:\n%s", Visualize(got), Visualize(want))
			}
		})
	}
}
