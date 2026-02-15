# API 文档

## 基础信息
- Base URL: `http://localhost:8080/api`
- 认证方式: Bearer Token (JWT)
- 登录方式: 仅支持微信小程序登录

---

## 认证相关接口

### 微信登录
- **接口**: `POST /auth/wechat/login`
- **描述**: 微信小程序登录，获取JWT token
- **请求头**: 无需认证
- **请求体**:
```json
{
  "code": "微信wx.login()返回的code",
  "nickname": "用户昵称（可选）",
  "avatar": "用户头像URL（可选）"
}
```
- **响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "jwt_token_string",
    "user": {
      "id": 1,
      "nickname": "用户昵称",
      "avatar": "头像URL",
      "default_campus_id": 1
    }
  }
}
```

### 获取个人信息
- **接口**: `GET /auth/profile`
- **描述**: 获取当前登录用户的详细信息
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "nickname": "用户昵称",
    "avatar": "头像URL",
    "default_campus_id": 1,
    "status": 1,
    "role": 0,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 更新个人信息
- **接口**: `PUT /auth/profile`
- **描述**: 更新当前登录用户的信息
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:
```json
{
  "nickname": "新昵称",
  "avatar": "新头像URL",
  "default_campus_id": 1
}
```

---

## 校区地点相关接口

### 获取校区列表
- **接口**: `GET /campuses`
- **描述**: 获取所有校区列表
- **请求头**: 无需认证
- **响应**:
```json
{
  "code": 200,
  "data": {
    "campuses": [
      {
        "id": 1,
        "name": "主校区",
        "is_active": true,
        "sort_order": 1
      }
    ]
  }
}
```

### 获取校区地点列表
- **接口**: `GET /campuses/:id/locations`
- **描述**: 获取指定校区的所有地点
- **请求头**: 无需认证
- **路径参数**:
  - `id`: 校区ID
- **响应**:
```json
{
  "code": 200,
  "data": {
    "campus": {
      "id": 1,
      "name": "主校区"
    },
    "locations": [
      {
        "id": 1,
        "name": "图书馆",
        "category": "library",
        "address": "主校区图书馆",
        "memory_count": 10
      }
    ]
  }
}
```

---

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
  "location_id": 1,
  "is_public": true,
  "tags": ["标签1", "标签2"],
  "images": ["图片URL1", "图片URL2"]
}
```

### 获取记忆详情
- **接口**: `GET /memories/:id`
- **描述**: 获取指定记忆的详细信息
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id`: 记忆ID
- **响应**:
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "title": "记忆标题",
    "content": "记忆内容",
    "location_name": "图书馆",
    "location_id": 1,
    "creator": {
      "id": 1,
      "nickname": "用户昵称",
      "avatar": "头像URL"
    },
    "is_public": true,
    "like_count": 10,
    "comment_count": 5,
    "view_count": 100,
    "is_liked": false,
    "tags": ["标签1", "标签2"],
    "images": ["图片URL1", "图片URL2"],
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 获取记忆列表
- **接口**: `GET /memories`
- **描述**: 获取记忆列表
- **请求头**: `Authorization: Bearer {token}`
- **查询参数**:
  - `page`: 页码（默认1）
  - `page_size`: 每页数量（默认10，最大100）
  - `location_id`: 地点ID（可选）
  - `is_public`: 是否公开（可选）
  - `creator_id`: 创建者ID（可选，查看我的记忆）
- **响应**:
```json
{
  "code": 200,
  "data": {
    "memories": [...],
    "total": 100,
    "page": 1,
    "page_size": 10,
    "has_more": true
  }
}
```

### 更新记忆
- **接口**: `PUT /memories/:id`
- **描述**: 更新记忆信息
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id`: 记忆ID
- **请求体**:
```json
{
  "title": "新标题",
  "content": "新内容",
  "is_public": false,
  "tags": ["新标签"]
}
```

### 删除记忆
- **接口**: `DELETE /memories/:id`
- **描述**: 删除记忆
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id`: 记忆ID

### 切换记忆点赞
- **接口**: `POST /memories/:id/like`
- **描述**: 切换记忆点赞状态（已点赞则取消，未点赞则点赞）
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id`: 记忆ID
- **响应**:
```json
{
  "code": 200,
  "data": {
    "is_liked": true,
    "like_count": 11
  }
}
```

---

## 评论相关接口

### 创建评论
- **接口**: `POST /memories/:id/comments`
- **描述**: 为记忆创建评论
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id`: 记忆ID
- **请求体**:
```json
{
  "content": "评论内容",
  "parent_id": null,
  "reply_to_user_id": null
}
```

### 获取评论列表
- **接口**: `GET /memories/:id/comments`
- **描述**: 获取记忆的评论列表
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id`: 记忆ID
- **查询参数**:
  - `page`: 页码（默认1）
  - `page_size`: 每页数量（默认20，最大100）
  - `parent_id`: 父评论ID（可选，获取子评论）
- **响应**:
```json
{
  "code": 200,
  "data": {
    "comments": [
      {
        "id": 1,
        "content": "评论内容",
        "creator": {
          "id": 1,
          "nickname": "用户昵称",
          "avatar": "头像URL"
        },
        "like_count": 5,
        "is_liked": false,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

### 删除评论
- **接口**: `DELETE /comments/:id`
- **描述**: 删除评论
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id`: 评论ID

### 切换评论点赞
- **接口**: `POST /comments/:id/like`
- **描述**: 切换评论点赞状态（已点赞则取消，未点赞则点赞）
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id`: 评论ID
- **响应**:
```json
{
  "code": 200,
  "data": {
    "is_liked": true,
    "like_count": 6
  }
}
```

---

## 地点相关接口

### 获取地点记忆列表
- **接口**: `GET /locations/:id/memories`
- **描述**: 获取指定地点的所有记忆
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id`: 地点ID
- **查询参数**:
  - `page`: 页码（默认1）
  - `page_size`: 每页数量（默认10，最大100）

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未登录或token无效 |
| 403 | 无权限访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 通用响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```
