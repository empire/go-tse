package tse

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/empire/go-tse/internal/symbols"
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

func DownloadAll(path string) {
	os.Mkdir("/tmp/symbols", 0755)
	ids := make(chan string, 10)
	for i := 0; i < 10; i++ {
		go download(ids)
	}
	queueTickers(ids)
}

func queueTickers(ids chan string) {
	counter := 0
	for ticker_index := range symbols.Symbols {
		counter++
		fmt.Printf("\r%d", counter)
		ids <- ticker_index
	}
	fmt.Println()
	close(ids)
}

func download(ids chan string) {
	for ticker_index := range ids {
		out := fmt.Sprintf("/tmp/symbols/%s.csv", symbols.Symbols[ticker_index])
		if err := downloadDailyRecord(out, ticker_index); err != nil {
			log.Println(err)
		}
	}
}

func downloadDailyRecord(out, tiker_index string) error {
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
