package middleware

import (
	"errors"
	"net/http"
	"strings"

	"lingua-evo/internal/config"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/token"
	"lingua-evo/runtime"
)

var (
	errNotFoundToken = errors.New("token not found")
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := handler.NewHandler(w, r)
		var auth string
		var err error
		if auth, err = handler.GetHeaderAuthorization(); err != nil {
			handler.SendError(http.StatusUnauthorized, errNotFoundToken)
			return
		}
		tokenStr := strings.Split(auth, " ")[1]
		claims, err := token.ValidateJWT(tokenStr, config.GetConfig().JWT.Secret)
		if err != nil {
			handler.SendError(http.StatusUnauthorized, err)
			return
		}
		r = r.WithContext(runtime.SetUserIDInContext(r.Context(), claims.UserID))
		next(w, r)
	})
}
