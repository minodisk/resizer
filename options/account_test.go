package options_test

import (
	"testing"

	"github.com/minodisk/resizer/options"
)

func TestAccountString(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name    string
		account options.ServiceAccount
		want    string
	}{
		{
			"general",
			options.ServiceAccount{
				Path:        "foo.json",
				ClientEmail: "bar@example.com",
				PrivateKey:  "XXXXXXXXX\nYYYYYYYYYY",
				ProjectID:   "baz",
			},
			"bar@example.com",
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			got := c.account.String()
			if got != c.want {
				t.Errorf("got: %s, want: %s", got, c.want)
			}
		})
	}
}
