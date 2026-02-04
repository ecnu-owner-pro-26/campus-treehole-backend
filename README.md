# 校园树洞后端项目

一个简单的校园记忆分享应用后端，支持用户在不同校区地点发表记忆、评论和点赞。

## 功能特性

- 用户注册登录
- 校区地点选择
- 记忆发表和浏览
- 评论和点赞功能

## 技术栈

- **语言**: Go 1.21
- **框架**: Gin
- **数据库**: SQLite + GORM
- **认证**: JWT

## 项目结构

```
├── api/                 # API层
│   ├── handler/         # HTTP处理器
│   └── router/          # 路由配置
├── application/         # 应用层
│   ├── dto/             # 数据传输对象
│   └── service/         # 业务逻辑
├── infra/               # 基础设施层
│   ├── model/           # 数据库模型
│   └── repo/            # 数据访问层
├── middleware/          # 中间件
├── utils/               # 工具函数
├── scripts/             # 脚本文件
└── docs/                # 文档
```

## 快速开始

### 1. 克隆项目
```bash
git clone <repository-url>
cd campus-treehole-backend
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 运行项目
```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动

### 4. 初始化数据（可选）
```bash
# 执行SQL脚本添加测试数据
sqlite3 ./data/campus_memory.db < scripts/init_campus_data.sql
```

## API接口

### 校区相关
- `GET /api/campuses` - 获取校区列表
- `GET /api/campuses/:id/locations` - 获取校区地点

### 用户认证
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/profile` - 获取个人信息

### 记忆管理
- `POST /api/memories` - 发表记忆
- `GET /api/memories/:id` - 获取记忆详情
- `GET /api/locations/:id/memories` - 获取地点记忆列表

### 互动功能
- `POST /api/memories/:id/like` - 点赞记忆
- `POST /api/memories/:id/comments` - 评论记忆
- `POST /api/comments/:id/like` - 点赞评论

## 开发说明

这是一个学习项目，代码结构简单清晰，适合Go初学者学习Web开发。

### 开发流程
1. 数据库基础设施
2. 用户认证系统
3. 校区地点管理
4. 记忆系统
5. 互动功能

## 许可证

MIT License