# 校园树洞项目

一个基于地图的校园记忆分享平台，用户可以在地图上标记位置并存储照片、文字、音频等记忆内容。

## 功能特性

- 📍 地图标记：点击地图任意位置创建记忆点
- 📷 多媒体支持：支持上传照片、文字、音频
- 🔒 隐私控制：可设置记忆为公开或私密
- 💬 留言互动：公开记忆支持他人留言
- 🗺️ 位置浏览：查看附近或特定位置的记忆

## 项目结构

```
campus-memory/
├── main.go                # 应用入口
├── api/                   # API层
│   ├── handler/          # HTTP请求处理器
│   ├── router/           # 路由定义
│   └── token/            # JWT令牌管理
├── application/           # 应用层
│   ├── assembler/        # 数据组装器
│   ├── dto/              # 数据传输对象
│   └── service/          # 业务服务
├── infra/                 # 基础设施层
│   ├── cache/            # 缓存（Redis）
│   ├── config/           # 配置管理
│   ├── model/            # 数据库模型
│   ├── repo/             # 数据仓储
│   ├── util/             # 工具函数
│   ├── database.go       # 数据库连接
│   └── storage.go        # 文件存储
├── provider/              # 依赖注入
├── types/                 # 类型定义
│   ├── consts/           # 常量定义
│   ├── errno/            # 错误定义
│   └── mapping/          # 数据映射
├── internal/              # 内部代码
│   └── middleware/       # 中间件
├── docs/                  # 文档和Swagger
├── etc/                   # 配置文件
└── storage/              # 文件存储目录
```

## 架构说明

项目采用DDD分层架构设计：

- **api**: API层，处理HTTP请求和路由
- **application**: 应用层，业务逻辑和服务编排
- **infra**: 基础设施层，数据库、缓存、存储等
- **provider**: 依赖注入，使用wire工具管理依赖
- **types**: 类型定义，包括常量、错误、映射等
- **internal/middleware**: 中间件，认证、日志等

## 技术栈

- TODO: 添加使用的技术栈说明

## 快速开始

- TODO: 添加项目启动说明

## API 文档

详见 [API.md](docs/API.md)

## 数据库设计

详见 [DATABASE.md](docs/DATABASE.md)
