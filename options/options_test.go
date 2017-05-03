package options_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/minodisk/resizer/options"
	"github.com/minodisk/resizer/testutil"
)

func TestMain(m *testing.M) {
	if err := testutil.CreateGoogleAuthFile(); err != nil {
		panic(err)
	}
	code := m.Run()
	if err := testutil.RemoveGoogleAuthFile(); err != nil {
		panic(err)
	}
	os.Exit(code)
}

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
		}{
			{
				"default",
				[]string{},
				options.Options{
					AllowedHosts:       options.Hosts{},
					Bucket:             "",
					DataSourceName:     "",
					MaxHTTPConnections: 0,
					ObjectPrefix:       "",
					Port:               80,
					ServiceAccountFile: "",
					Verbose:            false,
				},
			},
			{
				"overwrite",
				[]string{
					"-host", "foo",
					"-bucket", "bar",
					"-dsn", "baz",
					"-connections", "16",
					"-prefix", "qux",
					"-port", "9000",
					"-account", "google-auth.json",
					"-verbose", "true",
				},
				options.Options{
					AllowedHosts: options.Hosts{
						"foo",
					},
					Bucket:             "bar",
					DataSourceName:     "baz",
					MaxHTTPConnections: 16,
					ObjectPrefix:       "qux",
					Port:               9000,
					ServiceAccountFile: "google-auth.json",
					Verbose:            true,
				},
			},
		} {
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
