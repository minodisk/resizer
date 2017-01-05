package log

import (
	"fmt"
	"io"
	"os"
)

var (
	OutWriter, ErrorWriter io.Writer
)

// NopWriter write nothing.
type NopWriter struct{}

// Write do nothing.
func (w NopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func init() {
	switch os.Getenv("RESIZER_ENV") {
	default:
		OutWriter = os.Stdout
		ErrorWriter = os.Stderr
	case "test":
		n := NopWriter{}
		OutWriter = n
		ErrorWriter = n
	}
}

func Println(args ...interface{}) {
	fmt.Fprintln(OutWriter, args...)
}

func Printf(format string, args ...interface{}) {
	fmt.Fprintf(OutWriter, format, args)
}
