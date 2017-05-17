package orientation

import (
	"bytes"
	"image"
	"io"

	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
)

var (
	tagOrientMapper = map[int]func(image.Image) image.Image{
		1: Orient1,
		2: Orient2,
		3: Orient3,
		4: Orient4,
		5: Orient5,
		6: Orient6,
		7: Orient7,
		8: Orient8,
	}
)

// Apply reflects the image's EXIF orientation tag on the image pixels.
// If the image can not be decoded, it returns nil and an error.
// If an error occurs during subsequent processing,
// it returns the decoded image and an error.
func Apply(r io.Reader) (image.Image, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	src, err := Decode(tee)
	if err != nil {
		return src, err
	}
	tag, err := Tag(&buf)
	if err != nil {
		return src, err
	}
	dst, err := Orient(src, tag)
	if err != nil {
		return src, err
	}
	return dst, nil
}

// Decode decodes and returns the image.
// If decoding fails, it returns nil and DecodeError.
// If the format is not JPEG, it returns the decoded image and FormatError.
func Decode(r io.Reader) (image.Image, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, &DecodeError{err}
	}
	if format != "jpeg" {
		return img, &FormatError{errors.Errorf("format got %s, want jpeg", format)}
	}
	return img, nil
}

// Tag decodes the orientation tag from the image EXIF and returns it.
// If the EXIF can not be decoded, the orientation tag does not exist,
// or the orientation tag is not int, it returns TagError.
func Tag(r io.Reader) (int, error) {
	e, err := exif.Decode(r)
	if err != nil {
		return 0, &TagError{errors.Wrap(err, "fail to decode EXIF")}
	}
	tag, err := e.Get(exif.Orientation)
	if err != nil {
		return 0, &TagError{errors.Wrap(err, "orientation tag does not exist in EXIF")}
	}
	o, err := tag.Int(0)
	if err != nil {
		return 0, &TagError{errors.Wrap(err, "orientation tag is not int")}
	}
	return o, nil
}

// Orient reflects the rotation indicated by the tag in the img.
// If the specified orientation tag is unknown, returns OrientError.
func Orient(img image.Image, tag int) (image.Image, error) {
	fn, ok := tagOrientMapper[tag]
	if !ok {
		return nil, &OrientError{errors.Errorf("orientation tag got %d, want 1 to 8", tag)}
	}
	return fn(img), nil
}

// Orient1 returns s.
// Process like this:
// 	111    111
// 	100    100
// 	110 -> 110
// 	100    100
// 	100    100
func Orient1(s image.Image) image.Image {
	return s
}

// Orient2 returns a new image in which the right side replaced by the left side.
// Process like this:
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

// Orient3 returns a new image in which the bottom side is replaced by the upper
// side and the right side replaced by the left side.
// Process like this:
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

// Orient4 returns a new image in which the bottom side is replaced by the upper
// side.
// Process like this:
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

// Orient5 returns a new image in which the left side is replaced by the upper
// side and the upper side is replaced by the left side.
// Process like this:
// 	         111
// 	11111    100
// 	10100 -> 110
// 	10000    100
// 	         100
func Orient5(s image.Image) image.Image {
	sr := s.Bounds()
	dr := SwapSides(sr)
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

// Orient6 returns a new image in which the left side is replaced by the upper
// side and the bottom side is replaced by the left side.
// Process like this:
// 	         111
// 	10000    100
// 	10100 -> 110
// 	11111    100
// 	         100
func Orient6(s image.Image) image.Image {
	sr := s.Bounds()
	dr := SwapSides(sr)
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

// Orient7 returns a new image in which the right side is replaced by the upper
// side and the bottom side is replaced by the left side.
// Process like this:
// 	         111
// 	00001    100
// 	00101 -> 110
// 	11111    100
// 	         100
func Orient7(s image.Image) image.Image {
	sr := s.Bounds()
	dr := SwapSides(sr)
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

// Orient8 returns a new image in which the right side is replaced by the upper
// side and the upper side is replaced by the left side.
// Process like this:
// 	         111
// 	11111    100
// 	00101 -> 110
// 	00001    100
// 	         100
func Orient8(s image.Image) image.Image {
	sr := s.Bounds()
	dr := SwapSides(sr)
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

// SwapSides returns a new rectangle that swaps the width and height
// of the rectangle.
func SwapSides(r image.Rectangle) image.Rectangle {
	s := r.Size()
	return image.Rectangle{r.Min, image.Point{s.Y, s.X}}
}
