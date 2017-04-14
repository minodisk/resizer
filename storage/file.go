package storage

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"errors"
	"fmt"
	"image"
	"math"
	"time"

	"github.com/minodisk/resizer/input"
	uuid "github.com/satori/go.uuid"
)

type Image struct {
	ID               uint64
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ValidatedURL     string `sql:"type:text"`
	ValidatedMethod  string
	ValidatedFormat  string
	ValidatedWidth   int
	ValidatedHeight  int
	ValidatedQuality int
	ValidatedHash    string `sql:"size:32;index"`
	DestWidth        int
	DestHeight       int
	CanvasWidth      int
	CanvasHeight     int
	NormalizedHash   string `sql:"size:32;index"`
	ContentType      string `sql:"size:80"`
	ETag             string `sql:"size:32"`
	Filename         string
}

// New はクエリのマップ q から File を作成する。
// デフォルト値の存在するパラメーターに値が設定されていない場合は、デフォルト値を設定する。
func NewImage(input input.Input) (Image, error) {
	return Image{
		ValidatedURL:     input.URL,
		ValidatedMethod:  input.Method,
		ValidatedWidth:   input.Width,
		ValidatedHeight:  input.Height,
		ValidatedFormat:  input.Format,
		ValidatedQuality: input.Quality,
	}.serializeValidatedProps()
}

func (i Image) serializeValidatedProps() (Image, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(Image{
		ValidatedURL:     i.ValidatedURL,
		ValidatedMethod:  i.ValidatedMethod,
		ValidatedFormat:  i.ValidatedFormat,
		ValidatedQuality: i.ValidatedQuality,
		ValidatedWidth:   i.ValidatedWidth,
		ValidatedHeight:  i.ValidatedHeight,
	}); err != nil {
		return i, err
	}
	sum := md5.Sum(buf.Bytes())
	i.ValidatedHash = fmt.Sprintf("%x", sum)
	return i, nil
}

// Normalize は width, height の片方が 0 の場合に、
// 長さの指定されている辺と元画像のアスペクト比から長さの指定されていない辺の長さを推測する。
func (i Image) Normalize(src image.Point) (Image, error) {
	var err error
	i, err = i.normalize(src)
	if err != nil {
		return i, err
	}
	return i.serializeNormalizedProps()
}

func (i Image) normalize(src image.Point) (Image, error) {
	// 元画像の辺のいずれかが0ならアスペクト比の算出が不可能なのでエラーする。
	if src.X == 0 || src.Y == 0 {
		return i, errors.New("source size must not be a zero")
	}
	// 目的サイズが0の場合は許可しないのでエラーする。
	if i.ValidatedWidth == 0 && i.ValidatedHeight == 0 {
		return i, errors.New("target size must not be a zero")
	}

	// Width, Height どちらかが0なら元画像の縦横比を保つように補完する。
	var sx, sy, sr, dx, dy float64
	sx = float64(src.X)
	sy = float64(src.Y)
	sr = sx / sy
	dx = float64(i.ValidatedWidth)
	dy = float64(i.ValidatedHeight)
	if dx == 0 {
		dx = dy * sr
	}
	if dy == 0 {
		dy = dx / sr
	}

	// 目的サイズが元サイズより大きければ、
	// 他の条件を無視して元画像のサイズを採用する。
	if dx >= sx && dy >= sy {
		i.DestWidth = src.X
		i.DestHeight = src.Y
		i.CanvasWidth = i.DestWidth
		i.CanvasHeight = i.DestHeight
		return i, nil
	}

	// 目的サイズの片方の辺が0なら、
	// methodを無視してサイズを決定する。
	if i.ValidatedWidth == 0 || i.ValidatedHeight == 0 {
		i.DestWidth = int(dx)
		i.DestHeight = int(dy)
		i.CanvasWidth = i.DestWidth
		i.CanvasHeight = i.DestHeight
		return i, nil
	}

	switch i.ValidatedMethod {
	default:
		return i, fmt.Errorf("method %s isn't supported", i.ValidatedMethod)
	// 目的のサイズに完全に収まる大きさを計算する
	case input.MethodNormal:
		dr := dx / dy
		if dr == sr {
			i.DestWidth = i.ValidatedWidth
			i.DestHeight = i.ValidatedHeight
			i.CanvasWidth = i.DestWidth
			i.CanvasHeight = i.DestHeight
			return i, nil
		} else if dr > sr {
			i.DestWidth = int(dy * sr)
			i.DestHeight = i.ValidatedHeight
			i.CanvasWidth = i.DestWidth
			i.CanvasHeight = i.DestHeight
			return i, nil
		} else {
			i.DestWidth = i.ValidatedWidth
			i.DestHeight = int(dx / sr)
			i.CanvasWidth = i.DestWidth
			i.CanvasHeight = i.DestHeight
			return i, nil
		}
	// 目的のサイズを埋める大きさを計算する
	// 最終的なサイズが目的のサイズをはみだしてもよい
	case input.MethodThumbnail:
		rx := math.Min(1, dx/sx)
		ry := math.Min(1, dy/sy)
		if rx > ry {
			dx = sx * rx
			dy = sy * rx
		} else {
			dx = sx * ry
			dy = sy * ry
		}
		i.DestWidth = int(dx)
		i.DestHeight = int(dy)
		i.CanvasWidth = int(math.Min(float64(i.ValidatedWidth), dx))
		i.CanvasHeight = int(math.Min(float64(i.ValidatedHeight), dy))
		return i, nil
	}
}

func (i Image) serializeNormalizedProps() (Image, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(Image{
		ValidatedURL:     i.ValidatedURL,
		ValidatedMethod:  i.ValidatedMethod,
		ValidatedFormat:  i.ValidatedFormat,
		ValidatedQuality: i.ValidatedQuality,
		DestWidth:        i.DestWidth,
		DestHeight:       i.DestHeight,
	}); err != nil {
		return i, err
	}
	sum := md5.Sum(buf.Bytes())
	i.NormalizedHash = fmt.Sprintf("%x", sum)
	return i, nil
}

func (i Image) CreateFilename() string {
	id := uuid.NewV4().String()
	return fmt.Sprintf("%s.%s", id, i.ValidatedFormat)
}
