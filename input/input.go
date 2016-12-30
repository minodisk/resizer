package input

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	KeyURL     = "url"
	KeyMethod  = "method"
	KeyWidth   = "width"
	KeyHeight  = "height"
	KeyFormat  = "format"
	KeyQuality = "quality"

	MethodNormal    = "normal"
	MethodThumbnail = "thumbnail"
	MethodDefault   = MethodNormal

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

func (i Input) Validate(hosts []string) (Input, error) {
	var err error
	i, err = i.ValidateURL(hosts)
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

func (i Input) ValidateURL(allowedHosts []string) (Input, error) {
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
	var hosts []string
	for _, h := range allowedHosts {
		hosts = append(hosts, h)
		if strings.Index(h, ":") == -1 {
			hosts = append(hosts, fmt.Sprintf("%s:80", h))
		}
	}
	if !in(u.Host, hosts) {
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
	case MethodDefault, MethodThumbnail:
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
