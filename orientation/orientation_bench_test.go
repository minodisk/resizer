package orientation_test

import (
	"testing"

	"github.com/go-microservices/resizer/orientation"
	"github.com/go-microservices/resizer/processor"
)

func orient(b *testing.B, filename string, o int) {
	i, _, err := processor.Load(filename)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	orientation.Orient(i, o)
}

func BenchmarkOrient1(b *testing.B) {
	orient(b, "../fixtures/huge-1.jpg", 1)
}

func BenchmarkOrient2(b *testing.B) {
	orient(b, "../fixtures/huge-2.jpg", 2)
}

func BenchmarkOrient3(b *testing.B) {
	orient(b, "../fixtures/huge-3.jpg", 3)
}

func BenchmarkOrient4(b *testing.B) {
	orient(b, "../fixtures/huge-4.jpg", 4)
}

func BenchmarkOrient5(b *testing.B) {
	orient(b, "../fixtures/huge-5.jpg", 5)
}

func BenchmarkOrient6(b *testing.B) {
	orient(b, "../fixtures/huge-6.jpg", 6)
}

func BenchmarkOrient7(b *testing.B) {
	orient(b, "../fixtures/huge-7.jpg", 7)
}

func BenchmarkOrient8(b *testing.B) {
	orient(b, "../fixtures/huge-8.jpg", 8)
}
