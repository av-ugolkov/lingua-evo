package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/analytic"
	msgerror "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/token"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
)

func Auth(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken, err := ginExt.GetHeaderAuthorization(c, ginExt.AuthTypeBearer)
		if err != nil {
			ginExt.SendError(c, http.StatusUnauthorized,
				msgerror.NewError(fmt.Errorf("middleware.Auth: bearer token not found"),
					msgerror.ErrUnauthorized))
			return
		}
		claims, err := token.ValidateJWT(bearerToken, config.GetConfig().JWT.Secret)
		if err != nil {
			ginExt.SendError(c, http.StatusUnauthorized,
				msgerror.NewError(fmt.Errorf("middleware.Auth: %v", err),
					msgerror.ErrUnauthorized))
			return
		}
		c.Request = c.Request.WithContext(runtime.SetUserIDInContext(c.Request.Context(), claims.UserID))

		analytics.SendToKafka(claims.UserID, c.Request.URL.Path)

		next(c)
	}
}

func OptionalAuth(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken, err := ginExt.GetHeaderAuthorization(c, ginExt.AuthTypeBearer)
		if err != nil {
			slog.Warn(fmt.Sprintf("middleware.OptionalAuth: %v", err))
			next(c)
			return
		}
		claims, err := token.ValidateJWT(bearerToken, config.GetConfig().JWT.Secret)
		if err != nil {
			next(c)
			return
		}
		c.Request = c.Request.WithContext(runtime.SetUserIDInContext(c.Request.Context(), claims.UserID))

		analytics.SendToKafka(claims.UserID, c.Request.URL.Path)

		next(c)
	}
}
