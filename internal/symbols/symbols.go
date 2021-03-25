package symbols

import (
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
