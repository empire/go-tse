package main

import (
	"os"

	"github.com/empire/go-tse"
)

func main() {
	os.Mkdir("/tmp/symbols", 0755)

	ids := make(chan string, 10)
	for i := 0; i < 10; i++ {
		go tse.Download(ids)
	}
	tse.QueueTickers(ids)
}
