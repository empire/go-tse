package counter

import (
	"context"
	"errors"
	"fmt"
)

type Counter struct {
	err error
}

func (c *Counter) Err() error {
	return c.err
}

func New() *Counter {
	return &Counter{}
}

func (c *Counter) Show(ctx context.Context, input <-chan error) {
	defer fmt.Println()
	counter := 0
	for {
		select {
		case err, ok := <-input:
			if !ok {
				return
			}
			if errors.Is(err, context.DeadlineExceeded) {
				continue
			}
			if err != nil {
				c.err = err
				continue
			}
			counter++
			fmt.Printf("\rCount: %d", counter)
		case <-ctx.Done():
			return
		}
	}
}
