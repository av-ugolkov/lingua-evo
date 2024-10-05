package app

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

func testHandle(r *gin.Engine) {
	r.GET("/test", testGet)
	r.GET("/clear", testClear)
}

var count atomic.Int32

func testGet(c *gin.Context) {
	time.Sleep(5 * time.Second)

	count.Add(1)

	c.JSON(http.StatusOK, gin.H{"count": count.Load()})
}

func testClear(c *gin.Context) {
	count.Add(-count.Load())

	c.JSON(http.StatusOK, gin.H{"count": count.Load()})
}
