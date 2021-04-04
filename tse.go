package tse

import (
	"context"
	"os"

	"github.com/empire/go-tse/internal/cli/counter"
	"github.com/empire/go-tse/internal/download"
	"github.com/empire/go-tse/internal/generator"
	"github.com/empire/go-tse/internal/symbols"
)

func DownloadAll(ctx context.Context, path string) error {
	if err := os.Mkdir(path, 0755); !(err == nil || os.IsExist(err)) {
		return err
	}

	results := make(chan download.Result)
	errc := make(chan error)

	defer close(results)
	defer close(errc)

	ids := symbols.Iter(ctx)
	generator := generator.New(path)
	c := counter.New()

	go generator.Run(ctx, results, errc)
	go c.Show(ctx, errc)

	download.Download(ctx, ids, results)

	return c.Err()
}
