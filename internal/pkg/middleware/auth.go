package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/kafka"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/http/exchange"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/token"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"github.com/google/uuid"
)

type ExangerFunc func(ctx context.Context, ex *exchange.Exchanger)

var (
	topic          = []byte("user_action")
	middlewareAuth MiddlewareAuth
)

type Action struct {
	UserID    uuid.UUID `json:"uid"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
}

type MiddlewareAuth struct {
	writer *kafka.Writer
}

func NewMiddlewareAuth(w *kafka.Writer) {
	middlewareAuth = MiddlewareAuth{
		writer: w,
	}
}

func Auth(next ExangerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ex := exchange.NewExchanger(w, r)
		var bearerToken string
		var err error
		if bearerToken, err = ex.GetHeaderAuthorization(exchange.AuthTypeBearer); err != nil {
			ex.SendError(http.StatusUnauthorized, fmt.Errorf("middleware.Auth: %w", err))
			return
		}
		claims, err := token.ValidateJWT(bearerToken, config.GetConfig().JWT.Secret)
		if err != nil {
			ex.SendError(http.StatusUnauthorized, err)
			return
		}
		ctx := runtime.SetUserIDInContext(r.Context(), claims.UserID)

		go func() {
			b, err := json.Marshal(Action{
				UserID:    claims.UserID,
				Action:    r.URL.Path,
				CreatedAt: time.Now().UTC(),
			})
			if err != nil {
				slog.Error(fmt.Errorf("middleware.Auth: %v", err).Error())
				return
			}
			err = middlewareAuth.writer.SendMessage(context.Background(), topic, b)
			if err != nil {
				slog.Error(fmt.Errorf("middleware.MiddlewareAnalytics.SendData: %v", err).Error())
			}
		}()

		next(ctx, exchange.NewExchanger(w, r))
	})
}
