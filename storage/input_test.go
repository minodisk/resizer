package storage_test

import (
	"reflect"
	"testing"

	"github.com/go-microservices/resizer/storage"
)

func TestValidateURL(t *testing.T) {
	type Input struct {
		Input storage.Input
		Hosts []string
	}
	type Expected struct {
		Input storage.Input
		Error error
	}
	type Case struct {
		Spec     string
		Input    Input
		Expected Expected
	}

	cases := []Case{
		{
			Spec: "allow http",
			Input: Input{
				Input: storage.Input{
					URL: "http://example.com",
				},
				Hosts: []string{
					"example.com",
				},
			},
			Expected: Expected{
				Input: storage.Input{
					URL: "http://example.com",
				},
				Error: nil,
			},
		},
		{
			Spec: "allow https",
			Input: Input{
				Input: storage.Input{
					URL: "https://foo.example.com",
				},
				Hosts: []string{
					"example.com",
					"foo.example.com",
				},
			},
			Expected: Expected{
				Input: storage.Input{
					URL: "https://foo.example.com",
				},
				Error: nil,
			},
		},
		{
			Spec: "not allow any other scheme",
			Input: Input{
				Input: storage.Input{
					URL: "ftp://example.com",
				},
				Hosts: []string{
					"example.com",
				},
			},
			Expected: Expected{
				Input: storage.Input{
					URL: "ftp://example.com",
				},
				Error: storage.NewInvalidSchemeError("ftp"),
			},
		},
		{
			Spec: "not allow unspecified hosts",
			Input: Input{
				Input: storage.Input{
					URL: "http://foo.example.com",
				},
				Hosts: []string{
					"example.com",
				},
			},
			Expected: Expected{
				Input: storage.Input{
					URL: "http://foo.example.com",
				},
				Error: storage.NewInvalidHostError("foo.example.com"),
			},
		},
		{
			Spec: "Can specify multi-hosts",
			Input: Input{
				Input: storage.Input{
					URL: "http://foo.example.com",
				},
				Hosts: []string{
					"example.com",
					"foo.example.com",
				},
			},
			Expected: Expected{
				Input: storage.Input{
					URL: "http://foo.example.com",
				},
				Error: nil,
			},
		},
		{
			Spec: "allow port 80",
			Input: Input{
				Input: storage.Input{
					URL: "http://example.com:80",
				},
				Hosts: []string{
					"example.com",
				},
			},
			Expected: Expected{
				Input: storage.Input{
					URL: "http://example.com:80",
				},
				Error: nil,
			},
		},
		{
			Spec: "not allow any other port",
			Input: Input{
				Input: storage.Input{
					URL: "http://example.com:8080",
				},
				Hosts: []string{
					"example.com",
				},
			},
			Expected: Expected{
				Input: storage.Input{
					URL: "http://example.com:8080",
				},
				Error: storage.NewInvalidHostError("example.com:8080"),
			},
		},
		{
			Spec: "Can specify host with port",
			Input: Input{
				Input: storage.Input{
					URL: "http://example.com:8080",
				},
				Hosts: []string{
					"example.com:8080",
				},
			},
			Expected: Expected{
				Input: storage.Input{
					URL: "http://example.com:8080",
				},
				Error: nil,
			},
		},
	}

	for _, c := range cases {
		output, err := c.Input.Input.ValidateURL(c.Input.Hosts)
		if !reflect.DeepEqual(output, c.Expected.Input) {
			t.Errorf("ValidateURL() should %s.\nOutput expected `%v`, but actual `%v`", c.Spec, c.Expected.Input, output)
		}
		if !reflect.DeepEqual(err, c.Expected.Error) {
			t.Errorf("ValidateURL() should %s.\nError expected `%v`, but actual `%v`", c.Spec, c.Expected.Error, err)
		}
	}
}

func TestValidateSize(t *testing.T) {
	type Input struct {
		Input storage.Input
	}
	type Expected struct {
		Input storage.Input
		Error error
	}
	type Case struct {
		Spec     string
		Input    Input
		Expected Expected
	}

	cases := []Case{
		{
			Spec: "allow positive size",
			Input: Input{
				Input: storage.Input{
					Width:  100,
					Height: 200,
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Width:  100,
					Height: 200,
				},
				Error: nil,
			},
		},
		{
			Spec: "allow zero width and positive height",
			Input: Input{
				Input: storage.Input{
					Width:  0,
					Height: 200,
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Width:  0,
					Height: 200,
				},
				Error: nil,
			},
		},
		{
			Spec: "allow positive width and zero height",
			Input: Input{
				Input: storage.Input{
					Width:  100,
					Height: 0,
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Width:  100,
					Height: 0,
				},
				Error: nil,
			},
		},
		{
			Spec: "not allow negative width",
			Input: Input{
				Input: storage.Input{
					Width:  -100,
					Height: 200,
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Width:  -100,
					Height: 200,
				},
				Error: storage.NewInvalidSizeError(-100, 200),
			},
		},
		{
			Spec: "not allow negative height",
			Input: Input{
				Input: storage.Input{
					Width:  100,
					Height: -200,
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Width:  100,
					Height: -200,
				},
				Error: storage.NewInvalidSizeError(100, -200),
			},
		},
	}

	for _, c := range cases {
		output, err := c.Input.Input.ValidateSize()
		if !reflect.DeepEqual(output, c.Expected.Input) {
			t.Errorf("ValidateSize() should %s.\nOutput expected `%v`, but actual `%v`", c.Spec, c.Expected.Input, output)
		}
		if !reflect.DeepEqual(err, c.Expected.Error) {
			t.Errorf("ValidateSize() should %s.\nError expected `%v`, but actual `%v`", c.Spec, c.Expected.Error, err)
		}
	}
}

