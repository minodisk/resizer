package fetcher_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/go-microservices/resizer/fetcher"
)

var (
	mockServer *httptest.Server
)

func TestMain(m *testing.M) {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		p := path.Join(dir, "..", "fixtures", r.URL.Path[1:])
		log.Println(p)
		http.ServeFile(w, r, p)
	}))

	code := m.Run()
	os.Exit(code)
}

func TestInit(t *testing.T) {
	if err := fetcher.Init(); err != nil {
		log.Fatal(err)
	}
}

func TestFetchAndClean(t *testing.T) {
	// モックサーバから期待値となるファイルのデータを取得する
	url := fmt.Sprintf("%s/f-png24.png", mockServer.URL)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("fail to get file %s: error=%v", url, err)
	}
	defer resp.Body.Close()
	expected, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("fail to read response body: error=%v", err)
	}

	// fetcher.Fetchを実行し、戻り値のパスにファイルが存在していることをテストする
	// 同一のデータが保存されていることをテストする
	filename, err := fetcher.Fetch(url)
	if err != nil {
		t.Fatalf("fail to Fetch: error=%v", err)
	}
	actual, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("fail to read file %s: error=%v", filename, err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("deferrent content between server file and local file")
	}

	// fetcher.Cleanを実行し、パスにファイルが存在していないことをテストする
	if err := fetcher.Clean(filename); err != nil {
		t.Fatalf("fail to clear: error=%v", err)
	}
	if _, err := os.Stat(filename); err == nil {
		t.Errorf("%s was not cleaned", filename)
	}
}
