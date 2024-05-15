package log

import (
	"context"
	"io"
	"log"
	"log/slog"
)

type color string

var (
	Red     color = "\033[31m"
	Green   color = "\033[32m"
	Yellow  color = "\033[33m"
	Blue    color = "\033[34m"
	Magenta color = "\033[35m"
	Cyan    color = "\033[36m"
	Gray    color = "\033[37m"
	White   color = "\033[97m"
)

type CustomHandler struct {
	slog.Handler
	l *log.Logger
}

func NewCustomHandler(writers io.Writer, opts *slog.HandlerOptions) *CustomHandler {
	return &CustomHandler{
		Handler: slog.NewTextHandler(writers, opts),
		l:       log.New(writers, "", 0),
	}
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	timeStr := r.Time.Format("[02/01/2006 15:05:05.000]")

	h.l.Println(timeStr, level, r.Message)

	return nil
}

// func setColor(value string, c color) string {
// 	return fmt.Sprintf("%s%s\033[0m", c, value)
// }
