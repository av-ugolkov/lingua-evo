package gin_extension

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/av-ugolkov/lingua-evo/runtime"
)

const (
	cookiePathAuth = "/auth"
)

var (
	errNotFound = errors.New("not found query param")
)

func GetQuery(c *gin.Context, key string) (string, error) {
	value, ok := c.GetQuery(key)
	if !ok {
		return "", fmt.Errorf("gin_extension.GetQuery - %w [%s]", errNotFound, key)
	}

	return value, nil
}

func GetQueryUUID(c *gin.Context, key string) (uuid.UUID, error) {
	value, ok := c.GetQuery(key)
	if !ok {
		return uuid.Nil, fmt.Errorf("gin_extension.GetQueryUUID - %w [%s]", errNotFound, key)
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("gin_extension.GetQueryUUID - parse: %w", err)
	}
	return id, nil
}

func GetQueryInt(c *gin.Context, key string) (int, error) {
	value, ok := c.GetQuery(key)
	if !ok {
		return 0, fmt.Errorf("gin_extension.GetQueryInt - %w [%s]", errNotFound, key)
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("gin_extension.GetQueryInt - parse: %w", err)
	}
	return v, nil
}

func SetCookieRefreshToken(c *gin.Context, token uuid.UUID, maxAge time.Duration) {
	c.SetCookie(RefreshToken,
		token.String(),
		int(maxAge.Seconds()),
		cookiePathAuth,
		"",
		true,
		true)
}

func DeleteCookie(c *gin.Context, name string) {
	c.SetCookie(name, runtime.EmptyString, -1, "", "", true, true)
}

func GetHeaderAuthorization(c *gin.Context, typeAuth string) (string, error) {
	token := c.GetHeader("Authorization")
	if token == runtime.EmptyString {
		return runtime.EmptyString, fmt.Errorf("gin_extension.GetHeaderAuthorization: not found Authorization token")
	}

	if !strings.HasPrefix(token, string(typeAuth)) {
		return runtime.EmptyString, fmt.Errorf("gin_extension.GetHeaderAuthorization - invalid type auth [%s]: %s", typeAuth, token)
	}

	return token[len(string(typeAuth))+1:], nil
}

func SendError(c *gin.Context, httpStatus int, err error) {
	slog.Error(err.Error())
	c.JSON(httpStatus, gin.H{
		"error": err.Error()})
}

func GetCookieLanguageOrDefault(c *gin.Context) string {
	cookie, err := c.Cookie(Language)
	switch {
	case errors.Is(err, http.ErrNoCookie):
		return runtime.GetLanguage("en")
	case err != nil:
		slog.Error(fmt.Errorf("http.exchange.Exchanger.GetCookieLanguageOrDefault: %w", err).Error())
		return runtime.GetLanguage("en")
	default:
		return cookie
	}
}
