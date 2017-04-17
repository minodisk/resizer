package options

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type ServiceAccount struct {
	Path        string `json:"-"`
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

func (a *ServiceAccount) String() string {
	return fmt.Sprintf("%+v", *a)
}

func (a *ServiceAccount) Set(path string) error {
	if path == "" {
		return errors.New("path to Google service account JSON isn't specified")
	}
	var b []byte
	var err error
	b, err = ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "fail to read the file of Google service account JSON")
	}
	if err := json.Unmarshal(b, a); err != nil {
		return errors.Wrap(err, "fail to unmarshal JSON")
	}
	a.Path = path
	return nil
}

func (a *ServiceAccount) UnmarshalJSON(data []byte) error {
	type Alias ServiceAccount
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*a = ServiceAccount(alias)
	a.PrivateKey = strings.Replace(a.PrivateKey, `\n`, "\n", -1)
	return nil
}
