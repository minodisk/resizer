package options_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/minodisk/resizer/options"
)

func TestOptions(t *testing.T) {
	for _, c := range []struct {
		name string
		envs map[string]string
		args []string
		want *options.Options
	}{
		{
			"multiple hosts with comma separated",
			map[string]string{},
			[]string{
				"-host", "a.com,b.com",
			},
			&options.Options{
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
			map[string]string{},
			[]string{
				"-host", "a.com",
				"-host", "b.com",
			},
			&options.Options{
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
			map[string]string{},
			[]string{
				"-host", "a.com,b.com",
				"-host", "c.com",
			},
			&options.Options{
				AllowedHosts: []string{
					"a.com",
					"b.com",
					"c.com",
				},
				ObjectPrefix: "resized/",
				Port:         80,
			},
		},
		{
			"only env",
			map[string]string{
				options.EnvBucket: "foo",
			},
			[]string{},
			&options.Options{
				Bucket:       "foo",
				ObjectPrefix: "resized/",
				Port:         80,
			},
		},
		{
			"only args",
			map[string]string{},
			[]string{
				"-bucket", "bar",
			},
			&options.Options{
				Bucket:       "bar",
				ObjectPrefix: "resized/",
				Port:         80,
			},
		},
		{
			"envs and args",
			map[string]string{
				options.EnvBucket: "foo",
			},
			[]string{
				"-bucket", "bar",
			},
			&options.Options{
				Bucket:       "bar",
				ObjectPrefix: "resized/",
				Port:         80,
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			got := &options.Options{}

			for _, k := range options.Envs {
				os.Setenv(k, "")
			}
			for k, v := range c.envs {
				os.Setenv(k, v)
			}
			if err := got.Parse(c.args); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Error("ENVS:")
				for _, k := range options.Envs {
					t.Errorf("%s: %s\n", k, os.Getenv(k))
				}
				t.Errorf("\ngot:\n%+v\nwant:\n%+v", got, c.want)
			}
		})
	}
}
