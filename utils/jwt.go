package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT载荷结构
type Claims struct {
	UserID int64  `json:"user_id"` // 用户ID
	OpenID string `json:"openid"`  // 微信OpenID
	jwt.RegisteredClaims
}

// GetJWTSecret 从环境变量获取JWT密钥
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// 开发环境默认密钥（生产环境必须设置环境变量）
		secret = "campus-memory-default-secret-key"
	}
	return []byte(secret)
}

// GetJWTExpireHours 从环境变量获取JWT过期时间（小时）
func GetJWTExpireHours() int {
	hoursStr := os.Getenv("JWT_EXPIRE_HOURS")
	if hoursStr == "" {
		return 168 // 默认7天
	}
	
	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		return 168 // 解析失败使用默认值
	}
	
	return hours
}

// GenerateToken 生成JWT token
func GenerateToken(userID int64, openid string) (string, error) {
	// 设置过期时间
	expireHours := GetJWTExpireHours()
	expirationTime := time.Now().Add(time.Duration(expireHours) * time.Hour)
	
	// 创建声明
	claims := &Claims{
		UserID: userID,
		OpenID: openid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	
	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 签名并返回
	tokenString, err := token.SignedString(GetJWTSecret())
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// ParseToken 解析和验证JWT token
func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// 验证token并提取claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, jwt.ErrSignatureInvalid
}
