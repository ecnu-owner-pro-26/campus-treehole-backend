# 校园树洞后端 - 微信小程序版

基于微信小程序的校园记忆分享平台后端服务

## 项目特点

- **纯微信登录**：无需注册，使用微信一键登录
- **校区地点导航**：支持多校区、多地点的记忆分享
- **记忆分享**：支持文字、图片、音频的记忆发布
- **互动功能**：点赞、评论系统

## 技术栈

- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: PostgreSQL / MySQL
- **缓存**: Redis
- **认证**: JWT
- **文档**: Swagger

## 项目结构

```
.
├── api/                    # API层
│   ├── handler/           # 请求处理器
│   ├── router/            # 路由配置
│   └── token/             # Token管理
├── application/           # 应用层
│   ├── assembler/        # 数据组装器
│   ├── dto/              # 数据传输对象
│   └── service/          # 业务逻辑服务
├── infra/                # 基础设施层
│   ├── cache/            # 缓存
│   ├── model/            # 数据模型
│   ├── repo/             # 数据仓储
│   └── util/             # 工具函数
├── middleware/           # 中间件
├── types/                # 类型定义
├── utils/                # 通用工具
└── docs/                 # API文档

```

## 认证流程

### 微信登录流程

1. 小程序调用 `wx.login()` 获取 code
2. 前端携带 code 调用 `/api/auth/wechat/login`
3. 后端使用 code 换取 OpenID
4. 查询用户是否存在，不存在则创建
5. 生成 JWT token 返回给前端
6. 前端后续请求携带 token 访问受保护接口

### 用户数据模型

```go
type UserModel struct {
    ID              int64     // 用户ID
    OpenID          string    // 微信OpenID（唯一标识）
    UnionID         string    // 微信UnionID（可选）
    Nickname        string    // 用户昵称
    Avatar          string    // 用户头像
    DefaultCampusID *int64    // 默认校区
    Status          int8      // 状态（0-禁用 1-正常）
    Role            int8      // 角色（0-普通用户 1-管理员）
}
```

## 环境配置

创建 `.env` 文件：

```env
# 微信小程序配置
WECHAT_APPID=your_appid
WECHAT_SECRET=your_secret

# JWT配置
JWT_SECRET=your_jwt_secret
JWT_EXPIRE_HOURS=168

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=campus_treehole

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

## 快速开始

### 开发环境

```bash
# 安装依赖
go mod download

# 初始化数据库
make init-db

# 运行开发服务器
make dev
```

### Docker部署

```bash
# 构建并启动
docker-compose up -d

# 查看日志
docker-compose logs -f
```

## API文档

启动服务后访问：
- Swagger UI: http://localhost:8080/swagger/index.html
- API文档: [docs/API.md](docs/API.md)

## 核心功能模块

### 1. 认证模块
- 微信登录
- 用户信息管理
- JWT认证中间件

### 2. 记忆模块
- 创建/编辑/删除记忆
- 记忆列表查询
- 记忆点赞

### 3. 评论模块
- 发表评论
- 评论点赞
- 评论管理

### 4. 校区地点模块
- 校区列表
- 地点列表
- 地点记忆查询

## 开发规范

### 分层架构

- **API层**: 处理HTTP请求，参数验证
- **应用层**: 业务逻辑编排，DTO转换
- **基础设施层**: 数据持久化，外部服务调用

### 代码风格

- 遵循Go官方代码规范
- 使用有意义的变量和函数命名
- 添加必要的注释和文档

## License

MIT License
