package ginext

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	msgerror "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	//TODO need to be refactored
	cookiePathAuth = "/auth"
)

var (
	e = msgerror.ApiError{}

	errNotFound = errors.New("not found query param")
)

type Context struct {
	*gin.Context
}

func NewContext(c *gin.Context) *Context {
	return &Context{
		Context: c,
	}
}

func (c *Context) GetQuery(key string) (string, error) {
	value, ok := c.Context.GetQuery(key)
	if !ok {
		return "", fmt.Errorf("gin-ext.Context.GetQuery - %w [%s]", errNotFound, key)
	}

	return value, nil
}

func (c *Context) GetQueryUUID(key string) (uuid.UUID, error) {
	value, ok := c.Context.GetQuery(key)
	if !ok {
		return uuid.Nil, fmt.Errorf("gin_extension.GetQueryUUID - %w [%s]", errNotFound, key)
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("gin_extension.GetQueryUUID - parse: %w", err)
	}
	return id, nil
}

func (c *Context) GetQueryInt(key string) (int, error) {
	value, ok := c.Context.GetQuery(key)
	if !ok {
		return 0, fmt.Errorf("gin_extension.GetQueryInt - %w [%s]", errNotFound, key)
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("gin_extension.GetQueryInt - parse: %w", err)
	}
	return v, nil
}

func (c *Context) SetCookieRefreshToken(token uuid.UUID, maxAge time.Duration) {
	c.Context.SetCookie(RefreshToken,
		token.String(),
		int(maxAge.Seconds()),
		cookiePathAuth,
		"",
		true,
		true)
}

func (c *Context) DeleteCookie(name string) {
	c.Context.SetCookie(name, runtime.EmptyString, -1, "", "", true, true)
}

func (c *Context) GetHeaderAuthorization(typeAuth string) (string, error) {
	token := c.Context.GetHeader("Authorization")
	if token == runtime.EmptyString {
		return runtime.EmptyString, fmt.Errorf("gin_extension.GetHeaderAuthorization: not found Authorization token")
	}

	if !strings.HasPrefix(token, string(typeAuth)) {
		return runtime.EmptyString, fmt.Errorf("gin_extension.GetHeaderAuthorization - invalid type auth [%s]: %s", typeAuth, token)
	}

	tokenData := strings.Split(token, " ")
	if len(tokenData) != 2 {
		return runtime.EmptyString, fmt.Errorf("gin_extension.GetHeaderAuthorization - invalid token: %s", token)
	}

	return tokenData[1], nil
}

func (c *Context) SendError(httpStatus int, err error) {
	slog.Error(err.Error())
	switch {
	case errors.Is(err, &e):
		c.Context.JSON(httpStatus, e.Msg)
	default:
		c.Context.JSON(httpStatus, err.Error())
	}
}

func (c *Context) GetCookieLanguageOrDefault() string {
	cookie, err := c.Context.Cookie(Language)
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
