package orientation

import (
	"fmt"
	"image"
	"io"

	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
)

// Orientation tags.
const (
	Orientation1 = 1
	Orientation2 = 2
	Orientation3 = 3
	Orientation4 = 4
	Orientation5 = 5
	Orientation6 = 6
	Orientation7 = 7
	Orientation8 = 8
)

// Apply read orientation tag from file at filename
// and applies it to src image.
// Returns oriented image and any error occurred.
func Apply(r io.Reader, src image.Image) (image.Image, error) {
	o, err := Read(r)
	// When fails to read orientation tag from EXIF, finishes preprocessing.
	if err != nil {
		return src, nil
	}
	return Orient(src, o)
}

// Read decodes orientation tag in EXIF from r.
// Returns orientation tag and any error occurred.
func Read(r io.Reader) (int, error) {
	e, err := exif.Decode(r)
	if err != nil {
		return 0, errors.Wrap(err, "fail to decode EXIF")
	}
	tag, err := e.Get(exif.Orientation)
	if err != nil {
		return 0, errors.Wrap(err, "orientation tag doesn't exist in EXIF")
	}
	o, err := tag.Int(0)
	if err != nil {
		return 0, errors.Wrap(err, "orientation tag isn't int")
	}
	return o, nil
}

// Orient applies orientation tag o to s.
// Returns a new oriented image.
// When o is invalid as orientation tag, returns error.
func Orient(s image.Image, o int) (image.Image, error) {
	// rect := s.Bounds()
	// if o >= 5 && o <= 8 {
	// 	rect = SwapRect(rect)
	// }
	// d := image.NewRGBA64(rect)
	// a := affines[o]
	// a.TransformCenter(d, s, interp.Bilinear)
	// return d, nil

	switch o {
	case Orientation1:
		return Orient1(s), nil
	case Orientation2:
		return Orient2(s), nil
	case Orientation3:
		return Orient3(s), nil
	case Orientation4:
		return Orient4(s), nil
	case Orientation5:
		return Orient5(s), nil
	case Orientation6:
		return Orient6(s), nil
	case Orientation7:
		return Orient7(s), nil
	case Orientation8:
		return Orient8(s), nil
	default:
		return nil, fmt.Errorf("invalid orientation tag %d", o)
	}
}

// Orient1 returns s.
// 	111    111
// 	100    100
// 	110 -> 110
// 	100    100
// 	100    100
func Orient1(s image.Image) image.Image {
	return s
}

// Orient2 returns oriented s that is replaced
// right side to left side.
// 	111    111
// 	001    100
// 	011 -> 110
// 	001    100
// 	001    100
func Orient2(s image.Image) image.Image {
	sr := s.Bounds()
	dr := sr
	d := image.NewRGBA(dr)
	for y := sr.Min.Y; y < sr.Max.Y; y++ {
		for x := sr.Min.X; x < sr.Max.X; x++ {
			dx := sr.Max.X - 1 - x
			dy := y - sr.Min.Y
			d.Set(dr.Min.X+dx, dr.Min.Y+dy, s.At(x, y))
		}
	}
	return d
}

// Orient3 returns oriented s that is replaced
// bottom to top and right side to left side.
// 	001    111
// 	001    100
// 	011 -> 110
// 	001    100
// 	111    100
func Orient3(s image.Image) image.Image {
	sr := s.Bounds()
	dr := sr
	d := image.NewRGBA(dr)
	for y := sr.Min.Y; y < sr.Max.Y; y++ {
		for x := sr.Min.X; x < sr.Max.X; x++ {
			dx := sr.Max.X - 1 - x
			dy := sr.Max.Y - 1 - y
			d.Set(dr.Min.X+dx, dr.Min.Y+dy, s.At(x, y))
		}
	}
	return d
}

// Orient4 returns oriented s that is replaced
// bottom to top.
// 	100    111
// 	100    100
// 	110 -> 110
// 	100    100
// 	111    100
func Orient4(s image.Image) image.Image {
	sr := s.Bounds()
	dr := sr
	d := image.NewRGBA(dr)
	for y := sr.Min.Y; y < sr.Max.Y; y++ {
		for x := sr.Min.X; x < sr.Max.X; x++ {
			dx := x - sr.Min.X
			dy := sr.Max.Y - 1 - y
			d.Set(dr.Min.X+dx, dr.Min.Y+dy, s.At(x, y))
		}
	}
	return d
}

// Orient5 returns oriented s that is replaced
// left side to top and top to left side.
// 	         111
// 	11111    100
// 	10100 -> 110
// 	10000    100
// 	         100
func Orient5(s image.Image) image.Image {
	sr := s.Bounds()
	dr := SwapRect(sr)
	d := image.NewRGBA(dr)
	for y := sr.Min.Y; y < sr.Max.Y; y++ {
		for x := sr.Min.X; x < sr.Max.X; x++ {
			dx := y - sr.Min.Y
			dy := x - sr.Min.X
			d.Set(dr.Min.X+dx, dr.Min.Y+dy, s.At(x, y))
		}
	}
	return d
}

// Orient6 returns oriented s that is replaced
// left side to top and bottom to left side.
// 	         111
// 	10000    100
// 	10100 -> 110
// 	11111    100
// 	         100
func Orient6(s image.Image) image.Image {
	sr := s.Bounds()
	dr := SwapRect(sr)
	d := image.NewRGBA(dr)
	for y := sr.Min.Y; y < sr.Max.Y; y++ {
		for x := sr.Min.X; x < sr.Max.X; x++ {
			dx := sr.Max.Y - 1 - y
			dy := x - sr.Min.X
			d.Set(dr.Min.X+dx, dr.Min.Y+dy, s.At(x, y))
		}
	}
	return d
}

// Orient7 returns oriented s that is replaced
// right side to top and bottom to left side.
// 	         111
// 	00001    100
// 	00101 -> 110
// 	11111    100
// 	         100
func Orient7(s image.Image) image.Image {
	sr := s.Bounds()
	dr := SwapRect(sr)
	d := image.NewRGBA(dr)
	for y := sr.Min.Y; y < sr.Max.Y; y++ {
		for x := sr.Min.X; x < sr.Max.X; x++ {
			dx := sr.Max.Y - 1 - y
			dy := sr.Max.X - 1 - x
			d.Set(dr.Min.X+dx, dr.Min.Y+dy, s.At(x, y))
		}
	}
	return d
}

// Orient8 returns oriented s that is replaced
// right side to top and bottom to left side.
// 	         111
// 	11111    100
// 	00101 -> 110
// 	00001    100
// 	         100
func Orient8(s image.Image) image.Image {
	sr := s.Bounds()
	dr := SwapRect(sr)
	d := image.NewRGBA(dr)
	for y := sr.Min.Y; y < sr.Max.Y; y++ {
		for x := sr.Min.X; x < sr.Max.X; x++ {
			dx := y - sr.Min.Y
			dy := sr.Max.X - 1 - x
			d.Set(dr.Min.X+dx, dr.Min.Y+dy, s.At(x, y))
		}
	}
	return d
}

// SwapRect swaps width and height in rectangle.
// Returns new rectangle.
func SwapRect(r image.Rectangle) image.Rectangle {
	s := r.Size()
	return image.Rectangle{r.Min, image.Point{s.Y, s.X}}
}
