package middleware

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/analytic"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/pkg/http/gin_extension"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/token"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
)

func Auth(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken, err := ginExt.GetHeaderAuthorization(c, ginExt.AuthTypeBearer)
		if err != nil {
			ginExt.SendError(c, http.StatusUnauthorized,
				fmt.Errorf("middleware.Auth: bearer token not found"))
			return
		}
		claims, err := token.ValidateJWT(bearerToken, config.GetConfig().JWT.Secret)
		if err != nil {
			ginExt.SendError(c, http.StatusUnauthorized, fmt.Errorf("middleware.Auth: %w", err))
			return
		}
		c.Request = c.Request.WithContext(runtime.SetUserIDInContext(c.Request.Context(), claims.UserID))

		analytics.SendToKafka(claims.UserID, c.Request.URL.Path)

		next(c)
	}
}
