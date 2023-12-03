package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"lingua-evo/internal/config"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/token"
	"lingua-evo/runtime"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header["Authorization"]; !ok {
			handler.SendError(w, http.StatusUnauthorized, fmt.Errorf("token not found"))
			return
		}
		auth := r.Header["Authorization"][0]
		tokenStr := strings.Split(auth, " ")[1]
		claims, err := token.ValidateJWT(tokenStr, config.GetConfig().JWT.Secret)
		if err != nil {
			handler.SendError(w, http.StatusUnauthorized, err)
			return
		}

		r = r.WithContext(runtime.SetUserIDInContext(r.Context(), claims.UserID))
		next(w, r)
	})
}
