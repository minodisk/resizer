package input

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/minodisk/resizer/options"
)

const (
	KeyURL     = "url"
	KeyMethod  = "method"
	KeyWidth   = "width"
	KeyHeight  = "height"
	KeyFormat  = "format"
	KeyQuality = "quality"

	MethodContain = "contain"
	MethodCover   = "cover"
	MethodDefault = MethodContain

	FormatJPEG    = "jpeg"
	FormatPNG     = "png"
	FormatGIF     = "gif"
	FormatDefault = FormatJPEG

	QualityMax     = 100
	QualityMin     = 0
	QualityDefault = QualityMin
)

var (
	allowedSchemes = []string{
		"http",
		"https",
	}
)

type Input struct {
	URL     string
	Method  string
	Width   int
	Height  int
	Format  string
	Quality int
}

func New(q map[string][]string) (Input, error) {
	var o Input
	if len(q[KeyURL]) != 0 {
		o.URL = q[KeyURL][0]
	}
	if len(q[KeyMethod]) != 0 {
		o.Method = q[KeyMethod][0]
	}
	if len(q[KeyWidth]) != 0 {
		w, err := strconv.Atoi(q[KeyWidth][0])
		if err != nil {
			return o, err
		}
		o.Width = w
	}
	if len(q[KeyHeight]) != 0 {
		h, err := strconv.Atoi(q[KeyHeight][0])
		if err != nil {
			return o, err
		}
		o.Height = h
	}
	if len(q[KeyFormat]) != 0 {
		o.Format = q[KeyFormat][0]
	}
	if len(q[KeyQuality]) != 0 {
		var err error
		o.Quality, err = strconv.Atoi(q[KeyQuality][0])
		if err != nil {
			return o, err
		}
	} else {
		o.Quality = 100
	}
	return o, nil
}

func (i Input) Validate(allowedHosts options.Hosts) (Input, error) {
	var err error
	i, err = i.ValidateURL(allowedHosts)
	if err != nil {
		return i, err
	}
	i, err = i.ValidateSize()
	if err != nil {
		return i, err
	}
	i, err = i.ValidateMethod()
	if err != nil {
		return i, err
	}
	i, err = i.ValidateFormatAndQuality()
	if err != nil {
		return i, err
	}
	return i, nil
}

func (i Input) ValidateURL(allowedHosts options.Hosts) (Input, error) {
	if i.URL == "" {
		return i, fmt.Errorf("URL shouldn't be empty")
	}
	u, err := url.Parse(i.URL)
	if err != nil {
		return i, err
	}
	if !in(u.Scheme, allowedSchemes) {
		return i, NewInvalidSchemeError(u.Scheme)
	}
	if !allowedHosts.Contains(u.Host) {
		return i, NewInvalidHostError(u.Host)
	}
	return i, nil
}

func (i Input) ValidateSize() (Input, error) {
	if i.Width < 0 || i.Height < 0 || (i.Width == 0 && i.Height == 0) {
		return i, NewInvalidSizeError(i.Width, i.Height)
	}
	return i, nil
}

func (i Input) ValidateMethod() (Input, error) {
	switch i.Method {
	case "":
		i.Method = MethodDefault
	case MethodDefault, MethodCover:
	default:
		return i, NewInvalidMethodError(i.Method)
	}
	return i, nil
}

func (i Input) ValidateFormatAndQuality() (Input, error) {
	switch i.Format {
	case "":
		i.Format = FormatDefault
	case FormatJPEG, FormatPNG, FormatGIF:
	default:
		return i, NewInvalidFormatError(i.Format)
	}
	switch i.Format {
	case FormatJPEG:
		if i.Quality < QualityMin || QualityMax < i.Quality {
			return i, NewInvalidQualityError(i.Quality)
		}
	default:
		i.Quality = QualityDefault
	}
	return i, nil
}
