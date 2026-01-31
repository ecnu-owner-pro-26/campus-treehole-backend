# API 文档

## 基础信息
- Base URL: `http://localhost:8080/api`
- 认证方式: Bearer Token (JWT)

## 记忆相关接口

### 创建记忆
- **接口**: `POST /memories`
- **描述**: 创建新的记忆
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
```json
{
  "title": "记忆标题",
  "content": "记忆内容",
  "latitude": 39.9042,
  "longitude": 116.4074,
  "location_name": "天安门",
  "is_public": true,
  "allow_comments": true,
  "images": ["url1", "url2"],
  "audios": ["url1"]
}
```

### 获取记忆详情
- **接口**: `GET /memories/:id`
- **描述**: 获取指定记忆的详细信息

### 获取记忆列表
- **接口**: `GET /memories`
- **描述**: 获取记忆列表
- **查询参数**:
  - `page`: 页码（默认1）
  - `page_size`: 每页数量（默认10）
  - `is_public`: 是否公开（可选）

### 更新记忆
- **接口**: `PUT /memories/:id`
- **描述**: 更新记忆信息
- **请求头**: `Authorization: Bearer {token}`

### 删除记忆
- **接口**: `DELETE /memories/:id`
- **描述**: 删除记忆
- **请求头**: `Authorization: Bearer {token}`

## 留言相关接口

### 创建留言
- **接口**: `POST /memories/:id/comments`
- **描述**: 为记忆创建留言
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
```json
{
  "content": "留言内容"
}
```

### 获取留言列表
- **接口**: `GET /memories/:id/comments`
- **描述**: 获取记忆的留言列表

### 删除留言
- **接口**: `DELETE /comments/:id`
- **描述**: 删除留言
- **请求头**: `Authorization: Bearer {token}`

## 地图位置相关接口

### 根据位置获取记忆
- **接口**: `GET /locations/:lat/:lng/memories`
- **描述**: 获取指定位置附近的记忆
- **路径参数**:
  - `lat`: 纬度
  - `lng`: 经度
- **查询参数**:
  - `radius`: 半径（米，默认1000）
