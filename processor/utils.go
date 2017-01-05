package processor

import (
	"bufio"
	"image"
	"io"
	"os"

	"github.com/BurntSushi/graphics-go/graphics/interp"
	"github.com/go-microservices/resizer/input"
	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
)

func load(path string) (image.Image, string, error) {
	// ファイルをデコードする
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	src, format, err := image.Decode(bufio.NewReader(f))
	if err != nil {
		return src, format, errors.Wrap(err, "fail to decode image")
	}

	return src, format, nil
}

func applyOrientation(src image.Image, format, path string) (image.Image, error) {
	// JPEG以外ならピクセルをそのまま返す
	if format != input.FormatJPEG {
		return src, nil
	}

	// JPEGならEXIFの回転情報をピクセルに反映して返す
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	o, err := getOrientation(bufio.NewReader(f))
	// When fails to get orientation tag from EXIF, finishes preprocessing.
	if err != nil {
		return src, nil
	}

	dst, err := rotate(src, o)
	if err != nil {
		return nil, errors.Wrap(err, "fail to rotate image")
	}

	return dst, nil
}

func getOrientation(r io.Reader) (int, error) {
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

func rotate(src image.Image, o int) (image.Image, error) {
	// if o == 1 {
	// 	return i, nil
	// }

	rect := src.Bounds()
	// orientation=5~8 なら画像サイズの縦横を入れ替える
	if o >= 5 && o <= 8 {
		rect = RotateRect(rect)
	}
	d := image.NewRGBA64(rect)
	a := affines[o]
	a.TransformCenter(d, src, interp.Bilinear)

	return d, nil
}
