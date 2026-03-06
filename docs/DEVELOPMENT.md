# 校园记忆后端 - 开发指南

## 📋 目录

- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [技术栈](#技术栈)
- [开发环境搭建](#开发环境搭建)
- [开发流程](#开发流程)
- [API 测试](#api-测试)
- [代码规范](#代码规范)
- [数据库管理](#数据库管理)
- [常见问题](#常见问题)
- [部署指南](#部署指南)

## 🚀 快速开始

### 前置要求

- **Go 1.24.0** 或更高版本
- **Git** 版本控制
- **SQLite3** 数据库
- **Make** 构建工具
- 推荐 IDE: **VS Code** 或 **GoLand**

### 一键设置（Linux/Mac）

```bash
# 1. 克隆项目
git clone <repository-url>
cd campus-memory

# 2. 运行设置脚本
chmod +x scripts/dev_setup.sh
./scripts/dev_setup.sh

# 3. 启动项目
make run
```

### 一键设置（Windows）

```bash
# 1. 克隆项目
git clone <repository-url>
cd campus-memory

# 2. 运行设置脚本
scripts\dev_setup.bat

# 3. 启动项目
make run
```

### 手动设置

```bash
# 1. 安装依赖
make deps

# 2. 复制并配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置微信小程序 AppID 和 Secret

# 3. 初始化数据库
make init-db

# 4. 运行项目
make run
```

### 验证安装

```bash
# 检查服务是否启动
curl http://localhost:8080/api/campuses

# 或在浏览器访问
# http://localhost:8080/api/campuses
```

## 📁 项目结构

```
campus-memory/
├── api/                    # API 层
│   ├── handler/           # HTTP 请求处理器
│   │   ├── auth_handler.go          # 认证相关
│   │   ├── campus_handler.go        # 校区相关
│   │   ├── location_handler.go      # 地点相关
│   │   ├── memory_handler.go        # 记忆相关
│   │   ├── comment_handler.go       # 评论相关
│   │   ├── like_handler.go          # 点赞相关
│   │   ├── image_handler.go         # 图片上传相关
│   │   └── quicknav_handler.go      # 快速导航相关
│   └── router/            # 路由配置
│       └── router.go                # 路由注册
│
├── application/           # 应用层
│   ├── assembler/         # 数据组装器（Model ↔ DTO 转换）
│   │   ├── campus_assembler.go
│   │   ├── comment_assembler.go
│   │   ├── location_assembler.go
│   │   ├── memory_assembler.go
│   │   └── user_assembler.go
│   ├── dto/               # 数据传输对象
│   │   ├── auth_dto.go              # 认证相关 DTO
│   │   ├── campus_dto.go            # 校区相关 DTO
│   │   ├── location_dto.go          # 地点相关 DTO
│   │   ├── memory_dto.go            # 记忆相关 DTO
│   │   ├── comment_dto.go           # 评论相关 DTO
│   │   ├── like_dto.go              # 点赞相关 DTO
│   │   ├── image_dto.go             # 图片相关 DTO
│   │   └── quicknav_dto.go          # 快速导航 DTO
│   └── service/           # 业务服务层
│       ├── auth_service.go          # 认证业务逻辑
│       ├── campus_service.go        # 校区业务逻辑
│       ├── location_service.go      # 地点业务逻辑
│       ├── memory_service.go        # 记忆业务逻辑
│       ├── comment_service.go       # 评论业务逻辑
│       ├── like_service.go          # 点赞业务逻辑
│       ├── image_service.go         # 图片业务逻辑
│       └── quicknav_service.go      # 快速导航业务逻辑
│
├── infra/                 # 基础设施层
│   ├── cache/             # 缓存（Redis）
│   │   └── redis.go
│   ├── model/             # 数据库模型
│   │   ├── campus_model.go
│   │   ├── location_model.go
│   │   ├── memory_model.go
│   │   ├── comment_model.go
│   │   ├── like_model.go
│   │   ├── image_model.go
│   │   └── user_model.go
│   ├── repo/              # 数据仓库（数据访问层）
│   │   ├── campus_repo.go
│   │   ├── location_repo.go
│   │   ├── memory_repo.go
│   │   ├── comment_repo.go
│   │   ├── like_repo.go
│   │   ├── image_repo.go
│   │   └── user_repo.go
│   ├── util/              # 基础工具
│   │   └── response.go              # 统一响应格式
│   ├── database.go        # 数据库初始化
│   └── init.go            # 基础设施初始化
│
├── middleware/            # 中间件
│   ├── auth.go            # JWT 认证中间件
│   ├── cors.go            # 跨域处理中间件
│   └── logger.go          # 日志记录中间件
│
├── types/                 # 类型定义
│   └── errno/             # 错误码定义
│       └── error.go
│
├── utils/                 # 通用工具
│   ├── jwt.go             # JWT 工具
│   └── wechat.go          # 微信 API 工具
│
├── scripts/               # 脚本文件
│   ├── init_db.sql                  # 数据库初始化脚本
│   ├── init_campus_data.sql         # 测试数据脚本
│   ├── test_api.http                # API 测试文件
│   ├── dev_setup.sh                 # Linux/Mac 开发环境设置
│   ├── dev_setup.bat                # Windows 开发环境设置
│   └── build.sh                     # 构建脚本
│
├── docs/                  # 文档
│   ├── API.md             # API 接口文档
│   └── DEVELOPMENT.md     # 开发指南（本文件）
│
├── data/                  # 数据库文件目录
│   └── campus_memory.db   # SQLite 数据库文件
│
├── uploads/               # 文件上传目录
│   └── images/            # 图片存储
│
├── build/                 # 构建输出目录
│
├── main.go                # 应用程序入口
├── Makefile               # Make 构建配置
├── go.mod                 # Go 模块依赖
├── go.sum                 # 依赖校验文件
├── .env.example           # 环境变量示例
├── .gitignore             # Git 忽略配置
├── LICENSE                # 开源许可证
└── README.md              # 项目说明
```

### 架构说明

项目采用 **DDD（领域驱动设计）分层架构**：

#### 1. API 层（api/）
- **职责**: 处理 HTTP 请求和响应
- **组件**:
  - `handler`: 接收请求，调用 service，返回响应
  - `router`: 路由注册和中间件配置

#### 2. 应用层（application/）
- **职责**: 业务逻辑处理和编排
- **组件**:
  - `service`: 核心业务逻辑
  - `dto`: 前后端数据传输对象
  - `assembler`: Model 和 DTO 之间的转换

#### 3. 基础设施层（infra/）
- **职责**: 数据持久化和外部服务
- **组件**:
  - `model`: 数据库实体模型
  - `repo`: 数据访问接口实现
  - `cache`: 缓存实现
  - `util`: 基础工具函数

#### 4. 中间件层（middleware/）
- **职责**: 横切关注点
- **组件**:
  - JWT 认证
  - CORS 跨域处理
  - 日志记录

#### 5. 工具层（utils/、types/）
- **职责**: 通用工具和类型定义
- **组件**:
  - JWT 工具
  - 微信 API 工具
  - 错误码定义

### 数据流向

```
HTTP 请求
    ↓
Router（路由）
    ↓
Middleware（中间件：认证、日志等）
    ↓
Handler（处理器：参数验证）
    ↓
Service（业务逻辑）
    ↓
Repository（数据访问）
    ↓
Database（数据库）
```

## � 技术栈

### 后端框架
- **Gin** v1.11.0 - 高性能 Web 框架
- **GORM** v1.31.1 - ORM 框架
- **SQLite** v1.6.0 - 嵌入式数据库

### 认证与安全
- **JWT** v5.3.1 - JSON Web Token 认证
- **bcrypt** - 密码加密（golang.org/x/crypto）

### 数据验证
- **validator** v10.27.0 - 请求参数验证

### 可选组件
- **Redis** v9.18.0 - 缓存（可选）

### 开发工具
- **Air** - 热重载工具
- **golangci-lint** - 代码检查工具
- **Swag** - API 文档生成工具

## 🛠️ 开发环境搭建

### 1. 环境要求

```bash
# 检查 Go 版本（需要 1.24.0+）
go version

# 检查 Git
git --version

# 检查 SQLite
sqlite3 --version

# 检查 Make
make --version
```

### 2. 依赖管理

```bash
# 安装项目依赖
make deps

# 查看所有依赖
go list -m all

# 更新依赖
go get -u ./...

# 清理未使用的依赖
go mod tidy

# 下载依赖到本地缓存
go mod download
```

### 3. 开发工具安装

```bash
# 一键安装所有开发工具
make install-tools

# 或手动安装
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

**必需配置项**：

```bash
# 微信小程序配置（必需）
WECHAT_APPID=your_wechat_appid
WECHAT_SECRET=your_wechat_secret

# JWT 配置（必需）
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRE_HOURS=168

# 服务器配置
SERVER_PORT=8080
GIN_MODE=debug

# 数据库配置（SQLite）
DB_TYPE=sqlite
DB_PATH=./data/campus_memory.db
```

**可选配置项**：

```bash
# Redis 缓存（可选）
REDIS_ENABLED=false
REDIS_HOST=localhost
REDIS_PORT=6379

# 文件上传配置
UPLOAD_DIR=./uploads
MAX_FILE_SIZE=5

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=text
LOG_FILE=./logs/app.log

# 跨域配置
CORS_ALLOWED_ORIGINS=*
CORS_ALLOW_CREDENTIALS=true
```

### 5. 数据库初始化

```bash
# 初始化数据库和测试数据
make init-db

# 或手动执行 SQL 脚本
sqlite3 data/campus_memory.db < scripts/init_db.sql
sqlite3 data/campus_memory.db < scripts/init_campus_data.sql
```

### 6. 验证环境

```bash
# 检查环境配置
make check-env

# 运行项目
make run

# 测试 API
curl http://localhost:8080/api/campuses
```

## 🔄 开发流程

### 1. 分支管理

```bash
# 创建功能分支
git checkout -b feature/your-feature-name

# 创建修复分支
git checkout -b fix/bug-description

# 提交代码
git add .
git commit -m "feat: add your feature description"

# 推送分支
git push origin feature/your-feature-name
```

### 2. 开发命令

```bash
# 开发模式运行
make dev

# 普通模式运行
make run

# 代码格式化
make fmt

# 代码检查
make lint

# 运行测试
make test

# 构建可执行文件
make build

# 清理构建文件
make clean

# 查看所有可用命令
make help
```

### 3. 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

| 类型 | 说明 | 示例 |
|------|------|------|
| `feat` | 新功能 | `feat(auth): add wechat login` |
| `fix` | 修复 bug | `fix(memory): resolve creation validation` |
| `docs` | 文档更新 | `docs: update API documentation` |
| `style` | 代码格式调整 | `style: format code with gofmt` |
| `refactor` | 代码重构 | `refactor(service): optimize query logic` |
| `perf` | 性能优化 | `perf(db): add index for queries` |
| `test` | 测试相关 | `test(memory): add unit tests` |
| `chore` | 构建/工具变动 | `chore: update dependencies` |

**提交示例**：
```bash
git commit -m "feat(auth): add user registration endpoint"
git commit -m "fix(memory): resolve memory creation validation issue"
git commit -m "docs: update API documentation for comments"
git commit -m "refactor(repo): optimize database query performance"
```

### 4. 开发工作流

```
1. 创建功能分支
   ↓
2. 编写代码
   ↓
3. 运行测试（make test）
   ↓
4. 代码检查（make lint）
   ↓
5. 格式化代码（make fmt）
   ↓
6. 提交代码
   ↓
7. 推送到远程
   ↓
8. 创建 Pull Request
   ↓
9. 代码审查
   ↓
10. 合并到主分支
```

## 📝 代码规范

### 1. 命名规范

#### 包名（Package）
- 小写，简短，有意义
- 不使用下划线或驼峰
```go
package handler  // ✓
package userHandler  // ✗
```

#### 文件名
- 小写，下划线分隔
```go
auth_handler.go  // ✓
AuthHandler.go   // ✗
```

#### 函数名
- 驼峰命名
- 公开函数首字母大写
- 私有函数首字母小写
```go
func CreateMemory()  // ✓ 公开函数
func validateInput() // ✓ 私有函数
```

#### 变量名
- 驼峰命名
- 局部变量首字母小写
- 全局变量首字母大写
```go
var userID int64        // ✓ 局部变量
var GlobalConfig Config // ✓ 全局变量
```

#### 常量名
- 驼峰命名或全大写下划线分隔
```go
const MaxRetryCount = 3     // ✓ 驼峰
const MAX_RETRY_COUNT = 3   // ✓ 全大写
```

### 2. 注释规范

#### 包注释
```go
// Package handler 提供 HTTP 请求处理器
// 包含认证、记忆、评论等模块的处理逻辑
package handler
```

#### 类型注释
```go
// UserHandler 用户相关的 HTTP 处理器
// 负责处理用户注册、登录、信息更新等请求
type UserHandler struct {
    userService *service.UserService
}
```

#### 函数注释
```go
// CreateMemory 创建新记忆
// 参数:
//   - c: Gin 上下文
// 返回:
//   - 成功: 200 + 记忆详情
//   - 失败: 错误码 + 错误信息
func (h *MemoryHandler) CreateMemory(c *gin.Context) {
    // 实现逻辑
}
```

### 3. 错误处理

#### 统一错误处理
```go
func (s *MemoryService) CreateMemory(req *dto.CreateMemoryRequest, userID int64) (*dto.MemoryResponse, error) {
    // 1. 参数验证
    if req.Title == "" {
        return nil, errno.ErrBadRequest
    }
    
    // 2. 业务逻辑
    memory, err := s.memoryRepo.Create(req, userID)
    if err != nil {
        return nil, errno.ErrMemoryCreateFail
    }
    
    // 3. 返回结果
    return s.assembler.ToMemoryResponse(memory), nil
}
```

#### Handler 层错误处理
```go
func (h *MemoryHandler) CreateMemory(c *gin.Context) {
    // 调用服务层
    resp, err := h.memoryService.CreateMemory(&req, userID)
    if err != nil {
        // 判断是否是自定义错误
        if e, ok := err.(*errno.Error); ok {
            util.ErrorResponse(c, e.Code, e.Message)
        } else {
            util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
        }
        return
    }
    
    // 返回成功响应
    util.SuccessResponse(c, resp)
}
```

### 4. 数据库操作

#### 基本查询
```go
func (r *MemoryRepo) GetByID(id int64) (*model.MemoryModel, error) {
    var memory model.MemoryModel
    err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&memory).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errno.ErrMemoryNotFound
        }
        return nil, err
    }
    return &memory, nil
}
```

#### 使用事务
```go
func (r *MemoryRepo) CreateWithImages(memory *model.MemoryModel, images []string) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // 1. 创建记忆
        if err := tx.Create(memory).Error; err != nil {
            return err
        }
        
        // 2. 创建图片记录
        for _, url := range images {
            image := &model.ImageModel{
                MemoryID: memory.ID,
                URL:      url,
            }
            if err := tx.Create(image).Error; err != nil {
                return err
            }
        }
        
        return nil
    })
}
```

#### 分页查询
```go
func (r *MemoryRepo) List(page, pageSize int, locationID *int64) ([]*model.MemoryModel, int64, error) {
    var memories []*model.MemoryModel
    var total int64
    
    query := r.db.Model(&model.MemoryModel{}).Where("deleted_at IS NULL")
    
    // 条件筛选
    if locationID != nil {
        query = query.Where("location_id = ?", *locationID)
    }
    
    // 获取总数
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // 分页查询
    offset := (page - 1) * pageSize
    if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&memories).Error; err != nil {
        return nil, 0, err
    }
    
    return memories, total, nil
}
```

### 5. DTO 转换

#### Model 转 DTO
```go
func (a *MemoryAssembler) ToMemoryResponse(memory *model.MemoryModel, creator *model.UserModel, isLiked bool) *dto.MemoryResponse {
    return &dto.MemoryResponse{
        ID:           memory.ID,
        Title:        memory.Title,
        Content:      memory.Content,
        LocationName: memory.LocationName,
        LocationID:   memory.LocationID,
        Creator: dto.UserSimpleInfo{
            ID:       creator.ID,
            Nickname: creator.Nickname,
            Avatar:   creator.Avatar,
        },
        LikeCount:    memory.LikeCount,
        CommentCount: memory.CommentCount,
        ViewCount:    memory.ViewCount,
        IsLiked:      isLiked,
        Tags:         parseJSONArray(memory.Tags),
        CreatedAt:    memory.CreatedAt.Format(time.RFC3339),
    }
}
```

### 6. 代码组织

#### Service 层结构
```go
type MemoryService struct {
    memoryRepo   *repo.MemoryRepo
    userRepo     *repo.UserRepo
    locationRepo *repo.LocationRepo
    likeRepo     *repo.LikeRepo
    imageRepo    *repo.ImageRepo
}

func NewMemoryService(
    memoryRepo *repo.MemoryRepo,
    userRepo *repo.UserRepo,
    locationRepo *repo.LocationRepo,
    likeRepo *repo.LikeRepo,
    imageRepo *repo.ImageRepo,
) *MemoryService {
    return &MemoryService{
        memoryRepo:   memoryRepo,
        userRepo:     userRepo,
        locationRepo: locationRepo,
        likeRepo:     likeRepo,
        imageRepo:    imageRepo,
    }
}
```

### 7. 代码检查

```bash
# 格式化代码
make fmt

# 代码检查
make lint

# 或使用 go vet
go vet ./...

# 检查未使用的变量
go run golang.org/x/tools/cmd/deadcode@latest ./...
```

## 💾 数据库管理

### 1. 数据库结构

项目使用 SQLite 作为数据库，数据库文件位于 `data/campus_memory.db`

**核心表**：
- `users` - 用户表
- `campuses` - 校区表
- `locations` - 地点表
- `memories` - 记忆表
- `comments` - 评论表
- `likes` - 点赞表
- `images` - 图片表

### 2. 数据库操作

```bash
# 初始化数据库
make init-db

# 连接数据库
sqlite3 data/campus_memory.db

# 查看所有表
.tables

# 查看表结构
.schema memories

# 查询数据
SELECT * FROM memories LIMIT 10;

# 退出
.quit
```

### 3. 数据库迁移

```bash
# 备份数据库
cp data/campus_memory.db data/campus_memory.db.backup

# 执行迁移脚本
sqlite3 data/campus_memory.db < scripts/migrate.sql

# 恢复备份（如果需要）
cp data/campus_memory.db.backup data/campus_memory.db
```

### 4. 测试数据管理

```bash
# 清空测试数据
sqlite3 data/campus_memory.db "DELETE FROM memories;"
sqlite3 data/campus_memory.db "DELETE FROM comments;"
sqlite3 data/campus_memory.db "DELETE FROM likes;"

# 重新导入测试数据
sqlite3 data/campus_memory.db < scripts/init_campus_data.sql
```

### 5. 数据库调试

在 `.env` 文件中启用 SQL 日志：

```bash
DB_DEBUG=true
```

这将在控制台输出所有 SQL 查询语句，方便调试。

## ❓ 常见问题

### 1. 编译问题

**问题**: `go build` 失败
```bash
# 解决方案
go mod tidy
go mod download
go clean -cache
```

**问题**: 依赖版本冲突
```bash
# 查看依赖树
go mod graph

# 更新特定依赖
go get -u github.com/gin-gonic/gin@latest

# 清理并重新下载
go clean -modcache
go mod download
```

### 2. 数据库问题

**问题**: 数据库连接失败
```bash
# 检查数据库文件是否存在
ls -la data/campus_memory.db

# 检查文件权限
chmod 644 data/campus_memory.db

# 重新初始化数据库
rm data/campus_memory.db
make init-db
```

**问题**: 数据库锁定
```bash
# SQLite 数据库被锁定
# 解决方案：关闭所有连接到数据库的程序
pkill -f campus-memory
rm data/campus_memory.db-journal  # 删除日志文件
```

**问题**: 表不存在
```bash
# 检查表是否存在
sqlite3 data/campus_memory.db ".tables"

# 重新执行初始化脚本
sqlite3 data/campus_memory.db < scripts/init_db.sql
```

### 3. 运行时问题

**问题**: 端口被占用
```bash
# Linux/Mac 查找占用端口的进程
lsof -i :8080

# Windows 查找占用端口的进程
netstat -ano | findstr :8080

# 修改端口配置
# 编辑 .env 文件中的 SERVER_PORT
```

**问题**: JWT Token 验证失败
```bash
# 检查 JWT_SECRET 配置
cat .env | grep JWT_SECRET

# 确保客户端和服务端使用相同的密钥
# 检查 token 是否过期（默认 7 天）
```

**问题**: 微信登录失败
```bash
# 检查微信配置
cat .env | grep WECHAT

# 确保 WECHAT_APPID 和 WECHAT_SECRET 正确
# 检查网络连接是否正常
```

### 4. 开发工具问题

**问题**: Air 热重载不工作
```bash
# 检查 Air 是否安装
which air

# 重新安装 Air
go install github.com/cosmtrek/air@latest

# 检查 Air 配置（如果有 .air.toml）
cat .air.toml
```

**问题**: Make 命令不可用
```bash
# Windows 安装 Make
choco install make

# 或使用 Go 命令代替
go run main.go  # 代替 make run
go build        # 代替 make build
```

### 5. 图片上传问题

**问题**: 图片上传失败
```bash
# 检查上传目录是否存在
ls -la uploads/images/

# 创建上传目录
mkdir -p uploads/images

# 检查目录权限
chmod 755 uploads/images
```

**问题**: 图片访问 404
```bash
# 检查静态文件服务配置
# 在 router.go 中应该有：
# r.Static("/uploads", "./uploads")

# 检查图片文件是否存在
ls uploads/images/
```

## 🚀 部署指南

### 1. 构建生产版本

```bash
# 构建二进制文件
make build

# 输出文件位于 build/campus-memory 或 build/campus-memory.exe

# 构建 Docker 镜像
docker build -t campus-memory:latest .
```

### 2. 生产环境配置

创建生产环境配置文件 `.env.production`：

```bash
# 应用配置
APP_ENV=production
GIN_MODE=release
DEBUG=false

# 服务器配置
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# 微信小程序配置（必需）
WECHAT_APPID=your_production_appid
WECHAT_SECRET=your_production_secret

# JWT 配置（必需修改）
JWT_SECRET=your-production-secret-key-min-32-characters
JWT_EXPIRE_HOURS=168

# 数据库配置
DB_TYPE=sqlite
DB_PATH=./data/campus_memory.db

# 日志配置
LOG_LEVEL=warn
LOG_FORMAT=json
LOG_FILE=./logs/app.log

# 文件上传
UPLOAD_DIR=./uploads
MAX_FILE_SIZE=5

# 跨域配置
CORS_ALLOWED_ORIGINS=https://your-domain.com
CORS_ALLOW_CREDENTIALS=true
```

### 3. 数据库准备

```bash
# 备份开发数据库
cp data/campus_memory.db data/campus_memory.db.dev

# 清空测试数据（生产环境）
sqlite3 data/campus_memory.db "DELETE FROM memories;"
sqlite3 data/campus_memory.db "DELETE FROM comments;"
sqlite3 data/campus_memory.db "DELETE FROM likes;"
sqlite3 data/campus_memory.db "DELETE FROM users;"

# 或重新初始化
rm data/campus_memory.db
sqlite3 data/campus_memory.db < scripts/init_db.sql
sqlite3 data/campus_memory.db < scripts/init_campus_data.sql
```

### 4. 使用 Systemd 部署（Linux）

创建服务文件 `/etc/systemd/system/campus-memory.service`：

```ini
[Unit]
Description=Campus Memory Backend Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/campus-memory
ExecStart=/opt/campus-memory/campus-memory
Restart=on-failure
RestartSec=5s

# 环境变量
Environment="GIN_MODE=release"
EnvironmentFile=/opt/campus-memory/.env

# 日志
StandardOutput=append:/var/log/campus-memory/app.log
StandardError=append:/var/log/campus-memory/error.log

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
# 重新加载 systemd 配置
sudo systemctl daemon-reload

# 启用服务（开机自启）
sudo systemctl enable campus-memory

# 启动服务
sudo systemctl start campus-memory

# 查看状态
sudo systemctl status campus-memory

# 查看日志
sudo journalctl -u campus-memory -f
```

### 5. 使用 Docker 部署

**Dockerfile** (项目根目录已有):

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o campus-memory main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/campus-memory .
COPY --from=builder /app/scripts ./scripts
EXPOSE 8080
CMD ["./campus-memory"]
```

**docker-compose.yml**:

```yaml
version: '3.8'

services:
  campus-memory:
    build: .
    container_name: campus-memory
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
      - ./uploads:/app/uploads
      - ./logs:/app/logs
    environment:
      - GIN_MODE=release
      - SERVER_PORT=8080
    env_file:
      - .env.production
    restart: unless-stopped
```

部署命令：

```bash
# 构建并启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 重启服务
docker-compose restart
```

### 6. 使用 Nginx 反向代理

Nginx 配置文件 `/etc/nginx/sites-available/campus-memory`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 日志
    access_log /var/log/nginx/campus-memory-access.log;
    error_log /var/log/nginx/campus-memory-error.log;

    # 静态文件（图片上传）
    location /uploads {
        alias /opt/campus-memory/uploads;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # API 代理
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 前端（如果有）
    location / {
        root /var/www/campus-memory-frontend;
        try_files $uri $uri/ /index.html;
    }
}
```

启用配置：

```bash
# 创建软链接
sudo ln -s /etc/nginx/sites-available/campus-memory /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重载 Nginx
sudo systemctl reload nginx
```

### 7. HTTPS 配置（Let's Encrypt）

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

### 8. 监控和日志

```bash
# 查看应用日志
tail -f logs/app.log

# 查看系统服务日志
sudo journalctl -u campus-memory -f

# 查看 Nginx 日志
tail -f /var/log/nginx/campus-memory-access.log

# 监控资源使用
htop
docker stats  # 如果使用 Docker
```

### 9. 备份策略

```bash
# 数据库备份脚本
#!/bin/bash
BACKUP_DIR="/backup/campus-memory"
DATE=$(date +%Y%m%d_%H%M%S)

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份数据库
cp /opt/campus-memory/data/campus_memory.db $BACKUP_DIR/campus_memory_$DATE.db

# 备份上传文件
tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz /opt/campus-memory/uploads

# 删除 7 天前的备份
find $BACKUP_DIR -name "*.db" -mtime +7 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

添加到 crontab：

```bash
# 每天凌晨 2 点备份
0 2 * * * /opt/campus-memory/backup.sh
```

### 10. 性能优化

```bash
# 数据库优化
sqlite3 data/campus_memory.db "VACUUM;"
sqlite3 data/campus_memory.db "ANALYZE;"

# 启用 Gzip 压缩（Nginx）
gzip on;
gzip_types text/plain text/css application/json application/javascript;

# 启用缓存
# 在 Nginx 中配置静态资源缓存
```

## 📚 相关资源

- [Go 官方文档](https://golang.org/doc/)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [JWT 规范](https://jwt.io/)
- [RESTful API 设计指南](https://restfulapi.net/)
- [微信小程序开发文档](https://developers.weixin.qq.com/miniprogram/dev/framework/)
- [SQLite 文档](https://www.sqlite.org/docs.html)

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交代码 (`git commit -m 'feat: Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request
6. 等待代码审查

### 代码审查标准

- 代码符合项目规范
- 通过所有测试
- 添加必要的注释
- 更新相关文档
- 无明显的性能问题

## 📄 许可证

本项目采用 MIT 许可证，详见 [LICENSE](../LICENSE) 文件。

---

**最后更新**: 2024-01-01  
**维护者**: 开发团队  
**项目地址**: https://github.com/your-org/campus-memory