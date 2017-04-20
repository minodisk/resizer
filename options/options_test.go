package options_test

import (
	"reflect"
	"testing"

	"github.com/minodisk/resizer/options"
)

func TestParse(t *testing.T) {
	t.Parallel()

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		for _, c := range []struct {
			name string
			args []string
		}{
			{
				"undefined flag",
				[]string{
					"-xxx",
				},
			},
		} {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()
				_, err := options.Parse(c.args)
				if err == nil {
					t.Errorf("should return error")
				}
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		for _, c := range []struct {
			name string
			args []string
			want options.Options
		}{} {
			c := c
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()
				got, err := options.Parse(c.args)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, c.want) {
					t.Errorf("\n got: %+v\nwant: %+v", got, c.want)
				}
			})
		}
	})
}
