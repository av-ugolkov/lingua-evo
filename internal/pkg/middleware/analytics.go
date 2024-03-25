package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/kafka"
)

type MiddlewareAnalytics struct {
	writer *kafka.Writer
}

func NewMiddleware(w *kafka.Writer) *MiddlewareAnalytics {
	return &MiddlewareAnalytics{
		writer: w,
	}
}

func (m *MiddlewareAnalytics) SendData(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		go func() {
			err := m.writer.SendMessage(r.Context(), []byte("analytics"), []byte(r.URL.Path))
			if err != nil {
				slog.Error(fmt.Errorf("middleware.MiddlewareAnalytics.SendData: %w", err).Error())
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
