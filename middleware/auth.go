package middleware

import (
	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件（必需认证）
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现JWT认证逻辑
		// 1. 从Header获取Authorization
		// 2. 验证Bearer格式
		// 3. 解析和验证token
		// 4. 将用户信息存入context
		c.Next()
	}
}

// OptionalAuth 可选认证中间件（有token则验证，无token也放行）
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现可选认证逻辑
		c.Next()
	}
}
