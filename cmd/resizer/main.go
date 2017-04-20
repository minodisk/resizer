package main

import (
	"fmt"
	"os"

	"github.com/minodisk/resizer/options"
	"github.com/minodisk/resizer/server"
)

func main() {
	if err := _main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func _main() error {
	o, err := options.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	return server.Start(o)
}
