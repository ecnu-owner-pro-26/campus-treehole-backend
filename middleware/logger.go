package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 请求日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 请求路径
		path := c.Request.URL.Path
		
		// 请求方法
		method := c.Request.Method
		
		// 客户端IP
		clientIP := c.ClientIP()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		
		// 执行时间
		latency := endTime.Sub(startTime)
		
		// 状态码
		statusCode := c.Writer.Status()
		
		// 错误信息
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 记录日志
		log.Printf("[GIN] %s | %3d | %13v | %15s | %-7s %s %s",
			endTime.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			errorMessage,
		)
	}
}

// LoggerWithConfig 带配置的日志中间件
func LoggerWithConfig(skipPaths []string) gin.HandlerFunc {
	// 创建跳过路径的映射
	skipPathsMap := make(map[string]bool, len(skipPaths))
	for _, path := range skipPaths {
		skipPathsMap[path] = true
	}

	return func(c *gin.Context) {
		// 检查是否跳过该路径
		path := c.Request.URL.Path
		if skipPathsMap[path] {
			c.Next()
			return
		}

		// 开始时间
		startTime := time.Now()

		// 请求方法
		method := c.Request.Method
		
		// 客户端IP
		clientIP := c.ClientIP()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		
		// 执行时间
		latency := endTime.Sub(startTime)
		
		// 状态码
		statusCode := c.Writer.Status()
		
		// 错误信息
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 记录日志
		log.Printf("[GIN] %s | %3d | %13v | %15s | %-7s %s %s",
			endTime.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			errorMessage,
		)
	}
}
