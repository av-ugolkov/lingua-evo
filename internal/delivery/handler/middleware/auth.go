package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/analytic"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	msgerror "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/token"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

func Auth(next ginext.HandlerFunc) ginext.HandlerFunc {
	return func(c *ginext.Context) (int, any, error) {
		bearerToken, err := c.GetHeaderAuthorization(ginext.AuthTypeBearer)
		if err != nil {
			return http.StatusUnauthorized, nil,
				msgerror.NewError(fmt.Errorf("middleware.Auth: bearer token not found"),
					msgerror.ErrUnauthorized)
		}
		claims, err := token.ValidateJWT(bearerToken, config.GetConfig().JWT.Secret)
		if err != nil {
			return http.StatusUnauthorized, nil,
				msgerror.NewError(fmt.Errorf("middleware.Auth: %v", err),
					msgerror.ErrUnauthorized)
		}
		c.Request = c.Request.WithContext(runtime.SetUserIDInContext(c.Request.Context(), claims.UserID))

		analytics.SendToKafka(claims.UserID, c.Request.URL.Path)

		return next(c)
	}
}

func OptionalAuth(next ginext.HandlerFunc) ginext.HandlerFunc {
	return func(c *ginext.Context) (int, any, error) {
		bearerToken, err := c.GetHeaderAuthorization(ginext.AuthTypeBearer)
		if err != nil {
			slog.Warn(fmt.Sprintf("middleware.OptionalAuth: %v", err))
			return next(c)
		}
		claims, err := token.ValidateJWT(bearerToken, config.GetConfig().JWT.Secret)
		if err != nil {
			return next(c)
		}
		c.Request = c.Request.WithContext(runtime.SetUserIDInContext(c.Request.Context(), claims.UserID))

		analytics.SendToKafka(claims.UserID, c.Request.URL.Path)

		return next(c)
	}
}
