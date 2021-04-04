package main

import (
	"context"
	"log"
	"time"

	"github.com/empire/go-tse"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	err := tse.DownloadAll(ctx, "/tmp/symbols")
	if err != nil {
		log.Fatal(err)
	}
}
