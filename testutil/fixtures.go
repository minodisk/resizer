package testutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/storage"

	"google.golang.org/api/option"
)

const (
	DirFixtures = "fixtures"
)

func DownloadFixtures(names ...string) error {
	if err := os.MkdirAll(DirFixtures, 0755); err != nil {
		return err
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithServiceAccountFile(GoogleAuthFilename))
	if err != nil {
		return err
	}
	bkt := client.Bucket("resizer")

	var wg sync.WaitGroup
	var e chan error
	for _, name := range names {
		obj := bkt.Object(fmt.Sprintf("fixtures/%s", name))
		attrs, err := obj.Attrs(ctx)
		if err != nil {
			return err
		}
		if _, err := os.Stat(attrs.Name); err == nil {
			continue
		}

		wg.Add(1)
		go func(attrs *storage.ObjectAttrs) {
			if err := func() error {
				resp, err := http.Get(attrs.MediaLink)
				if err != nil {
					return err
				}
				file, err := os.Create(attrs.Name)
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
		}(attrs)
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
