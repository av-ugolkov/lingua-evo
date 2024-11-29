package ginext

import (
	"errors"
	"log/slog"
	"net/http"

	msgerr "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *Context) (status int, obj any, err error)

type Engine struct {
	*gin.Engine
}

func NewEngine(eng *gin.Engine) *Engine {
	eng.UseH2C = true
	return &Engine{
		Engine: eng,
	}
}

func execHandlers(gc *gin.Context, handlers ...HandlerFunc) {
	for _, h := range handlers {
		c := NewContext(gc)
		status, resp, err := h(c)
		switch status {
		case http.StatusPermanentRedirect, http.StatusTemporaryRedirect:
			c.Redirect(status, resp.(string))
		default:
			obj := gin.H{"resp": resp}
			if err != nil {
				slog.Error(err.Error())
				var e *msgerr.ApiError
				switch {
				case errors.As(err, &e):
					obj["msg"] = e.Msg
				default:
					obj["msg"] = err.Error()
				}
			}
			c.JSON(status, obj)
		}

	}
}

func (e *Engine) GET(relativePath string, handlers ...HandlerFunc) {
	e.Engine.GET(relativePath, func(c *gin.Context) {
		execHandlers(c, handlers...)
	})
}

func (e *Engine) POST(relativePath string, handlers ...HandlerFunc) {
	e.Engine.POST(relativePath, func(c *gin.Context) {
		execHandlers(c, handlers...)
	})
}

func (e *Engine) DELETE(relativePath string, handlers ...HandlerFunc) {
	e.Engine.DELETE(relativePath, func(c *gin.Context) {
		execHandlers(c, handlers...)
	})
}

func (e *Engine) PUT(relativePath string, handlers ...HandlerFunc) {
	e.Engine.PUT(relativePath, func(c *gin.Context) {
		execHandlers(c, handlers...)
	})
}

func (e *Engine) PATCH(relativePath string, handlers ...HandlerFunc) {
	e.Engine.PATCH(relativePath, func(c *gin.Context) {
		execHandlers(c, handlers...)
	})
}
