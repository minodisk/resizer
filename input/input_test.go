package input_test

import (
	"reflect"
	"testing"

	"github.com/go-microservices/resizer/input"
)

func TestValidateURL(t *testing.T) {
	type Input struct {
		Input input.Input
		Hosts []string
	}
	type Expected struct {
		Input input.Input
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
				Input: input.Input{
					URL: "http://example.com",
				},
				Hosts: []string{
					"example.com",
				},
			},
			Expected: Expected{
				Input: input.Input{
					URL: "http://example.com",
				},
				Error: nil,
			},
		},
		{
			Spec: "allow https",
			Input: Input{
				Input: input.Input{
					URL: "https://foo.example.com",
				},
				Hosts: []string{
					"example.com",
					"foo.example.com",
				},
			},
			Expected: Expected{
				Input: input.Input{
					URL: "https://foo.example.com",
				},
				Error: nil,
			},
		},
		{
			Spec: "not allow any other scheme",
			Input: Input{
				Input: input.Input{
					URL: "ftp://example.com",
				},
				Hosts: []string{
					"example.com",
				},
			},
			Expected: Expected{
				Input: input.Input{
					URL: "ftp://example.com",
				},
				Error: input.NewInvalidSchemeError("ftp"),
			},
		},
		{
			Spec: "not allow unspecified hosts",
			Input: Input{
				Input: input.Input{
					URL: "http://foo.example.com",
				},
				Hosts: []string{
					"example.com",
				},
			},
			Expected: Expected{
				Input: input.Input{
					URL: "http://foo.example.com",
				},
				Error: input.NewInvalidHostError("foo.example.com"),
			},
		},
		{
			Spec: "Can specify multi-hosts",
			Input: Input{
				Input: input.Input{
					URL: "http://foo.example.com",
				},
				Hosts: []string{
					"example.com",
					"foo.example.com",
				},
			},
			Expected: Expected{
				Input: input.Input{
					URL: "http://foo.example.com",
				},
				Error: nil,
			},
		},
		{
			Spec: "allow port 80",
			Input: Input{
				Input: input.Input{
					URL: "http://example.com:80",
				},
				Hosts: []string{
					"example.com",
				},
			},
			Expected: Expected{
				Input: input.Input{
					URL: "http://example.com:80",
				},
				Error: nil,
			},
		},
		{
			Spec: "not allow any other port",
			Input: Input{
				Input: input.Input{
					URL: "http://example.com:8080",
				},
				Hosts: []string{
					"example.com",
				},
			},
			Expected: Expected{
				Input: input.Input{
					URL: "http://example.com:8080",
				},
				Error: input.NewInvalidHostError("example.com:8080"),
			},
		},
		{
			Spec: "Can specify host with port",
			Input: Input{
				Input: input.Input{
					URL: "http://example.com:8080",
				},
				Hosts: []string{
					"example.com:8080",
				},
			},
			Expected: Expected{
				Input: input.Input{
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
		Input input.Input
	}
	type Expected struct {
		Input input.Input
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
				Input: input.Input{
					Width:  100,
					Height: 200,
				},
			},
			Expected: Expected{
				Input: input.Input{
					Width:  100,
					Height: 200,
				},
				Error: nil,
			},
		},
		{
			Spec: "allow zero width and positive height",
			Input: Input{
				Input: input.Input{
					Width:  0,
					Height: 200,
				},
			},
			Expected: Expected{
				Input: input.Input{
					Width:  0,
					Height: 200,
				},
				Error: nil,
			},
		},
		{
			Spec: "allow positive width and zero height",
			Input: Input{
				Input: input.Input{
					Width:  100,
					Height: 0,
				},
			},
			Expected: Expected{
				Input: input.Input{
					Width:  100,
					Height: 0,
				},
				Error: nil,
			},
		},
		{
			Spec: "not allow negative width",
			Input: Input{
				Input: input.Input{
					Width:  -100,
					Height: 200,
				},
			},
			Expected: Expected{
				Input: input.Input{
					Width:  -100,
					Height: 200,
				},
				Error: input.NewInvalidSizeError(-100, 200),
			},
		},
		{
			Spec: "not allow negative height",
			Input: Input{
				Input: input.Input{
					Width:  100,
					Height: -200,
				},
			},
			Expected: Expected{
				Input: input.Input{
					Width:  100,
					Height: -200,
				},
				Error: input.NewInvalidSizeError(100, -200),
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
		Input input.Input
	}
	type Expected struct {
		Input input.Input
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
				Input: input.Input{
					Method: input.MethodNormal,
				},
			},
			Expected: Expected{
				Input: input.Input{
					Method: input.MethodNormal,
				},
				Error: nil,
			},
		},
		{
			Spec: "allow thumbnail method",
			Input: Input{
				Input: input.Input{
					Method: input.MethodThumbnail,
				},
			},
			Expected: Expected{
				Input: input.Input{
					Method: input.MethodThumbnail,
				},
				Error: nil,
			},
		},
		{
			Spec: "not allow any other method",
			Input: Input{
				Input: input.Input{
					Method: "foo",
				},
			},
			Expected: Expected{
				Input: input.Input{
					Method: "foo",
				},
				Error: input.NewInvalidMethodError("foo"),
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
		Input input.Input
	}
	type Expected struct {
		Input input.Input
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
				Input: input.Input{
					Format: "",
				},
			},
			Expected: Expected{
				Input: input.Input{
					Format: input.FormatJPEG,
				},
				Error: nil,
			},
		},
		{
			Spec: "allow jpeg format",
			Input: Input{
				Input: input.Input{
					Format: input.FormatJPEG,
				},
			},
			Expected: Expected{
				Input: input.Input{
					Format: input.FormatJPEG,
				},
				Error: nil,
			},
		},
		{
			Spec: "allow png format",
			Input: Input{
				Input: input.Input{
					Format: input.FormatPNG,
				},
			},
			Expected: Expected{
				Input: input.Input{
					Format: input.FormatPNG,
				},
				Error: nil,
			},
		},
		{
			Spec: "allow gif format",
			Input: Input{
				Input: input.Input{
					Format: input.FormatGIF,
				},
			},
			Expected: Expected{
				Input: input.Input{
					Format: input.FormatGIF,
				},
				Error: nil,
			},
		},
		{
			Spec: "not allow any other format",
			Input: Input{
				Input: input.Input{
					Format: "foo",
				},
			},
			Expected: Expected{
				Input: input.Input{
					Format: "foo",
				},
				Error: input.NewInvalidFormatError("foo"),
			},
		},
		{
			Spec: "not allow negative quality",
			Input: Input{
				Input: input.Input{
					Format:  input.FormatJPEG,
					Quality: -1,
				},
			},
			Expected: Expected{
				Input: input.Input{
					Format:  input.FormatJPEG,
					Quality: -1,
				},
				Error: input.NewInvalidQualityError(-1),
			},
		},
		{
			Spec: "not allow over 100 quality",
			Input: Input{
				Input: input.Input{
					Format:  input.FormatJPEG,
					Quality: 101,
				},
			},
			Expected: Expected{
				Input: input.Input{
					Format:  input.FormatJPEG,
					Quality: 101,
				},
				Error: input.NewInvalidQualityError(101),
			},
		},
		{
			Spec: "fill quality as 0 with any other format",
			Input: Input{
				Input: input.Input{
					Format:  input.FormatPNG,
					Quality: 80,
				},
			},
			Expected: Expected{
				Input: input.Input{
					Format:  input.FormatPNG,
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
