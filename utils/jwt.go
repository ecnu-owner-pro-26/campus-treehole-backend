package utils

import (
	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT载荷结构
type Claims struct {
	// TODO: 定义JWT载荷字段
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken() (string, error) {
	// TODO: 实现生成JWT token逻辑
	return "", nil
}

// ParseToken 解析和验证JWT token
func ParseToken(tokenString string) (*Claims, error) {
	// TODO: 实现解析和验证JWT token逻辑
	return nil, nil
}
