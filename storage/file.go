package storage

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"image"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	MethodNormal    = "normal"
	MethodThumbnail = "thumbnail"
	MethodDefault   = MethodNormal

	FormatJpeg    = "jpeg"
	FormatPng     = "png"
	FormatGif     = "gif"
	FormatDefault = FormatJpeg

	QualityMax     = 100
	QualityMin     = 0
	QualityDefault = QualityMin
)

var (
	schemes = []string{
		"http",
		"https",
	}
)

type Image struct {
	ID               uint64
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ValidatedURL     string `sql:"type:text"`
	ValidatedMethod  string
	ValidatedFormat  string
	ValidatedQuality uint8
	ValidatedWidth   int
	ValidatedHeight  int
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
func NewImage(q map[string][]string, hosts []string) (Image, error) {
	i := Image{}
	if len(q["url"]) != 0 {
		i.ValidatedURL = q["url"][0]
	}
	if len(q["method"]) != 0 {
		i.ValidatedMethod = q["method"][0]
	}
	if len(q["width"]) != 0 {
		w, err := strconv.Atoi(q["width"][0])
		if err != nil {
			return i, err
		}
		i.ValidatedWidth = w
	}
	if len(q["height"]) != 0 {
		h, err := strconv.Atoi(q["height"][0])
		if err != nil {
			return i, err
		}
		i.ValidatedHeight = h
	}
	if len(q["format"]) != 0 {
		i.ValidatedFormat = q["format"][0]
	}
	if len(q["quality"]) != 0 {
		q, err := strconv.ParseUint(q["quality"][0], 10, 8)
		if err != nil {
			return i, err
		}
		i.ValidatedQuality = uint8(q)
	}
	return i.validate(hosts)
}

// validate はパラメーターに正しい値が入っているかを検査します。
// 間違った値が入っている場合はエラーを返す。
func (i Image) validate(hosts []string) (Image, error) {
	var err error
	i, err = i.v(hosts)
	if err != nil {
		return i, err
	}
	return i.serializeValidatedProps()
}

func (i Image) v(hosts []string) (Image, error) {
	if i.ValidatedURL == "" {
		return i, fmt.Errorf("url shouldn't be empty")
	}
	u, err := url.Parse(i.ValidatedURL)
	if err != nil {
		return i, err
	}
	isValidScheme := func() bool {
		for _, scheme := range schemes {
			if scheme == u.Scheme {
				return true
			}
		}
		return false
	}()
	if !isValidScheme {
		return i, fmt.Errorf("the scheme of url is allowed %s", strings.Join(schemes, " or "))
	}
	host, isValidHost := func() (string, bool) {
		host := strings.Split(u.Host, ":")[0]
		for _, h := range hosts {
			index := strings.LastIndex(host, h)
			if index != -1 && index == len(host)-len(h) {
				return host, true
			} else if host == "localhost" {
				return host, true
			}
		}
		r := regexp.MustCompile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}")
		if r.MatchString(host) {
			return host, true
		}
		return host, false
	}()
	if !isValidHost {
		return i, fmt.Errorf("the host '%s' isn't allowed", host)
	}

	if i.ValidatedWidth < 0 || i.ValidatedHeight < 0 {
		return i, fmt.Errorf("size should be specified with positive value: width=%d, height=%d", i.ValidatedWidth, i.ValidatedHeight)
	}
	if i.ValidatedWidth == 0 && i.ValidatedHeight == 0 {
		return i, fmt.Errorf("size should not be zero: width=%d, height=%d", i.ValidatedWidth, i.ValidatedHeight)
	}

	if i.ValidatedMethod == "" {
		i.ValidatedMethod = MethodDefault
	} else if !(i.ValidatedMethod == MethodDefault || i.ValidatedMethod == MethodThumbnail) {
		return i, fmt.Errorf("method '%s' isn't allowed, specify with %s or %s", i.ValidatedMethod, MethodDefault, MethodThumbnail)
	}

	if i.ValidatedFormat == "" {
		i.ValidatedFormat = FormatDefault
	} else if !(i.ValidatedFormat == FormatJpeg || i.ValidatedFormat == FormatPng || i.ValidatedFormat == FormatGif) {
		return i, fmt.Errorf("format '%s' isn't allowed, specify with %s, %s or %s", i.ValidatedFormat, FormatJpeg, FormatPng, FormatGif)
	}

	if i.ValidatedFormat != FormatJpeg {
		i.ValidatedQuality = QualityDefault
	} else if i.ValidatedQuality == QualityDefault {
		i.ValidatedQuality = 100
	} else if i.ValidatedQuality < QualityMin || QualityMax < i.ValidatedQuality {
		return i, fmt.Errorf("quality %d should be specify between %d and %d", i.ValidatedQuality, QualityMin, QualityMax)
	}

	return i, nil
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
	case MethodNormal:
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
	case MethodThumbnail:
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
	uuid := uuid.NewV4().String()
	return fmt.Sprintf("%s.%s", uuid, i.ValidatedFormat)
}
