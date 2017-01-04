package option_test

import (
	"reflect"
	"testing"

	"github.com/go-microservices/resizer/option"
	"github.com/pkg/errors"
)

func TestLoad(t *testing.T) {
	os, err := option.Load("../fixtures/config.yml")
	if err != nil {
		t.Fatal(errors.Wrap(err, "fail to load config"))
	}

	if a, e := len(os), 2; a != e {
		t.Errorf("the length of options expected %d, but actual %d", e, a)
	}

	type Case struct {
		Key    string
		Option option.NewOption
	}

	cases := []Case{
		{
			Key: "*",
			Option: option.NewOption{
				MaxHTTPConnections: 7,
				GoogleCloud: option.GoogleCloud{
					ProjectID:         "syoya-test",
					ServiceAccount:    "/secret/gcloud.json",
					StorageBucketName: "resizer",
				},
				MySQL: option.MySQL{
					DataSourceName: "root:@tcp(mysql:3306)/resizer?charset=utf8&parseTime=True",
				},
			},
		},
		{
			Key: "foo.bar",
			Option: option.NewOption{
				MaxHTTPConnections: 7,
				GoogleCloud: option.GoogleCloud{
					ProjectID:         "foo-bar",
					ServiceAccount:    "/foo/bar.json",
					StorageBucketName: "foo",
				},
				MySQL: option.MySQL{
					DataSourceName: "foo:bar@tcp(mysql:3306)/baz",
				},
			},
		},
	}

	for _, c := range cases {
		o, ok := os[c.Key]
		if !ok {
			t.Errorf("the key `%s` expected existing, but actual not")
			continue
		}
		if a, e := o, c.Option; !reflect.DeepEqual(a, e) {
			t.Errorf("option is expected:\n%+v\nbut actual:\n%+v\n", e, a)
		}
	}
}
