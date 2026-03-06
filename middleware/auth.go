package middleware

import (
	"campus-memory/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从Header获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "未提供认证token",
			})
			c.Abort()
			return
		}

		// 2. 验证Bearer格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "认证格式错误，应为: Bearer {token}",
			})
			c.Abort()
			return
		}

		// 3. 提取token字符串
		tokenString := parts[1]

		// 4. 解析和验证token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "token无效或已过期",
			})
			c.Abort()
			return
		}

		// 5. 将用户信息存入context
		c.Set("user_id", claims.UserID)
		c.Set("openid", claims.OpenID)

		// 6. 继续处理请求
		c.Next()
	}
}

// OptionalAuth 可选认证中间件（有token则验证，无token则跳过）
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		
		// 如果没有提供token，直接跳过
		if authHeader == "" {
			c.Next()
			return
		}

		// 如果提供了token，则验证
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]
			claims, err := utils.ParseToken(tokenString)
			if err == nil {
				// token有效，存入context
				c.Set("user_id", claims.UserID)
				c.Set("openid", claims.OpenID)
			}
		}

		c.Next()
	}
}
