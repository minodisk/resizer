package processor_test

import (
	"bytes"
	"image"
	"math"
	"testing"

	"github.com/go-microservices/resizer/processor"
	"github.com/go-microservices/resizer/storage"
)

const (
	u   = 3
	png = "../fixtures/f-png24.png"
)

var (
	formats = []string{
		"../fixtures/f.jpg",
		"../fixtures/f-png8.png",
		"../fixtures/f-png24.png",
		"../fixtures/f.gif",
	}
	orientations = []string{
		"../fixtures/f-orientation-1.jpg",
		"../fixtures/f-orientation-2.jpg",
		"../fixtures/f-orientation-3.jpg",
		"../fixtures/f-orientation-4.jpg",
		"../fixtures/f-orientation-5.jpg",
		"../fixtures/f-orientation-6.jpg",
		"../fixtures/f-orientation-7.jpg",
		"../fixtures/f-orientation-8.jpg",
	}
	raw = []int{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
)

func diff(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}

func isNear(a, b uint32) bool {
	return diff(a, b) <= math.MaxUint8*4
}

func evalPixels(t *testing.T, i image.Image, p image.Point, colors []int) {
	for y := 0; y < p.Y; y++ {
		for x := 0; x < p.X; x++ {
			var er, eg, eb, ea uint32
			e := colors[p.X*y+x]
			if e == 1 {
				er = 0xffff
				eg = 0xffff
				eb = 0xffff
			} else {
				er = 0
				eg = 0
				eb = 0
			}
			ea = 0xffff

			a := i.At(u/2>>0+u*x, u/2>>0+u*y)
			ar, ag, ab, aa := a.RGBA()

			if !(isNear(er, ar) && isNear(eg, ag) && isNear(eb, ab) && isNear(ea, aa)) {
				t.Errorf(
					"wrong color at (%d, %d) expected {%d, %d, %d, %d}, but actual {%d, %d, %d, %d}",
					x, y,
					er, eg, eb, ea,
					ar, ag, ab, aa,
				)
			}
		}
	}
}

func eval(t *testing.T, path string, f storage.Image, size image.Point, colors []int) string {
	var b []byte
	w := bytes.NewBuffer(b)
	p := processor.New()
	f.ValidatedWidth *= u
	f.ValidatedHeight *= u
	pixels, err := p.Preprocess(path)
	if err != nil {
		t.Fatalf("cannot preprocess image: error=%v", err)
	}

	f, err = f.Normalize(pixels.Bounds().Size())
	if err != nil {
		t.Fatalf("fail to normalize: error=%v", err)
	}

	if _, err := p.Process(pixels, w, f); err != nil {
		t.Fatalf("cannot process image: %v", err)
		return ""
	}

	r := bytes.NewReader(w.Bytes())
	img, format, err := image.Decode(r)
	if err != nil {
		t.Fatalf("cannot decode image: %v", err)
		return ""
	}

	expectedSize := size.Mul(u)
	rect := img.Bounds()
	actualSize := rect.Size()
	if !actualSize.Eq(expectedSize) {
		t.Fatalf("wrong size expected %v, but actual %v", expectedSize, actualSize)
		return ""
	}

	evalPixels(t, img, size, colors)

	return format
}

func TestFormats(t *testing.T) {
	size := image.Point{5, 7}
	colors := []int{
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		0, 1, 0, 0, 0,
		0, 1, 1, 1, 0,
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 0, 0, 0,
	}
	for _, format := range []string{storage.FormatPng, storage.FormatJpeg, storage.FormatGif} {
		f := storage.Image{
			ValidatedMethod:  storage.MethodNormal,
			ValidatedWidth:   5,
			ValidatedHeight:  7,
			ValidatedFormat:  format,
			ValidatedQuality: 100,
		}
		for _, path := range formats {
			eval(t, path, f, size, colors)
		}
	}
}

func TestOrientations(t *testing.T) {
	f := storage.Image{
		ValidatedMethod: storage.MethodNormal,
		ValidatedWidth:  5,
		ValidatedHeight: 7,
		ValidatedFormat: storage.FormatPng,
	}
	size := image.Point{5, 7}
	colors := []int{
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		0, 1, 0, 0, 0,
		0, 1, 1, 1, 0,
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 0, 0, 0,
	}
	for _, path := range orientations {
		eval(t, path, f, size, colors)
	}
}

func TestFormatNormal(t *testing.T) {
	size := image.Point{5, 7}
	colors := []int{
		0, 0, 0, 0, 0,
		0, 1, 1, 1, 0,
		0, 1, 0, 0, 0,
		0, 1, 1, 1, 0,
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 0, 0, 0, 0,
	}
	eval(t, png, storage.Image{
		ValidatedMethod: storage.MethodNormal,
		ValidatedWidth:  5,
		ValidatedHeight: 100,
		ValidatedFormat: storage.FormatPng,
	}, size, colors)
	eval(t, png, storage.Image{
		ValidatedMethod: storage.MethodNormal,
		ValidatedWidth:  100,
		ValidatedHeight: 7,
		ValidatedFormat: storage.FormatPng,
	}, size, colors)
}

func TestFormatThumbnail(t *testing.T) {
	eval(t, png, storage.Image{
		ValidatedMethod: storage.MethodThumbnail,
		ValidatedWidth:  3,
		ValidatedHeight: 7,
		ValidatedFormat: storage.FormatPng,
	}, image.Point{3, 7}, []int{
		0, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
		1, 0, 0,
		1, 0, 0,
		0, 0, 0,
	})

	eval(t, png, storage.Image{
		ValidatedMethod: storage.MethodThumbnail,
		ValidatedWidth:  5,
		ValidatedHeight: 3,
		ValidatedFormat: storage.FormatPng,
	}, image.Point{5, 3}, []int{
		0, 1, 0, 0, 0,
		0, 1, 1, 1, 0,
		0, 1, 0, 0, 0,
	})

	eval(t, png, storage.Image{
		ValidatedMethod: storage.MethodThumbnail,
		ValidatedWidth:  100,
		ValidatedHeight: 100,
		ValidatedFormat: storage.FormatPng,
	}, image.Point{10, 14}, []int{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	})

	eval(t, png, storage.Image{
		ValidatedMethod: storage.MethodThumbnail,
		ValidatedWidth:  6,
		ValidatedHeight: 100,
		ValidatedFormat: storage.FormatPng,
	}, image.Point{6, 14}, []int{
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1,
		1, 1, 0, 0, 0, 0,
		1, 1, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1,
		1, 1, 0, 0, 0, 0,
		1, 1, 0, 0, 0, 0,
		1, 1, 0, 0, 0, 0,
		1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	})

	eval(t, png, storage.Image{
		ValidatedMethod: storage.MethodThumbnail,
		ValidatedWidth:  2,
		ValidatedHeight: 100,
		ValidatedFormat: storage.FormatPng,
	}, image.Point{2, 14}, []int{
		0, 0,
		0, 0,
		1, 1,
		1, 1,
		0, 0,
		0, 0,
		1, 1,
		1, 1,
		0, 0,
		0, 0,
		0, 0,
		0, 0,
		0, 0,
		0, 0,
	})

	eval(t, png, storage.Image{
		ValidatedMethod: storage.MethodThumbnail,
		ValidatedWidth:  100,
		ValidatedHeight: 10,
		ValidatedFormat: storage.FormatPng,
	}, image.Point{10, 10}, []int{
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
	})

	eval(t, png, storage.Image{
		ValidatedMethod: storage.MethodThumbnail,
		ValidatedWidth:  100,
		ValidatedHeight: 6,
		ValidatedFormat: storage.FormatPng,
	}, image.Point{10, 6}, []int{
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
	})

	eval(t, png, storage.Image{
		ValidatedMethod: storage.MethodThumbnail,
		ValidatedWidth:  100,
		ValidatedHeight: 2,
		ValidatedFormat: storage.FormatPng,
	}, image.Point{10, 2}, []int{
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 0, 0,
	})
}
