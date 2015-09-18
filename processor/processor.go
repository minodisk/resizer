// Package processor では画像処理が実装されています。
//
// このアプリケーションのメインの機能である画像のリサイズを処理します。
package processor

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"reflect"

	"code.google.com/p/graphics-go/graphics/interp"

	"github.com/go-microservices/resizer/log"
	"github.com/go-microservices/resizer/storage"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
)

type Processor struct{}

func New() *Processor {
	return &Processor{}
}

// Preprocess はリサイズ処理の前処理を行います。
// 画像をデコードし、jpegのEXIFの回転情報をピクセルに反映して返します。
func (self *Processor) Preprocess(path string) (image.Image, error) {
	t := log.Start()
	defer log.End(t)

	// ファイルをデコードする
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	i, format, err := image.Decode(bufio.NewReader(f))
	if err != nil {
		log.Println("fail to decode")
		return nil, err
	}

	// jpeg以外ならピクセルをそのまま返す
	if format != storage.FormatJpeg {
		return i, nil
	}

	// jpegならEXIFの回転情報をピクセルに反映して返す
	f, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	io, err := orient(bufio.NewReader(f), i)
	// 回転に失敗した場合、元のピクセルを返す
	if err != nil {
		log.Println("cancel to apply orientation")
		return i, nil
	}
	log.Debugf("%s -> %s", reflect.TypeOf(i), reflect.TypeOf(io))
	return io, nil
}

func orient(r io.Reader, i image.Image) (image.Image, error) {
	t := log.Start()
	defer log.End(t)

	e, err := exif.Decode(r)
	if err != nil {
		log.Printf("fail to decode exif: %v", err)
		return nil, err
	}
	tag, err := e.Get(exif.Orientation)
	// Orientationタグが存在しない場合、処理を完了する
	if err != nil {
		log.Println("oritentation tag doesn't exist")
		return nil, err
	}
	o, err := tag.Int(0)
	if err != nil {
		log.Println("oritentation tag is't int")
		return nil, err
	}

	rect := i.Bounds()
	// orientation=5~8 なら画像サイズの縦横を入れ替える
	if o >= 5 && o <= 8 {
		rect = RotateRect(rect)
	}
	d := image.NewRGBA64(rect)
	a := affines[o]
	a.TransformCenter(d, i, interp.Bilinear)

	return d, nil
}

// Process はリサイズ処理を行い、エンコードしたデータを返します。
func (self *Processor) Process(i image.Image, w io.Writer, f storage.Image) (*image.Point, error) {
	var ir image.Image
	switch f.ValidatedMethod {
	default:
		return nil, fmt.Errorf("Unsupported method: %s", f.ValidatedMethod)
	case storage.MethodNormal:
		ir = resize.Resize(uint(f.DestWidth), uint(f.DestHeight), i, resize.Lanczos3)
	case storage.MethodThumbnail:
		cr := image.Rect(0, 0, f.CanvasWidth, f.CanvasHeight)
		src := resize.Resize(uint(f.DestWidth), uint(f.DestHeight), i, resize.Lanczos3)
		dst := image.NewRGBA(cr)
		draw.Draw(dst, cr, src, image.Point{int((f.DestWidth - f.CanvasWidth) / 2), int((f.DestHeight - f.CanvasHeight) / 2)}, draw.Src)
		ir = dst
	}

	switch f.ValidatedFormat {
	default:
		return nil, fmt.Errorf("Unsupported format: %s", f.ValidatedFormat)
	case storage.FormatJpeg:
		if err := jpeg.Encode(w, ir, &jpeg.Options{int(f.ValidatedQuality)}); err != nil {
			return nil, err
		}
	case storage.FormatPng:
		e := png.Encoder{CompressionLevel: png.DefaultCompression}
		if err := e.Encode(w, ir); err != nil {
			return nil, err
		}
	case storage.FormatGif:
		if err := gif.Encode(w, ir, &gif.Options{NumColors: 256}); err != nil {
			return nil, err
		}
	}

	size := ir.Bounds().Size()
	return &size, nil
}
