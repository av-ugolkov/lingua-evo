package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/analytic"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/http/exchange"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/token"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

type ExchangerFunc func(ctx context.Context, ex *exchange.Exchanger)

func Auth(next ExchangerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		analytics.SendToKafka(claims.UserID, r.URL.Path)

		next(ctx, ex)
	}
}
