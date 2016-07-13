package storage_test

import (
	"os"
	"testing"

	"github.com/go-microservices/resizer/option"
	"github.com/go-microservices/resizer/storage"
)

var s *storage.Storage

func TestInit(t *testing.T) {
	var err error
	o, err := option.New(os.Args[1:])
	if err != nil {
		t.Fatalf("fail to create options: error=%v", err)
	}
	s, err = storage.New(o)
	if err != nil {
		t.Fatalf("Can't create instance: %v", err)
	}
}

func TestCache(t *testing.T) {
	file := storage.Image{
		ValidatedURL:     "http://example.com/foo.jpg",
		ValidatedWidth:   400,
		ValidatedHeight:  300,
		ValidatedMethod:  storage.MethodDefault,
		ValidatedFormat:  storage.FormatDefault,
		ValidatedQuality: 100,
	}
	s.NewRecord(file)
	s.Create(&file)
}
