package fetcher

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"time"

	"github.com/go-microservices/resizer/log"
	"github.com/pkg/errors"
)

const (
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36"
)

var (
	tempDir = path.Join(os.TempDir(), "resizer")
	expires time.Duration
	client  *http.Client
)

func init() {
	if err := _init(); err != nil {
		panic(err)
	}
}

func _init() error {
	var err error
	expires, err = time.ParseDuration("1h")
	if err != nil {
		return err
	}
	if err := os.RemoveAll(tempDir); err != nil {
		return err
	}
	if err := os.MkdirAll(tempDir, 0777); err != nil {
		return err
	}

	client = new(http.Client)

	return nil
}

func Fetch(url string) (string, error) {
	sum := md5.Sum([]byte(fmt.Sprintf("%s-%d", url, time.Now().UnixNano())))
	f := fmt.Sprintf("%x", sum)
	filename := path.Join(tempDir, f)

	log.Printf("file is temporary saved as %s\n", filename)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "fail to new request")
	}
	req.Header.Set("User-Agent", UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "fail to GET")
	}

	dump, _ := httputil.DumpRequest(req, true)

	log.Printf("Dump: %s\n", dump)

	if resp.StatusCode != http.StatusOK {
		log.Printf("not ok: StatusCode=%d\n", resp.StatusCode)
		return "", fmt.Errorf("can't fetch image %s", url)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(errors.Wrap(err, "fail to close response body"))
		}
	}()
	log.Printf("ok: StatusCode=%d\n", resp.StatusCode)

	file, err := os.Create(filename)
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(errors.Wrap(err, "fail to close file"))
		}
	}()
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", err
	}

	return filename, nil
}

func Clean(filename string) error {
	return os.Remove(filename)
}
