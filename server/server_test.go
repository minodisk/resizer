package server_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/minodisk/resizer/options"
	"github.com/minodisk/resizer/server"
	"github.com/minodisk/resizer/testutil"
	"github.com/pkg/errors"
)

var (
	appServer      *httptest.Server
	fixturesServer *httptest.Server
)

func TestMain(m *testing.M) {
	if err := testutil.CreateGoogleAuthFile(); err != nil {
		panic(err)
	}
	if err := testutil.DownloadFixtures("f-png24.png"); err != nil {
		panic(err)
	}

	fixturesServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(testutil.DirFixtures, r.URL.Path[1:]))
	}))
	u, err := url.Parse(fixturesServer.URL)
	if err != nil {
		panic(err)
	}

	h, err := server.NewHandler(options.Options{
		ServiceAccount: options.ServiceAccount{
			Path: testutil.GoogleAuthFilename,
		},
		DataSourceName: "root:@tcp(mysql:3306)/resizer?charset=utf8&parseTime=True",
		AllowedHosts:   []string{u.Host},
	})
	if err != nil {
		panic(err)
	}
	appServer = httptest.NewServer(http.HandlerFunc(h.ServeHTTP))

	c := m.Run()

	appServer.Close()
	fixturesServer.Close()

	if err := testutil.RemoveGoogleAuthFile(); err != nil {
		panic(err)
	}
	if err := testutil.RemoveFixtures(); err != nil {
		panic(err)
	}

	os.Exit(c)
}

func TestNew(t *testing.T) {
	t.Run("1st time", func(t *testing.T) {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				t.Fatalf("shouldn't be redirected")
				return nil
			},
		}
		resp, err := client.Get(fmt.Sprintf("%s?width=15&url=%s/f-png24.png", appServer.URL, fixturesServer.URL))
		if err != nil {
			t.Fatalf("fail to get resized image: %+v", err)
		}
		if resp.StatusCode != http.StatusOK {
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("body is empty")
			}
			t.Errorf("status code isn't OK: %s", b)
		}
	})

	time.Sleep(time.Second * 5)

	func() {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if strings.Contains(req.URL.Host, "googleapis") != true {
					t.Error("request at the 2nd time shouldn't be redirected")
				}
				return nil
			},
		}
		resp, err := client.Get(fmt.Sprintf("%s?width=15&url=%s/f-png24.png", appServer.URL, fixturesServer.URL))
		if err != nil {
			t.Fatalf("fail to get resized image at the 2st time: error=%v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("response at the 2nd time should be OK: expected %d, but actual %d", http.StatusOK, resp.StatusCode)
		}
	}()

	func() {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if strings.Contains(req.URL.Host, "googleapis") != true {
					t.Error("request at the 3rd time shouldn't be redirected")
				}
				return nil
			},
		}
		resp, err := client.Get(fmt.Sprintf("%s?height=21&url=%s/f-png24.png", appServer.URL, fixturesServer.URL))
		if err != nil {
			t.Errorf("fail to get resized image at the 3rd time: error=%v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("response at the 3rd time should be OK: expected %d, but actual %d", http.StatusOK, resp.StatusCode)
		}
	}()
}

var (
	rTitle   = regexp.MustCompile(`<title>(\d+ .+)<\/title>`)
	rH1      = regexp.MustCompile(`<h1>(.+)<\/h1>`)
	rP       = regexp.MustCompile(`<p>(.+)<\/p>`)
	rAddress = regexp.MustCompile(`<address>(.+)<\/address>`)
)

func TestFail(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s", appServer.URL))
	if err != nil {
		t.Fatalf("fail to get resized image at the 1st time: error=%v", err)
	}
	if a, e := resp.StatusCode, http.StatusBadRequest; a != e {
		t.Errorf("the status text is expected `%d`, but actual `%d`", e, a)
	}
	if a, e := fmt.Sprintf("%d %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)), resp.Status; a != e {
		t.Errorf("the status is expected `%s`, but actual `%s`", e, a)
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(errors.Wrap(err, "fail to read the response body"))
	}
	body := string(buf)
	if a, e := rTitle.FindStringSubmatch(body)[1], "400 Bad Request"; a != e {
		t.Errorf("the status in <title> is expected `%s`, but actual `%s`", e, a)
	}
	if a, e := rH1.FindStringSubmatch(body)[1], "Bad Request"; a != e {
		t.Errorf("<h1> is expected `%s`, but actual `%s`", e, a)
	}
	if a, e := rP.FindStringSubmatch(body)[1], "URL shouldn't be empty"; a != e {
		t.Errorf("<p> is expected `%s`, but actual `%s`", e, a)
	}
	if a, e := rAddress.FindStringSubmatch(body)[1], "Resizer"; a != e {
		t.Errorf("the application name in <address> is expected `%s`, but actual `%s`", e, a)
	}
}
