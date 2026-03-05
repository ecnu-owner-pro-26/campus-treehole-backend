package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		log.Printf("[%s] %s %s %d %v %s",
			reqMethod,
			reqUri,
			clientIP,
			statusCode,
			latencyTime,
			c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)
	}
}