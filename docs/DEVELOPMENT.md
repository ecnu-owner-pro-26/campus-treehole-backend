# 校园记忆树洞后端 - 开发指南

## 📋 目录

- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [开发环境搭建](#开发环境搭建)
- [开发流程](#开发流程)
- [API 测试](#api-测试)
- [代码规范](#代码规范)
- [常见问题](#常见问题)
- [部署指南](#部署指南)

## 🚀 快速开始

### 前置要求

- Go 1.21 或更高版本
- Git
- SQLite3
- 推荐使用 VS Code 或 GoLand

### 一键设置

```bash
# 克隆项目
git clone <repository-url>
cd campus-treehole-backend

# 运行设置脚本
chmod +x scripts/dev_setup.sh
./scripts/dev_setup.sh

# 启动项目
make run
```

### 手动设置

```bash
# 1. 安装依赖
make deps

# 2. 复制环境变量文件
cp .env.example .env

# 3. 初始化数据库
make init-db

# 4. 运行项目
make run
```

## 📁 项目结构

```
campus-treehole-backend/
├── api/                    # API 层
│   ├── handler/           # HTTP 处理器
│   ├── router/            # 路由配置
│   └── token/             # JWT 相关
├── application/           # 应用层
│   ├── assembler/         # 数据组装器
│   ├── dto/               # 数据传输对象
│   └── service/           # 业务服务
├── infra/                 # 基础设施层
│   ├── cache/             # 缓存
│   ├── model/             # 数据模型
│   ├── repo/              # 数据仓库
│   └── util/              # 工具类
├── middleware/            # 中间件
├── types/                 # 类型定义
├── utils/                 # 通用工具
├── scripts/               # 脚本文件
├── docs/                  # 文档
├── data/                  # 数据库文件
├── logs/                  # 日志文件
├── storage/               # 文件存储
└── build/                 # 构建输出
```

### 架构说明

项目采用分层架构设计：

- **API 层**: 处理 HTTP 请求和响应
- **应用层**: 业务逻辑处理
- **基础设施层**: 数据访问和外部服务
- **中间件**: 横切关注点（认证、日志等）

## 🛠️ 开发环境搭建

### 1. 环境要求

```bash
# 检查 Go 版本
go version  # 需要 1.21+

# 检查 Git
git --version

# 检查 SQLite
sqlite3 --version
```

### 2. 依赖管理

```bash
# 查看所有依赖
go list -m all

# 更新依赖
go get -u ./...

# 清理依赖
go mod tidy
```

### 3. 开发工具安装

```bash
# 热重载工具
go install github.com/cosmtrek/air@latest

# 代码检查工具
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# API 文档生成工具
go install github.com/swaggo/swag/cmd/swag@latest
```

### 4. 环境变量配置

复制 `.env.example` 到 `.env` 并修改相应配置：

```bash
cp .env.example .env
```

重要配置项：

- `SERVER_PORT`: 服务器端口（默认 8080）
- `DB_PATH`: 数据库文件路径
- `JWT_SECRET`: JWT 签名密钥（生产环境必须修改）
- `LOG_LEVEL`: 日志级别

## 🔄 开发流程

### 1. 分支管理

```bash
# 创建功能分支
git checkout -b feature/your-feature-name

# 提交代码
git add .
git commit -m "feat: add your feature description"

# 推送分支
git push origin feature/your-feature-name
```

### 2. 代码开发

```bash
# 开发模式运行（热重载）
make dev

# 普通模式运行
make run

# 代码格式化
make fmt

# 代码检查
make lint

# 运行测试
make test
```

### 3. 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

示例：
```bash
git commit -m "feat(auth): add user registration endpoint"
git commit -m "fix(memory): resolve memory creation validation issue"
git commit -m "docs: update API documentation"
```

## 🧪 API 测试

### 1. 使用 HTTP 文件测试

项目提供了完整的 API 测试文件：`scripts/test_api.http`

在 VS Code 中安装 REST Client 插件，然后打开测试文件进行测试。

### 2. 测试流程

```bash
# 1. 启动服务
make run

# 2. 运行用户注册
POST http://localhost:8080/api/auth/register

# 3. 运行用户登录获取 token
POST http://localhost:8080/api/auth/login

# 4. 使用 token 测试其他接口
```

### 3. 自动化测试

```bash
# 运行单元测试
make test

# 运行集成测试
go test -tags=integration ./...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📝 代码规范

### 1. 命名规范

- **包名**: 小写，简短，有意义
- **文件名**: 小写，下划线分隔
- **函数名**: 驼峰命名，公开函数首字母大写
- **变量名**: 驼峰命名，局部变量首字母小写
- **常量名**: 全大写，下划线分隔

### 2. 注释规范

```go
// Package handler 提供 HTTP 请求处理器
package handler

// UserHandler 用户相关的 HTTP 处理器
type UserHandler struct {
    userService *service.UserService
}

// CreateUser 创建新用户
// @Summary 创建用户
// @Description 创建一个新的用户账户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "用户信息"
// @Success 200 {object} dto.CreateUserResponse
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // 实现逻辑
}
```

### 3. 错误处理

```go
// 统一错误处理
func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
    // 参数验证
    if err := s.validator.Validate(req); err != nil {
        return nil, errno.ErrInvalidParams.WithDetails(err.Error())
    }
    
    // 业务逻辑
    user, err := s.userRepo.Create(req)
    if err != nil {
        return nil, errno.ErrUserCreateFail.WithError(err)
    }
    
    return &dto.CreateUserResponse{User: user}, nil
}
```

### 4. 数据库操作

```go
// 使用事务
func (r *UserRepo) CreateUserWithProfile(user *model.User, profile *model.Profile) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(user).Error; err != nil {
            return err
        }
        
        profile.UserID = user.ID
        if err := tx.Create(profile).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

## ❓ 常见问题

### 1. 编译问题

**问题**: `go build` 失败
```bash
# 解决方案
go mod tidy
go mod download
```

**问题**: 依赖版本冲突
```bash
# 解决方案
go mod edit -replace=problematic-package@version=./local-path
go mod tidy
```

### 2. 数据库问题

**问题**: 数据库连接失败
```bash
# 检查数据库文件权限
ls -la data/campus_memory.db

# 重新初始化数据库
rm data/campus_memory.db
make init-db
```

**问题**: 数据库迁移失败
```bash
# 手动执行 SQL 脚本
sqlite3 data/campus_memory.db < scripts/init_db.sql
```

### 3. 运行时问题

**问题**: 端口被占用
```bash
# 查找占用端口的进程
lsof -i :8080

# 修改端口配置
# 编辑 .env 文件中的 SERVER_PORT
```

**问题**: JWT Token 验证失败
```bash
# 检查 JWT_SECRET 配置
# 确保客户端和服务端使用相同的密钥
```

### 4. 开发工具问题

**问题**: Air 热重载不工作
```bash
# 检查 Air 配置
cat .air.toml

# 重新安装 Air
go install github.com/cosmtrek/air@latest
```

## 🚀 部署指南

### 1. 构建生产版本

```bash
# 构建二进制文件
make build

# 构建 Docker 镜像
docker build -t campus-memory:latest .
```

### 2. 环境配置

生产环境需要修改的配置：

```bash
# .env 文件
APP_ENV=production
DEBUG=false
JWT_SECRET=your-production-secret-key
LOG_LEVEL=warn
```

### 3. 数据库迁移

```bash
# 备份数据库
cp data/campus_memory.db data/campus_memory.db.backup

# 执行迁移脚本
sqlite3 data/campus_memory.db < scripts/migrate.sql
```

### 4. 服务部署

```bash
# 使用 systemd (Linux)
sudo cp scripts/campus-memory.service /etc/systemd/system/
sudo systemctl enable campus-memory
sudo systemctl start campus-memory

# 使用 Docker
docker run -d \
  --name campus-memory \
  -p 8080:8080 \
  -v ./data:/app/data \
  -v ./logs:/app/logs \
  campus-memory:latest

# 使用 Docker Compose
docker-compose up -d
```

### 5. 监控和日志

```bash
# 查看应用日志
tail -f logs/app.log

# 查看系统服务状态
sudo systemctl status campus-memory

# 监控资源使用
htop
```

## 📚 相关资源

- [Go 官方文档](https://golang.org/doc/)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [JWT 规范](https://jwt.io/)
- [RESTful API 设计指南](https://restfulapi.net/)

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交代码
4. 创建 Pull Request
5. 等待代码审查

## 📄 许可证

本项目采用 MIT 许可证，详见 [LICENSE](../LICENSE) 文件。