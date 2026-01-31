# 架构设计文档

## 分层架构

### API层
- 负责路由定义和注册
- 将HTTP请求分发到对应的handler

### Handler层 (internal/handler)
- 处理HTTP请求和响应
- 参数验证和数据转换
- 调用application层处理业务逻辑

### Application层 (application)
- 核心业务逻辑
- 协调provider层和infra层
- 事务管理

### Provider层 (provider)
- 数据访问层
- 封装数据库CRUD操作
- 数据查询和持久化

### Types层 (types)
- 定义数据结构
- 请求/响应模型
- 领域模型

### Infra层 (infra)
- 基础设施
- 数据库连接
- 文件存储
- 配置管理

### Middleware层 (internal/middleware)
- 认证授权
- 日志记录
- 跨域处理
- 错误处理

## 数据流向

```
HTTP请求 -> API路由 -> Handler -> Application -> Provider -> Database
                                      ↓
                                    Infra (Storage/Config)
```
