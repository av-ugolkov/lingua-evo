package jwt

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"lingua-evo/internal/config"

	"github.com/cristalhq/jwt/v3"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

func Middleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			slog.Error("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("malformed token"))
			return
		}
		slog.Debug("create jwt verifier")
		jwtToken := authHeader[1]
		key := []byte(config.GetConfig().JWT.Secret)
		verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
		if err != nil {
			unauthorized(w, err)
			return
		}

		slog.Debug("parse and verify token")
		token, err := jwt.ParseAndVerifyString(jwtToken, verifier)
		if err != nil {
			unauthorized(w, err)
			return
		}
		slog.Debug("parse user claims")
		var uc UserClaims
		err = json.Unmarshal(token.RawClaims(), &uc)
		if err != nil {
			unauthorized(w, err)
			return
		}
		if valid := uc.IsValidAt(time.Now()); !valid {
			slog.Error("token has been expired")
			unauthorized(w, err)
			return
		}

		var userUUID ContextKey = "user_uuid"
		ctx := context.WithValue(r.Context(), userUUID, uc.ID)
		h(w, r.WithContext(ctx))
	}
}

type ContextKey string

func unauthorized(w http.ResponseWriter, err error) {
	slog.Error(err.Error())
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte("unauthorized"))
}
