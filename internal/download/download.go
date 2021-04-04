package download

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/empire/go-tse/internal/symbols"
)

const (
	TSE_TICKER_EXPORT_DATA_ADDRESS = "http://tsetmc.com/tsev2/data/Export-txt.aspx?t=i&a=1&b=0&i=%s"
)

type Result struct {
	Ticker  string
	Content []byte
	Err     error
}

func Download(ctx context.Context, ids <-chan string, out chan<- Result) {
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			download(ctx, ids, out)
		}()
	}
	wg.Wait()
}

func download(ctx context.Context, in <-chan string, out chan<- Result) {
	for tickerIndex := range in {
		content, err := downloadDailyRecord(ctx, tickerIndex)
		select {
		case out <- Result{Ticker: symbols.Symbols[tickerIndex], Content: content, Err: err}:
		case <-ctx.Done():
			break
		}
	}
}

func downloadDailyRecord(ctx context.Context, tiker_index string) ([]byte, error) {
	url := fmt.Sprintf(TSE_TICKER_EXPORT_DATA_ADDRESS, tiker_index)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("Received non 200 response code")
	}

	out, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return out, nil
}
