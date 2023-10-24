package middleware

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*jwtToken, err := jwt.NewToken(config.GetConfig().JWT.Secret)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		r.Header.Set("jwt", jwtToken)*/
		next.ServeHTTP(w, r)
	})
}
