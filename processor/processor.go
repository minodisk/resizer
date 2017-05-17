// Package processor では画像処理が実装されています。
//
// このアプリケーションのメインの機能である画像のリサイズを処理します。
package processor

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"sync"

	"github.com/minodisk/orientation"
	"github.com/minodisk/resizer/input"
	"github.com/minodisk/resizer/storage"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

var (
	mutex sync.Mutex
)

type Processor struct{}

func New() *Processor {
	return &Processor{}
}

func (p *Processor) Process(path string, w io.Writer, f storage.Image) (*image.Point, error) {
	c := make(chan Result)
	go p.process(&mutex, c, path, w, f)
	for res := range c {
		return res.Point, res.Error
	}
	return nil, nil
}

type Result struct {
	Point *image.Point
	Error error
}

func (p *Processor) process(m *sync.Mutex, c chan Result, path string, w io.Writer, f storage.Image) {
	m.Lock()
	defer m.Unlock()

	i, err := p.Preprocess(path)
	if err != nil {
		c <- Result{nil, err}
	}
	pt, err := p.Resize(i, w, f)
	c <- Result{pt, err}
}

// Preprocess load image and EXIF from file at filename.
// When orientation tag exists in EXIF, orient pixels in
// image.
func (self *Processor) Preprocess(filename string) (image.Image, error) {
	src, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := orientation.Apply(src)
	if err != nil {
		if err, ok := err.(*orientation.DecodeError); ok {
			return nil, errors.Wrap(err, "fail to apply orientation")
		}
	}

	return dst, nil
}

// Load decodes image from file at filename.
// It returns decoded image, the format of image, and any error occurred.
func Load(filename string) (image.Image, string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", errors.Wrap(err, "fail to open file")
	}
	defer f.Close()
	src, format, err := image.Decode(f)
	if err != nil {
		return src, format, errors.Wrap(err, "fail to decode file as image")
	}
	return src, format, nil
}

// Process writes resized i to w with options in f.
// Returns the size of resized i and any error occurred.
func (self *Processor) Resize(i image.Image, w io.Writer, f storage.Image) (*image.Point, error) {
	log.Printf("dest image: %+v\n", f)

	var ir image.Image
	switch f.ValidatedMethod {
	default:
		return nil, fmt.Errorf("Unsupported method: %s", f.ValidatedMethod)
	case input.MethodContain:
		ir = resize.Resize(uint(f.DestWidth), uint(f.DestHeight), i, resize.Lanczos3)
	case input.MethodCover:
		cr := image.Rect(0, 0, f.CanvasWidth, f.CanvasHeight)
		src := resize.Resize(uint(f.DestWidth), uint(f.DestHeight), i, resize.Lanczos3)
		dst := image.NewRGBA(cr)
		draw.Draw(dst, cr, src, image.Point{int((f.DestWidth - f.CanvasWidth) / 2), int((f.DestHeight - f.CanvasHeight) / 2)}, draw.Src)
		ir = dst
	}

	switch f.ValidatedFormat {
	default:
		return nil, fmt.Errorf("Unsupported format: %s", f.ValidatedFormat)
	case input.FormatJPEG:
		if err := jpeg.Encode(w, ir, &jpeg.Options{Quality: int(f.ValidatedQuality)}); err != nil {
			return nil, err
		}
	case input.FormatPNG:
		e := png.Encoder{CompressionLevel: png.DefaultCompression}
		if err := e.Encode(w, ir); err != nil {
			return nil, err
		}
	case input.FormatGIF:
		if err := gif.Encode(w, ir, &gif.Options{NumColors: 256}); err != nil {
			return nil, err
		}
	}

	size := ir.Bounds().Size()
	return &size, nil
}
