package processor_test

import (
	"image"
	"image/color"
	"reflect"
	"testing"

	"github.com/go-microservices/resizer/processor"
)

type Pixels struct {
	size   image.Point
	matrix []int
}

var (
	colors = map[int]color.Color{
		0: color.Black,
		1: color.White,
	}
)

func NewImage(pixels Pixels) image.Image {
	input := image.NewRGBA(image.Rect(0, 0, pixels.size.X, pixels.size.Y))
	for y := 0; y < pixels.size.Y; y++ {
		for x := 0; x < pixels.size.X; x++ {
			at := pixels.size.X*y + x
			color := colors[pixels.matrix[at]]
			input.Set(x, y, color)
		}
	}
	return input
}

func TestRotate(t *testing.T) {
	type Case struct {
		spec        string
		orientation int
		input       Pixels
		expected    Pixels
	}

	cases := []Case{
		{
			spec:        "Can rotate correctly with orientation tag 1",
			orientation: 1,
			input: Pixels{
				size: image.Pt(5, 6),
				matrix: []int{
					0, 0, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 0, 0, 0, 0,
				},
			},
			expected: Pixels{
				size: image.Pt(5, 6),
				matrix: []int{
					0, 0, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 0, 0, 0, 0,
				},
			},
		},
		{
			spec:        "Can rotate correctly with orientation tag 2",
			orientation: 2,
			input: Pixels{
				size: image.Pt(5, 6),
				matrix: []int{
					0, 0, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 0, 0, 0, 0,
				},
			},
			expected: Pixels{
				size: image.Pt(5, 6),
				matrix: []int{
					0, 0, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 0, 0, 1, 0,
					0, 1, 1, 1, 0,
					0, 0, 0, 1, 0,
					0, 0, 0, 0, 0,
				},
			},
		},
		{
			spec:        "Can rotate correctly with orientation tag 3",
			orientation: 3,
			input: Pixels{
				size: image.Pt(5, 6),
				matrix: []int{
					0, 0, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 0, 0, 0, 0,
				},
			},
			expected: Pixels{
				size: image.Pt(5, 6),
				matrix: []int{
					0, 0, 0, 0, 0,
					0, 0, 0, 1, 0,
					0, 1, 1, 1, 0,
					0, 0, 0, 1, 0,
					0, 1, 1, 1, 0,
					0, 0, 0, 0, 0,
				},
			},
		},
		{
			spec:        "Can rotate correctly with orientation tag 4",
			orientation: 4,
			input: Pixels{
				size: image.Pt(5, 6),
				matrix: []int{
					0, 0, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 0, 0, 0, 0,
				},
			},
			expected: Pixels{
				size: image.Pt(5, 6),
				matrix: []int{
					0, 0, 0, 0, 0,
					0, 1, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 0, 0, 0, 0,
				},
			},
		},
		{
			spec:        "Can rotate correctly with orientation tag 5",
			orientation: 5,
			input: Pixels{
				size: image.Pt(5, 6),
				matrix: []int{
					0, 0, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 0, 0, 0, 0,
				},
			},
			expected: Pixels{
				size: image.Pt(6, 5),
				matrix: []int{
					0, 0, 0, 0, 0, 0,
					0, 1, 1, 1, 1, 0,
					0, 1, 0, 1, 0, 0,
					0, 1, 0, 1, 0, 0,
					0, 0, 0, 0, 0, 0,
				},
			},
		},
		{
			spec:        "Can rotate correctly with orientation tag 6",
			orientation: 6,
			input: Pixels{
				size: image.Pt(5, 6),
				matrix: []int{
					0, 0, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 0, 0, 0, 0,
				},
			},
			expected: Pixels{
				size: image.Pt(6, 5),
				matrix: []int{
					0, 0, 0, 0, 0, 0,
					0, 1, 1, 1, 1, 0,
					0, 0, 1, 0, 1, 0,
					0, 0, 1, 0, 1, 0,
					0, 0, 0, 0, 0, 0,
				},
			},
		},
		{
			spec:        "Can rotate correctly with orientation tag 7",
			orientation: 7,
			input: Pixels{
				size: image.Pt(5, 6),
				matrix: []int{
					0, 0, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 0, 0, 0, 0,
				},
			},
			expected: Pixels{
				size: image.Pt(6, 5),
				matrix: []int{
					0, 0, 0, 0, 0, 0,
					0, 0, 1, 0, 1, 0,
					0, 0, 1, 0, 1, 0,
					0, 1, 1, 1, 1, 0,
					0, 0, 0, 0, 0, 0,
				},
			},
		},
		{
			spec:        "Can rotate correctly with orientation tag 8",
			orientation: 8,
			input: Pixels{
				size: image.Pt(5, 6),
				matrix: []int{
					0, 0, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 1, 1, 1, 0,
					0, 1, 0, 0, 0,
					0, 0, 0, 0, 0,
				},
			},
			expected: Pixels{
				size: image.Pt(6, 5),
				matrix: []int{
					0, 0, 0, 0, 0, 0,
					0, 1, 0, 1, 0, 0,
					0, 1, 0, 1, 0, 0,
					0, 1, 1, 1, 1, 0,
					0, 0, 0, 0, 0, 0,
				},
			},
		},
	}

	for _, c := range cases {
		a, err := processor.Rotate(NewImage(c.input), c.orientation)
		if err != nil {
			t.Errorf("%s. Fail to rotate image with error: %s", err)
			continue
		}
		if e := NewImage(c.expected); !reflect.DeepEqual(a, e) {
			t.Errorf("%s. Rotated image is expected:\n%+v\nbut actual:\n%+v\n", c.spec, e, a)
		}
	}
}
