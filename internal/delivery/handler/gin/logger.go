package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{Skip: func(c *gin.Context) bool {
		return c.Request.Method == http.MethodOptions
	}})

}
