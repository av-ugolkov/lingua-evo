package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"lingua-evo/internal/config"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/token"
	"lingua-evo/runtime"

	"github.com/google/uuid"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, ok := BearerAuth(r)
		if !ok {
			handler.SendError(w, http.StatusUnauthorized, fmt.Errorf("token not found"))
			return
		}
		claims, err := token.ValidateJWT(tokenStr, config.GetConfig().JWT.Secret)
		if err != nil {
			handler.SendError(w, http.StatusUnauthorized, err)
			return
		}

		r = r.WithContext(runtime.SetUserIDInContext(r.Context(), uuid.MustParse(claims.Subject)))
		next(w, r)
	})
}

func BearerAuth(r *http.Request) (token string, ok bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", false
	}
	return parseBasicAuth(auth)
}

func parseBasicAuth(auth string) (token string, ok bool) {
	const prefix = "Bearer "
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", false
	}
	return auth[len(prefix):], true
}
