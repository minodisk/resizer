package processor_test

import (
	"testing"

	"github.com/go-microservices/resizer/input"
	"github.com/go-microservices/resizer/processor"
	"github.com/go-microservices/resizer/storage"
)

type NopWriter struct{}

func (w NopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func BenchmarkHuge(b *testing.B) {
	if err := process("../fixtures/huge.jpg"); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkLarge(b *testing.B) {
	if err := process("../fixtures/large.jpg"); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkMedium(b *testing.B) {
	if err := process("../fixtures/medium.jpg"); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkSmall(b *testing.B) {
	if err := process("../fixtures/small.jpg"); err != nil {
		b.Fatal(err)
	}
}

// func BenchmarkMemory(b *testing.B) {
// 	go func() {
// 		var mem runtime.MemStats
// 		for {
// 			runtime.ReadMemStats(&mem)
// 			fmt.Printf("%d\t%d\t%d\t%d\n", mem.Alloc, mem.TotalAlloc, mem.HeapAlloc, mem.HeapSys)
// 			time.Sleep(time.Second * 1)
// 		}
// 	}()
//
// 	wg := new(sync.WaitGroup)
// 	for i := 0; i < 100; i++ {
// 		wg.Add(1)
// 		go func() {
// 			if err := process(b, "../fixtures/medium.jpg"); err != nil {
// 				b.Fatal(err)
// 			}
// 			wg.Done()
// 		}()
// 	}
// 	wg.Wait()
// }

func process(path string) error {
	p := processor.New()
	i, err := storage.NewImage(input.Input{
		URL:   "http://example.com/test.jpg",
		Width: 800,
	})
	if err != nil {
		return err
	}

	var w NopWriter
	if _, err = p.Process(path, &w, i); err != nil {
		return err
	}
	return nil
}
