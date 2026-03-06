# 校园树洞 

基于微信小程序的校园记忆分享平台后端

## 功能特性

- 微信一键登录,无需注册
- 多校区、多地点记忆分享
- 图片上传与管理
- 点赞、评论互动系统
- 快速导航与地点搜索

## 技术栈

- Go 1.24 + Gin
- SQLite (GORM)
- JWT 认证
- Redis (可选，暂未实现)

## 项目结构

```
├── api/              # API 层 (handler, router)
├── application/      # 应用层 (service, dto, assembler)
├── infra/            # 基础设施层 (model, repo, cache)
├── middleware/       # 中间件 (auth, cors, logger)
├── utils/            # 工具函数 (jwt, wechat)
├── types/            # 类型定义
├── docs/             # API 文档
└── data/             # SQLite 数据库文件
```

## 快速开始

### 环境要求

- Go 1.24+
- 微信小程序 AppID 和 Secret

### 配置

创建 `.env` 文件:

```env
# 微信小程序
WECHAT_APPID=your_appid
WECHAT_SECRET=your_secret

# JWT
JWT_SECRET=your_jwt_secret
JWT_EXPIRE_HOURS=168
```

### 运行

```bash
# 安装依赖
go mod download

# 初始化数据库
sqlite3 data/campus_memory.db < scripts/init_db.sql
sqlite3 data/campus_memory.db < scripts/init_campus_data.sql

# 启动服务
go run main.go
# 或使用 make
make dev
```

服务运行在 `http://localhost:8080`

### Docker 部署

```bash
docker-compose up -d
```

## API 文档

详细 API 文档: [docs/API.md](docs/API.md)

### 核心接口

**认证**
- `POST /api/auth/wechat/login` - 微信登录
- `GET /api/auth/profile` - 获取个人信息
- `PUT /api/auth/profile` - 更新个人信息

**校区地点**
- `GET /api/campuses` - 获取校区列表
- `GET /api/campuses/:id/locations` - 获取校区地点
- `GET /api/quicknav/tree` - 获取导航树

**记忆**
- `POST /api/memories` - 创建记忆
- `GET /api/memories` - 获取记忆列表
- `GET /api/memories/:id` - 获取记忆详情
- `PUT /api/memories/:id` - 更新记忆
- `DELETE /api/memories/:id` - 删除记忆

**互动**
- `POST /api/comments` - 发表评论
- `POST /api/memories/:id/like` - 点赞记忆
- `POST /api/comments/:id/like` - 点赞评论

**图片**
- `POST /api/images/upload` - 上传图片
- `GET /api/images/memory/:id` - 获取记忆图片

## 数据模型

### 核心实体

- **User**: 用户 (OpenID, 昵称, 头像, 默认校区)
- **Campus**: 校区 (名称, 状态)
- **Location**: 地点 (名称, 类别, 所属校区, 记忆数)
- **Memory**: 记忆 (标题, 内容, 地点, 创建者, 点赞数, 评论数)
- **Comment**: 评论 (内容, 记忆, 创建者, 父评论, 点赞数)
- **Like**: 点赞 (用户, 目标类型, 目标ID)
- **Image**: 图片 (URL, 记忆, 大小)

### 关系

```
Campus (1) ─── (N) Location
Location (1) ─── (N) Memory
User (1) ─── (N) Memory
Memory (1) ─── (N) Comment
Memory (1) ─── (N) Image
User (1) ─── (N) Like
```

## 认证流程

1. 小程序调用 `wx.login()` 获取 code
2. 前端发送 code 到 `/api/auth/wechat/login`
3. 后端用 code 换取 OpenID
4. 生成 JWT token 返回
5. 后续请求携带 token: `Authorization: Bearer {token}`

## 开发

### 项目架构

采用 DDD 分层架构:
- **API 层**: HTTP 请求处理、参数验证
- **应用层**: 业务逻辑编排、DTO 转换
- **基础设施层**: 数据持久化、外部服务

### 代码规范

- 遵循 Go 官方代码规范
- 使用有意义的命名
- 添加必要的注释

### 常用命令

```bash
# 开发
make dev

# 构建
make build

# 测试
make test

# 清理
make clean
```

## License

MIT License
