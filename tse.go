package tse

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-gota/gota/dataframe"
)

const (
	TSE_TICKER_EXPORT_DATA_ADDRESS = "http://tsetmc.com/tsev2/data/Export-txt.aspx?t=i&a=1&b=0&i=%s"
	TSE_TICKER_ADDRESS             = "http://tsetmc.com/Loader.aspx?ParTree=151311&i=%s"
	TSE_ISNT_INFO_URL              = "http://www.tsetmc.com/tsev2/data/instinfofast.aspx?i=%s&c=57+"
	TSE_CLIENT_TYPE_DATA_URL       = "http://www.tsetmc.com/tsev2/data/clienttype.aspx?i=%s"
	TSE_SYMBOL_ID_URL              = "http://www.tsetmc.com/tsev2/data/search.aspx?skey=%s"
	TSE_SHAREHOLDERS_URL           = "http://www.tsetmc.com/Loader.aspx?Partree=15131T&c=%s"
)

var symbols map[string]string

func init() {
	var indexToSymbol map[string]string
	symbols = make(map[string]string)
	if err := json.Unmarshal(content, &indexToSymbol); err != nil {
		log.Fatal(err)
	}
	for k, v := range indexToSymbol {
		symbols[v] = k
	}
}

//go:embed symbols/symbols_name.json
var content []byte

func QueueTickers(ids chan string) {
	counter := 0
	for ticker_index := range symbols {
		counter++
		fmt.Printf("\r%d", counter)
		ids <- ticker_index
	}
	fmt.Println()
	close(ids)
}

func Download(ids chan string) {
	for ticker_index := range ids {
		out := fmt.Sprintf("/tmp/symbols/%s.csv", symbols[ticker_index])
		if err := download_daily_record(out, ticker_index); err != nil {
			log.Println(err)
		}
	}
}

func download_daily_record(out, tiker_index string) error {
	url := fmt.Sprintf(TSE_TICKER_EXPORT_DATA_ADDRESS, tiker_index)
	return downloadFile(url, out)
}

func downloadFile(url, fileName string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	df := dataframe.ReadCSV(response.Body)
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
	df.WriteCSV(file)
	return df.Err
}
