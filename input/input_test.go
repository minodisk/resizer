package input_test

import (
	"reflect"
	"testing"

	"github.com/minodisk/resizer/input"
)

func TestValidateURL(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		name  string
		input input.Input
		hosts []string
		want  input.Input
		err   error
	}{
		{
			"allow http",
			input.Input{
				URL: "http://example.com",
			},
			[]string{
				"example.com",
			},
			input.Input{
				URL: "http://example.com",
			},
			nil,
		},
		{
			"allow https",
			input.Input{
				URL: "https://foo.example.com",
			},
			[]string{
				"example.com",
				"foo.example.com",
			},
			input.Input{
				URL: "https://foo.example.com",
			},
			nil,
		},
		{
			"not allow any other scheme",
			input.Input{
				URL: "ftp://example.com",
			},
			[]string{
				"example.com",
			},
			input.Input{
				URL: "ftp://example.com",
			},
			input.NewInvalidSchemeError("ftp"),
		},
		{
			"not allow unspecified hosts",
			input.Input{
				URL: "http://foo.example.com",
			},
			[]string{
				"example.com",
			},
			input.Input{
				URL: "http://foo.example.com",
			},
			input.NewInvalidHostError("foo.example.com"),
		},
		{
			"Can specify multi-hosts",
			input.Input{
				URL: "http://foo.example.com",
			},
			[]string{
				"example.com",
				"foo.example.com",
			},
			input.Input{
				URL: "http://foo.example.com",
			},
			nil,
		},
		{
			"not allow any other port",
			input.Input{
				URL: "http://example.com:8080",
			},
			[]string{
				"example.com",
			},
			input.Input{
				URL: "http://example.com:8080",
			},
			input.NewInvalidHostError("example.com:8080"),
		},
		{
			"Can specify host with port",
			input.Input{
				URL: "http://example.com:8080",
			},
			[]string{
				"example.com:8080",
			},
			input.Input{
				URL: "http://example.com:8080",
			},
			nil,
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got, err := c.input.ValidateURL(c.hosts)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("result\n got: %+v\nwant: %+v", got, c.want)
			}
			if !reflect.DeepEqual(err, c.err) {
				t.Errorf("error\n got: %+v\nwant: %+v", err, c.err)
			}
		})
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
