package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Func func(ctx context.Context) error

type Closer struct {
	mu    sync.Mutex
	funcs []Func
}

func NewCloser() *Closer {
	return &Closer{}
}

func (c *Closer) Add(f Func) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		complete = make(chan struct{}, 1)
		errs     = make([]string, 0, len(c.funcs))
	)

	go func() {
		for _, f := range c.funcs {
			if err := f(ctx); err != nil {
				errs = append(errs, err.Error())
			}

			complete <- struct{}{}
		}
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return fmt.Errorf("shutdown cancelled: %v", ctx.Err())
	}

	if len(errs) > 0 {
		return fmt.Errorf(
			"shutdown finished with error(s): \n%s",
			strings.Join(errs, "\n"),
		)
	}

	return nil
}
