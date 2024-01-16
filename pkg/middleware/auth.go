package middleware

import (
	"fmt"
	"net/http"

	"lingua-evo/internal/config"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/http/handler/common"
	"lingua-evo/pkg/token"
	"lingua-evo/runtime"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := handler.NewHandler(w, r)
		var bearerToken string
		var err error
		if bearerToken, err = handler.GetHeaderAuthorization(common.AuthTypeBearer); err != nil {
			handler.SendError(http.StatusUnauthorized, fmt.Errorf("middleware.Auth: %w", err))
			return
		}
		claims, err := token.ValidateJWT(bearerToken, config.GetConfig().JWT.Secret)
		if err != nil {
			handler.SendError(http.StatusUnauthorized, err)
			return
		}
		r = r.WithContext(runtime.SetUserIDInContext(r.Context(), claims.UserID))
		next(w, r)
	})
}
