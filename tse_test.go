package tse

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
)

const sample = `<TICKER>,<DTYYYYMMDD>,<FIRST>,<HIGH>,<LOW>,<CLOSE>,<VALUE>,<VOL>,<OPENINT>,<PER>,<OPEN>,<LAST>
Sample.Ticker,20210314,116028.00,116028.00,116028.00,114127.00,9208562220,79365,61,D,112649.00,116028.00`

type transporter struct {
	urls []*url.URL
}

func (t *transporter) RoundTrip(r *http.Request) (*http.Response, error) {
	t.urls = append(t.urls, r.URL)
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(sample)),
		Header:     make(http.Header),
	}, nil
}

func TestA(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	http.DefaultClient.Transport = &transporter{}
	err := DownloadAll(ctx, "/tmp/out")
	if err != nil {
		t.Error(err)
	}
}
