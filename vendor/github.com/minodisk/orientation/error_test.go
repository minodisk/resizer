package orientation_test

import (
	"errors"
	"testing"

	"github.com/minodisk/orientation"
)

type causer interface {
	Cause() error
}

func TestError(t *testing.T) {
	want := errors.New("want")
	for _, c := range []struct {
		name string
		err  error
	}{
		{
			"DecodeError",
			&orientation.DecodeError{want},
		},
		{
			"FormatError",
			&orientation.FormatError{want},
		},
		{
			"TagError",
			&orientation.TagError{want},
		},
		{
			"OrientError",
			&orientation.OrientError{want},
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			err, ok := c.err.(causer)
			if !ok {
				t.Fatal("%s does not implements causer", c.name)
			}
			got := err.Cause()
			t.Run("Cause()", func(t *testing.T) {
				if got != want {
					t.Errorf("got %v, want %v", got, want)
				}
			})
			t.Run("Error()", func(t *testing.T) {
				if got.Error() != want.Error() {
					t.Errorf("got %s, want %s", got.Error(), want.Error())
				}
			})
		})
	}
}
