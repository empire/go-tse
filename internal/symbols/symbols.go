package symbols

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
)

//go:embed symbols_name.json
var content []byte

var Symbols map[string]string

func init() {
	var symbolToIndex map[string]string
	Symbols = make(map[string]string)
	if err := json.Unmarshal(content, &symbolToIndex); err != nil {
		log.Fatal(err)
	}
	for k, v := range symbolToIndex {
		Symbols[v] = k
	}
}

func Iter(ctx context.Context) <-chan string {
	ids := make(chan string)
	go func() {
		iter(ctx, ids)
		close(ids)
	}()
	return ids
}

func iter(ctx context.Context, ids chan<- string) {
	for tickerIndex := range Symbols {
		select {
		case <-ctx.Done():
			break
		case ids <- tickerIndex:
		}
	}
}
