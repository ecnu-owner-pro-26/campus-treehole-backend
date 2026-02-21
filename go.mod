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
	
	// 配置管理
	github.com/joho/godotenv v1.5.1
	
	// 数据验证
	github.com/go-playground/validator/v10 v10.16.0
	
	// 日志
	github.com/sirupsen/logrus v1.9.3
	
	// 测试
	github.com/stretchr/testify v1.8.4
)
