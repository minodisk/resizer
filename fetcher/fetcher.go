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
)

const (
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36"
)

var (
	tempDir string
	expires time.Duration
	client  *http.Client
)

func Init() error {
	tempDir = path.Join(os.TempDir(), "")
	log.Println(tempDir)

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
	t := log.Start()
	defer log.End(t)

	sum := md5.Sum([]byte(fmt.Sprintf("%s-%d", url, time.Now().UnixNano())))
	f := fmt.Sprintf("%x", sum)
	filename := path.Join(tempDir, f)
	log.Debugf("file is temporary saved as %s", filename)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Errorf("fail to new request: error=%v", err)
		return "", err
	}
	req.Header.Set("User-Agent", UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("fail to GET %s: error=%v", url, err)
		return "", err
	}

	dump, _ := httputil.DumpRequest(req, true)
	log.Debugf("%s", dump)

	if resp.StatusCode != http.StatusOK {
		log.Printf("not ok: StatusCode=%d", resp.StatusCode)
		return "", fmt.Errorf("can't fetch image %s", url)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error(err)
		}
	}()
	log.Printf("ok: StatusCode=%d", resp.StatusCode)

	file, err := os.Create(filename)
	defer func() {
		if err := file.Close(); err != nil {
			log.Error(err)
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
