package generator

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/empire/go-tse/internal/download"
	"github.com/go-gota/gota/dataframe"
)

type Geneator struct {
	path string
}

func New(path string) *Geneator {
	return &Geneator{
		path: path,
	}
}

type geneatorErrorWriter struct {
	*Geneator
	err error
}

func (g *Geneator) Run(ctx context.Context, input chan download.Result, output chan<- error) {
	for {
		select {
		case <-ctx.Done():
		case result, ok := <-input:
			if !ok {
				return
			}
			w := geneatorErrorWriter{g, result.Err}
			w.toCsv(result.Ticker, result.Content)
			output <- w.err
		}
	}
}

func (gew *geneatorErrorWriter) toCsv(ticker string, out []byte) {
	if gew.err != nil {
		return
	}

	fileName := fmt.Sprintf("%s/%s.csv", gew.path, ticker)
	gew.err = toCsv(fileName, out)
}

func toCsv(fileName string, out []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	df := dataframe.ReadCSV(bytes.NewReader(out))
	HISTORY_FIELD_MAPPINGS := map[string]string{
		"<CLOSE>":      "adjClose",
		"<DTYYYYMMDD>": "date",
		"<FIRST>":      "open",
		"<HIGH>":       "high",
		"<LAST>":       "close",
		"<LOW>":        "low",
		"<OPENINT>":    "count",
		"<VALUE>":      "value",
		"<VOL>":        "volume",
	}
	for oldName, newName := range HISTORY_FIELD_MAPPINGS {
		df = df.Rename(newName, oldName)
	}
	df = df.Drop([]string{"<PER>", "<OPEN>", "<TICKER>"})
	return df.WriteCSV(file)
}
