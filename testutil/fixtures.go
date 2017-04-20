package testutil

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const (
	DirFixtures = "fixtures"
)

func DownloadFixtures(names ...string) error {
	if err := os.MkdirAll(DirFixtures, 0755); err != nil {
		return err
	}

	var wg sync.WaitGroup
	var e chan error
	for _, name := range names {
		if _, err := os.Stat(name); err == nil {
			continue
		}

		wg.Add(1)
		go func(filename string) {
			if err := func() error {
				url := fmt.Sprintf("https://resizer.storage.googleapis.com/fixtures/%s", filename)
				resp, err := http.Get(url)
				if err != nil {
					return err
				}
				file, err := os.Create(filepath.Join(DirFixtures, filename))
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				if _, err := io.Copy(file, resp.Body); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				e <- err
			}
			wg.Done()
		}(name)
	}
	go func() {
		for {
			select {
			case err := <-e:
				panic(err)
			}
		}
	}()
	wg.Wait()

	return nil
}

func RemoveFixtures() error {
	return os.RemoveAll(DirFixtures)
}
