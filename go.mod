module campus-memory

go 1.21

require (
	// Web框架
	github.com/gin-gonic/gin v1.9.1
	
	// 数据库ORM
	gorm.io/gorm v1.25.5
	gorm.io/driver/sqlite v1.5.4
	
	// JWT认证
	github.com/golang-jwt/jwt/v5 v5.2.0
	
	// 跨域支持
	github.com/gin-contrib/cors v1.4.0
	
	// 密码加密
	golang.org/x/crypto v0.17.0
)
