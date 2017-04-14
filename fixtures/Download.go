package fixtures

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
)

const (
	BucketURL = "https://resizer-fixtures.storage.googleapis.com/"
)

var (
	Status    = "none"
	Filenames = []string{
		"huge.jpg",
	}
)

func Download() {
	var wg sync.WaitGroup
	var e chan error
	for _, f := range Filenames {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			resp, err := http.Get(path.Join(BucketURL, filename))
			if err != nil {
				e <- err
				return
			}
			file, err := os.Create(filename)
			if err != nil {
				e <- err
				return
			}
			defer resp.Body.Close()
			n, err := io.Copy(file, resp.Body)
			if err != nil {
				e <- err
				return
			}
			fmt.Printf("complete to download %s: %d bytes\n", filename, n)
		}(f)
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
}
