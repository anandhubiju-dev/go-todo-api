package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		fmt.Println("Request started",
			c.Request.Method,
			c.Request.URL.Path)

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		fmt.Println("Request completed in:",
			c.Request.Method,
			c.Request.URL.Path,
			status,
			"in",
			duration,
		)
	}
}
