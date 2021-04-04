package counter

import (
	"context"
	"errors"
	"fmt"
)

type Counter struct {
	errs chan error
	err  error
}

func (c *Counter) Input() chan<- error {
	return c.errs
}

func (c *Counter) Err() error {
	return c.err
}

func New() *Counter {
	return &Counter{
		errs: make(chan error),
	}
}

func (c *Counter) Show() {
	counter := 0
	for err := range c.errs {
		if errors.Is(err, context.DeadlineExceeded) {
			continue
		}
		if err != nil {
			c.err = err
			continue
		}
		counter++
		fmt.Printf("\rCount: %d", counter)
	}
	fmt.Println()
}

func (c *Counter) Close() {
	close(c.errs)
}
