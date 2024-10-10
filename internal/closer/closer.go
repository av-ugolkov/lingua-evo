package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Func func(ctx context.Context) error

var (
	mu    sync.Mutex
	funcs []Func
)

func Add(f Func) {
	mu.Lock()
	defer mu.Unlock()
	funcs = append(funcs, f)
}

func Close(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	var (
		complete = make(chan struct{}, 1)
		errs     = make([]string, 0, len(funcs))
	)

	go func() {
		for _, f := range funcs {
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
