package middleware

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/pkg/http/exchange"
	"github.com/av-ugolkov/lingua-evo/pkg/token"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
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
		r = r.WithContext(runtime.SetUserIDInContext(r.Context(), claims.UserID))
		next(w, r)
	})
}
