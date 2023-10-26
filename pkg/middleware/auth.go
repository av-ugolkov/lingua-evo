package middleware

import (
	"fmt"
	"lingua-evo/internal/config"
	"lingua-evo/pkg/jwt"
	"lingua-evo/pkg/tools"
	"log/slog"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header["Authorization"]; !ok {
			tools.SendError(w, http.StatusUnauthorized, fmt.Errorf("token not found"))
			return
		}
		auth := r.Header["Authorization"][0]
		token := strings.Split(auth, " ")[1]
		claims, err := jwt.ValidateToken(token, config.GetConfig().JWT.Secret)
		if err != nil {
			tools.SendError(w, http.StatusUnauthorized, err)
			return
		}

		slog.Info(claims.ID.String())

		c, err := r.Cookie("refresh_token")
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
		}
		fmt.Printf("token: %v\n", c)
		next.ServeHTTP(w, r)
	})
}
