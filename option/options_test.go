package option_test

import (
	"os"
	"strings"
	"testing"

	"github.com/go-microservices/resizer/option"
)

func TestFlags(t *testing.T) {
	o, err := option.New([]string{
		"--id", "AAAA",
		"--bucket", "BBBB",
		"--json", "CCCC",
		"--dbuser", "DDDD",
		"--dbpassword", "EEEE",
		"--dbprotocol", "FFFF",
		"--dbaddress", "GGGG",
		"--dbname", "HHHH",
		"--host", "IIII",
		"--host", "JJJJ",
	})
	if err != nil {
		t.Fatal(err)
	}
	if o.GCProjectID != "AAAA" {
		t.Error("wrong ProjectID:", o.GCProjectID)
	}
	if o.GCStorageBucket != "BBBB" {
		t.Error("wrong Bucket:", o.GCStorageBucket)
	}
	if o.JSON != "CCCC" {
		t.Error("wrong JSON:", o.JSON)
	}
	if o.DBUser != "DDDD" {
		t.Error("wrong DBUser:", o.DBUser)
	}
	if o.DBPassword != "EEEE" {
		t.Error("wrong DBPassword:", o.DBPassword)
	}
	if o.DBProtocol != "FFFF" {
		t.Error("wrong DBProtocol:", o.DBProtocol)
	}
	if o.DBAddress != "GGGG" {
		t.Error("wrong DBAddress:", o.DBAddress)
	}
	if o.DBName != "HHHH" {
		t.Error("wrong DBName:", o.DBName)
	}
	if o.AllowedHosts[0] != "IIII" {
		t.Error("wrong Hosts[0]:", o.AllowedHosts[0])
	}
	if o.AllowedHosts[1] != "JJJJ" {
		t.Error("wrong Hosts[1]:", o.AllowedHosts[0])
	}
}

func TestEnvar(t *testing.T) {
	for key, value := range map[string]string{
		"RESIZER_PROJECT_ID":  "AAAA",
		"RESIZER_BUCKET":      "BBBB",
		"RESIZER_JSON":        "CCCC",
		"RESIZER_DB_USER":     "DDDD",
		"RESIZER_DB_PASSWORD": "EEEE",
		"RESIZER_DB_PROTOCOL": "FFFF",
		"RESIZER_DB_ADDRESS":  "GGGG",
		"RESIZER_DB_NAME":     "HHHH",
		"RESIZER_HOSTS":       "IIII,JJJJ",
	} {
		if err := os.Setenv(key, value); err != nil {
			t.Fatal(err)
		}
	}
	o, err := option.New([]string{})
	if err != nil {
		t.Fatal(err)
	}
	if o.GCProjectID != os.Getenv("RESIZER_PROJECT_ID") {
		t.Error("wrong ProjectID:", o.GCProjectID)
	}
	if o.GCStorageBucket != os.Getenv("RESIZER_BUCKET") {
		t.Error("wrong Bucket:", o.GCStorageBucket)
	}
	if o.JSON != os.Getenv("RESIZER_JSON") {
		t.Error("wrong JSON:", o.JSON)
	}
	if o.DBUser != os.Getenv("RESIZER_DB_USER") {
		t.Error("wrong DBUser:", o.DBUser)
	}
	if o.DBPassword != os.Getenv("RESIZER_DB_PASSWORD") {
		t.Error("wrong DBPassword:", o.DBPassword)
	}
	if o.DBProtocol != os.Getenv("RESIZER_DB_PROTOCOL") {
		t.Error("wrong DBProtocol:", o.DBProtocol)
	}
	if o.DBAddress != os.Getenv("RESIZER_DB_ADDRESS") {
		t.Error("wrong DBAddress:", o.DBAddress)
	}
	if o.DBName != os.Getenv("RESIZER_DB_NAME") {
		t.Error("wrong DBName:", o.DBName)
	}
	SplitedHosts := strings.Split(os.Getenv("RESIZER_HOSTS"), ",")
	if o.AllowedHosts[0] != SplitedHosts[0] {
		t.Error("wrong Hosts:", o.AllowedHosts, "expect:", SplitedHosts)
	}
}
