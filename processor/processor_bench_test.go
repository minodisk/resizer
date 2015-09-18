package processor_test

import (
	"bytes"
	"testing"

	"github.com/go-microservices/resizer/fetcher"
	"github.com/go-microservices/resizer/processor"
	"github.com/go-microservices/resizer/storage"
)

func BenchmarkProcess(b *testing.B) {
	if err := fetcher.Init(); err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		o, err := storage.NewImage(map[string][]string{
			"url":    []string{""},
			"width":  []string{"800"},
			"height": []string{"0"},
		})
		if err != nil {
			b.Fatalf("%s", err)
		}

		path, err := fetcher.Fetch(o.ValidatedURL)
		if err != nil {
			b.Fatalf("%s", err.Error())
		}

		var bs []byte
		w := bytes.NewBuffer(bs)
		p := processor.New()

		image, err := p.Preprocess(path)
		if err != nil {
			b.Fatalf("%s", err.Error())
		}

		if _, err := p.Process(image, w, o); err != nil {
			b.Fatalf("%s", err.Error())
		}
	}
}
