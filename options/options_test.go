package options_test

import (
	"reflect"
	"testing"

	"github.com/minodisk/resizer/options"
)

func TestOptions(t *testing.T) {
	for _, c := range []struct {
		name string
		args []string
		want options.Options
	}{
		{
			"multiple hosts with comma separated",
			[]string{
				"-host", "a.com,b.com",
			},
			options.Options{
				AllowedHosts: []string{
					"a.com",
					"b.com",
				},
				ObjectPrefix: "resized/",
				Port:         80,
			},
		},
		{
			"multiple hosts with specified multiple times",
			[]string{
				"-host", "a.com",
				"-host", "b.com",
			},
			options.Options{
				AllowedHosts: []string{
					"a.com",
					"b.com",
				},
				ObjectPrefix: "resized/",
				Port:         80,
			},
		},
		{
			"multiple hosts with both way",
			[]string{
				"-host", "a.com,b.com",
				"-host", "c.com",
			},
			options.Options{
				AllowedHosts: []string{
					"a.com",
					"b.com",
					"c.com",
				},
				ObjectPrefix: "resized/",
				Port:         80,
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := options.Parse(c.args)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("got %+v, want %+v", got, c.want)
			}
		})
	}
}
