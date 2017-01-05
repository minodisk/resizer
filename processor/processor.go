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
	"sync"

	"github.com/go-microservices/resizer/input"
	"github.com/go-microservices/resizer/log"
	"github.com/go-microservices/resizer/storage"
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

// Preprocess はリサイズ処理の前処理を行います。
// 画像をデコードし、JPEGのEXIFデータの回転情報をピクセルに反映して返します。
func (self *Processor) Preprocess(path string) (image.Image, error) {
	src, format, err := load(path)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to load image")
	}

	dst, err := applyOrientation(src, format, path)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to apply orientation")
	}

	return dst, nil
}

// Process はリサイズ処理を行い、エンコードしたデータを返します。
func (self *Processor) Resize(i image.Image, w io.Writer, f storage.Image) (*image.Point, error) {
	log.Printf("dest image: %+v\n", f)

	var ir image.Image
	switch f.ValidatedMethod {
	default:
		return nil, fmt.Errorf("Unsupported method: %s", f.ValidatedMethod)
	case input.MethodNormal:
		ir = resize.Resize(uint(f.DestWidth), uint(f.DestHeight), i, resize.Lanczos3)
	case input.MethodThumbnail:
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
		if err := jpeg.Encode(w, ir, &jpeg.Options{int(f.ValidatedQuality)}); err != nil {
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
