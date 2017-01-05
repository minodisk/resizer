package processor_test

import (
	"image"
	"testing"

	"github.com/go-microservices/resizer/input"
	"github.com/go-microservices/resizer/processor"
	"github.com/go-microservices/resizer/storage"
)

var (
	huge image.Image
)

func init() {
	var err error
	huge, _, err = processor.Load("../fixtures/huge.jpg")
	if err != nil {
		panic(err)
	}
}

type NopWriter struct{}

func (w NopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func BenchmarkProcessHuge(b *testing.B) {
	if err := process("../fixtures/huge.jpg"); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkProcessLarge(b *testing.B) {
	if err := process("../fixtures/large.jpg"); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkProcessMedium(b *testing.B) {
	if err := process("../fixtures/medium.jpg"); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkProcessSmall(b *testing.B) {
	if err := process("../fixtures/small.jpg"); err != nil {
		b.Fatal(err)
	}
}

func process(path string) error {
	p := processor.New()
	input := input.Input{
		URL:   "http://example.com/test.jpg",
		Width: 800,
	}
	input, err := input.Validate([]string{"example.com"})
	i, err := storage.NewImage(input)
	if err != nil {
		return err
	}

	var w NopWriter
	if _, err = p.Process(path, &w, i); err != nil {
		return err
	}
	return nil
}

func BenchmarkRotate1(b *testing.B) {
	if err := rotateHuge(1); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkRotate2(b *testing.B) {
	if err := rotateHuge(2); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkRotate3(b *testing.B) {
	if err := rotateHuge(3); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkRotate4(b *testing.B) {
	if err := rotateHuge(4); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkRotate5(b *testing.B) {
	if err := rotateHuge(5); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkRotate6(b *testing.B) {
	if err := rotateHuge(6); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkRotate7(b *testing.B) {
	if err := rotateHuge(7); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkRotate8(b *testing.B) {
	if err := rotateHuge(8); err != nil {
		b.Fatal(err)
	}
}

func rotateHuge(o int) error {
	if _, err := processor.Rotate(huge, o); err != nil {
		return err
	}

	return nil
}
