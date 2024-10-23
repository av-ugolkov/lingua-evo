package ginext

import (
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
		if status, resp, err := h(c); err != nil {
			c.SendError(status, err)
		} else {
			c.JSON(status, resp)
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
