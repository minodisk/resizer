package fetcher_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/minodisk/resizer/fetcher"
	"github.com/minodisk/resizer/testutil"
)

var (
	mockServer *httptest.Server
)

func TestMain(m *testing.M) {
	if err := testutil.DownloadFixtures("f-png24.png"); err != nil {
		panic(err)
	}

	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, testutil.DirFixtures)
	}))

	code := m.Run()
	if err := testutil.RemoveFixtures(); err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
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
