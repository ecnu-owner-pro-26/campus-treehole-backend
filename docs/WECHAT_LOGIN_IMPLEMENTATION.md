# 微信登录功能实现说明

## 已完成的文件

### 1. 工具层
- ✅ `utils/jwt.go` - JWT token生成和解析
- ✅ `utils/wechat.go` - 微信API调用
- ✅ `infra/util/response.go` - 统一响应格式

### 2. 数据层
- ✅ `infra/repo/user_repo.go` - 用户数据访问

### 3. 业务层
- ✅ `application/service/auth_service.go` - 认证业务逻辑

### 4. 接口层
- ✅ `middleware/auth.go` - JWT认证中间件
- ✅ `api/handler/auth_handler.go` - HTTP处理器

### 5. 测试
- ✅ `scripts/test_api.http` - API测试用例

---

## 环境变量配置

在 `.env` 文件中添加以下配置：

```env
# 微信小程序配置
WECHAT_APPID=your_wechat_appid
WECHAT_SECRET=your_wechat_secret

# JWT配置
JWT_SECRET=your_jwt_secret_key_change_in_production
JWT_EXPIRE_HOURS=168
```

---

## 路由配置

需要在 `api/router/router.go` 中添加路由：

```go
// 创建handler实例
authHandler := handler.NewAuthHandler(authService)

// 公开路由
auth := api.Group("/auth")
{
    auth.POST("/wechat/login", authHandler.WechatLogin)
}

// 需要认证的路由
authenticated := api.Group("")
authenticated.Use(middleware.JWTAuth())
{
    authRoutes := authenticated.Group("/auth")
    {
        authRoutes.GET("/profile", authHandler.GetProfile)
        authRoutes.PUT("/profile", authHandler.UpdateProfile)
    }
}
```

---

## 依赖注入

在 `main.go` 中初始化服务：

```go
// 初始化数据库
db, err := infra.InitDatabase("data/campus_memory.db")
if err != nil {
    log.Fatalf("Failed to init database: %v", err)
}

// 创建repo
userRepo := repo.NewUserRepo(db.DB)

// 创建service
authService := service.NewAuthService(userRepo)

// 创建handler
authHandler := handler.NewAuthHandler(authService)

// 设置路由
r := gin.Default()
router.SetupRoutes(r, authHandler)
```

---

## API使用示例

### 1. 微信登录

**请求：**
```http
POST /api/auth/wechat/login
Content-Type: application/json

{
  "code": "微信wx.login()返回的code",
  "nickname": "用户昵称",
  "avatar": "头像URL"
}
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "nickname": "用户昵称",
      "avatar": "头像URL",
      "default_campus_id": null
    }
  }
}
```

### 2. 获取个人信息

**请求：**
```http
GET /api/auth/profile
Authorization: Bearer {token}
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "nickname": "用户昵称",
    "avatar": "头像URL",
    "default_campus_id": null,
    "status": 1,
    "role": 0,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 3. 更新个人信息

**请求：**
```http
PUT /api/auth/profile
Authorization: Bearer {token}
Content-Type: application/json

{
  "nickname": "新昵称",
  "avatar": "新头像URL",
  "default_campus_id": 1
}
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "更新成功"
  }
}
```

---

## 功能特点

1. **自动创建用户**：首次登录自动创建用户记录
2. **信息更新**：再次登录时更新昵称和头像
3. **JWT认证**：使用JWT token进行身份验证
4. **UnionID支持**：支持微信UnionID跨应用识别
5. **可选认证**：提供OptionalAuth中间件用于可选认证场景

---

## 测试步骤

1. 配置环境变量（`.env`文件）
2. 启动服务器：`go run main.go`
3. 使用 `scripts/test_api.http` 测试接口
4. 检查数据库中的用户记录

---

## 注意事项

1. **生产环境**：必须修改 `JWT_SECRET` 为强密码
2. **微信配置**：需要在微信公众平台获取真实的 AppID 和 Secret
3. **测试环境**：可以使用模拟的 code 进行测试
4. **OpenID安全**：OpenID 不会返回给前端（在模型中标记为 `json:"-"`）

---

## 下一步

1. 在 `router.go` 中配置路由
2. 在 `main.go` 中初始化依赖注入
3. 配置 `.env` 文件
4. 运行测试验证功能

---

## 相关文档

- [API文档](./API.md)
- [开发指南](./DEVELOPMENT.md)
- [项目变更记录](./CHANGES.md)
