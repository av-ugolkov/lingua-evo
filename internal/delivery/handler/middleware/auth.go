package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/analytic"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/fext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/router"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/token"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gofiber/fiber/v2"
)

func Auth(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bearerToken, err := GetTokenAuth(c, router.AuthTypeBearer)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
		}

		claims, err := token.ValidateJWT(bearerToken, config.GetConfig().JWT.Secret)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
		}

		fext.SetUserIDToContext(c, claims.UserID)
		analytics.SendToKafka(claims.UserID, c.OriginalURL())

		return next(c)
	}
}

func OptionalAuth(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bearerToken, err := GetTokenAuth(c, router.AuthTypeBearer)
		if err != nil {
			slog.Warn(fmt.Sprintf("delivery.handler.middleware.OptionalAuth: %v", err))
			return next(c)
		}
		claims, err := token.ValidateJWT(bearerToken, config.GetConfig().JWT.Secret)
		if err != nil {
			slog.Warn(fmt.Sprintf("delivery.handler.middleware.OptionalAuth: %v", err))
			return next(c)
		}

		fext.SetUserIDToContext(c, claims.UserID)
		analytics.SendToKafka(claims.UserID, c.OriginalURL())

		return next(c)
	}
}

func GetTokenAuth(c *fiber.Ctx, authType string) (string, error) {
	headerAuthorization, ok := c.GetReqHeaders()[router.HeaderAuthorization]
	if !ok {
		return runtime.EmptyString,
			fmt.Errorf("delivery.handler.middleware.Auth: not found header [%s]", router.HeaderAuthorization)
	}
	if len(headerAuthorization) != 1 {
		return runtime.EmptyString,
			fmt.Errorf("delivery.handler.middleware.Auth: invalid header [%s]", router.HeaderAuthorization)
	}

	authorization := headerAuthorization[0]
	if !strings.HasPrefix(authorization, authType) {
		return runtime.EmptyString,
			fmt.Errorf("delivery.handler.middleware.Auth: invalid type auth [%s] for token [%s]", authType, authorization)
	}

	tokenData := strings.Split(authorization, " ")
	if len(tokenData) != 2 {
		return runtime.EmptyString,
			fmt.Errorf("delivery.handler.middleware.Auth: invalid token [%s] for type auth [%s]", authorization, router.AuthTypeBearer)
	}

	return tokenData[1], nil
}
