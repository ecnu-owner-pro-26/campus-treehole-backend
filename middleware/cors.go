package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的来源（生产环境应该配置具体的域名）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		
		// 允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		
		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		
		// 允许携带凭证（cookies）
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		
		// 预检请求的缓存时间（秒）
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		
		// 允许客户端访问的响应头
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		// 处理 OPTIONS 预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// 继续处理请求
		c.Next()
	}
}
