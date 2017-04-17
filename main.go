package main

import "github.com/minodisk/resizer/server"

func main() {
	if err := server.Start(); err != nil {
		panic(err)
	}
}