func TestValidateMethod(t *testing.T) {
	type Input struct {
		Input storage.Input
	}
	type Expected struct {
		Input storage.Input
		Error error
	}
	type Case struct {
		Spec     string
		Input    Input
		Expected Expected
	}

	cases := []Case{
		{
			Spec: "allow normal method",
			Input: Input{
				Input: storage.Input{
					Method: storage.MethodNormal,
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Method: storage.MethodNormal,
				},
				Error: nil,
			},
		},
		{
			Spec: "allow thumbnail method",
			Input: Input{
				Input: storage.Input{
					Method: storage.MethodThumbnail,
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Method: storage.MethodThumbnail,
				},
				Error: nil,
			},
		},
		{
			Spec: "not allow any other method",
			Input: Input{
				Input: storage.Input{
					Method: "foo",
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Method: "foo",
				},
				Error: storage.NewInvalidMethodError("foo"),
			},
		},
	}

	for _, c := range cases {
		output, err := c.Input.Input.ValidateMethod()
		if !reflect.DeepEqual(output, c.Expected.Input) {
			t.Errorf("ValidateMethod() should %s.\nOutput expected `%v`, but actual `%v`", c.Spec, c.Expected.Input, output)
		}
		if !reflect.DeepEqual(err, c.Expected.Error) {
			t.Errorf("ValidateMethod() should %s.\nError expected `%v`, but actual `%v`", c.Spec, c.Expected.Error, err)
		}
	}
}

func TestValidateFormatAndQuality(t *testing.T) {
	type Input struct {
		Input storage.Input
	}
	type Expected struct {
		Input storage.Input
		Error error
	}
	type Case struct {
		Spec     string
		Input    Input
		Expected Expected
	}

	cases := []Case{
		{
			Spec: "fill empty format with jpeg",
			Input: Input{
				Input: storage.Input{
					Format: "",
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Format: storage.FormatJpeg,
				},
				Error: nil,
			},
		},
		{
			Spec: "allow jpeg format",
			Input: Input{
				Input: storage.Input{
					Format: storage.FormatJpeg,
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Format: storage.FormatJpeg,
				},
				Error: nil,
			},
		},
		{
			Spec: "allow png format",
			Input: Input{
				Input: storage.Input{
					Format: storage.FormatPng,
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Format: storage.FormatPng,
				},
				Error: nil,
			},
		},
		{
			Spec: "allow gif format",
			Input: Input{
				Input: storage.Input{
					Format: storage.FormatGif,
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Format: storage.FormatGif,
				},
				Error: nil,
			},
		},
		{
			Spec: "not allow any other format",
			Input: Input{
				Input: storage.Input{
					Format: "foo",
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Format: "foo",
				},
				Error: storage.NewInvalidFormatError("foo"),
			},
		},
		{
			Spec: "not allow negative quality",
			Input: Input{
				Input: storage.Input{
					Format:  storage.FormatJpeg,
					Quality: -1,
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Format:  storage.FormatJpeg,
					Quality: -1,
				},
				Error: storage.NewInvalidQualityError(-1),
			},
		},
		{
			Spec: "not allow over 100 quality",
			Input: Input{
				Input: storage.Input{
					Format:  storage.FormatJpeg,
					Quality: 101,
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Format:  storage.FormatJpeg,
					Quality: 101,
				},
				Error: storage.NewInvalidQualityError(101),
			},
		},
		{
			Spec: "fill quality as 0 with any other format",
			Input: Input{
				Input: storage.Input{
					Format:  storage.FormatPng,
					Quality: 80,
				},
			},
			Expected: Expected{
				Input: storage.Input{
					Format:  storage.FormatPng,
					Quality: 0,
				},
				Error: nil,
			},
		},
	}

	for _, c := range cases {
		output, err := c.Input.Input.ValidateFormatAndQuality()
		if !reflect.DeepEqual(output, c.Expected.Input) {
			t.Errorf("ValidateFormatAndQuality() should %s.\nOutput expected `%v`, but actual `%v`", c.Spec, c.Expected.Input, output)
		}
		if !reflect.DeepEqual(err, c.Expected.Error) {
			t.Errorf("ValidateFormatAndQuality() should %s.\nError expected `%v`, but actual `%v`", c.Spec, c.Expected.Error, err)
		}
	}
}
