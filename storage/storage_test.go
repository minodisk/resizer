package storage_test

import (
	"testing"

	"github.com/go-microservices/resizer/storage"
)

var s *storage.Storage

func TestInit(t *testing.T) {
	var err error
	s, err = storage.New()
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
