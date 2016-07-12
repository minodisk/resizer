package uploader_test

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/go-microservices/resizer/storage"
	"github.com/go-microservices/resizer/uploader"
)

var u *uploader.Uploader

func TestNew(t *testing.T) {
	var err error
	u, err = uploader.New()
	if err != nil {
		t.Fatalf("fail to new: error=%v", err)
	}
}

func TestUpload(t *testing.T) {
	content := "test"

	f := storage.Image{
		ContentType: "text/plain; charset=utf-8",
		ETag:        fmt.Sprintf("%x", md5.Sum([]byte(content))),
		Filename:    "test/test.txt",
	}

	buf := bytes.NewBufferString(content)
	url, err := u.Upload(buf, f)
	if err != nil {
		t.Fatalf("fail to upload: error=%v", err)
	}

	// httpで取得してアップロードしたファイルをチェックする
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("fail to upload: error=%v", err)
	}
	defer resp.Body.Close()
	ct := resp.Header.Get("Content-Type")
	if ct != f.ContentType {
		t.Fatalf("wrong Content-Type: expected %s, but actual %s", f.ContentType, ct)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("fail to read body: error=%v", err)
	}
	b := string(body)
	if b != content {
		t.Fatalf("wrong body: expected %s, but actual %s", content, b)
	}
}

func TestCreateURL(t *testing.T) {
	bucket := os.Getenv(uploader.EnvGCSBucket)
	path := "baz"
	expected := fmt.Sprintf("https://%s.storage.googleapis.com/%s", bucket, path)
	actual := u.CreateURL(path)
	if actual != expected {
		t.Fatalf("fail to create URL: expected %s, but actual %s", expected, actual)
	}
}
