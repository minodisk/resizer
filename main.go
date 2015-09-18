package main

import (
	"runtime"

	"github.com/go-microservices/resizer/fetcher"
	"github.com/go-microservices/resizer/log"
	"github.com/go-microservices/resizer/server"
)

func main() {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Println("RECOVER!!!!!")
	// 		w, err := os.Create("panic.log")
	// 		if err != nil {
	// 			os.Exit(1)
	// 		}
	// 		os.Stderr = w
	// 		debug.PrintStack()
	// 	}

	// runtime.SetBlockProfileRate(1)
	// go func() {
	// 	log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	// }()

	runtime.GOMAXPROCS(runtime.NumCPU())

	filename, err := log.Init()
	if err != nil {
		panic(err)
	}
	if filename != "" {
		log.Infof("%s is created", filename)
	}

	if err := fetcher.Init(); err != nil {
		panic(err)
	}

	if err := server.Start(); err != nil {
		panic(err)
	}
}
