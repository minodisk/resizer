package testutil

import (
	"io/ioutil"
	"os"
)

const (
	GoogleAuthFilename = "google-auth.json"
)

func CreateGoogleAuthFile() error {
	return ioutil.WriteFile(GoogleAuthFilename, []byte(os.Getenv("GOOGLE_AUTH_JSON")), 0664)
}

func RemoveGoogleAuthFile() error {
	return os.RemoveAll(GoogleAuthFilename)
}
