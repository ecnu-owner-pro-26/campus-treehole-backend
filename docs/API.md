# 校园记忆 API 文档

## 📖 目录

1. [基础信息](#基础信息)
2. [快速开始 - 用户使用流程](#快速开始---用户使用流程)
3. [接口详细说明](#接口详细说明)
   - [认证相关接口](#认证相关接口)
   - [校区地点相关接口](#校区地点相关接口)
   - [记忆相关接口](#记忆相关接口)
   - [评论相关接口](#评论相关接口)
   - [地点相关接口](#地点相关接口)
4. [数据模型](#数据模型)
5. [错误码说明](#错误码说明)

---

## 基础信息

- **Base URL**: `http://localhost:8080`
- **API 前缀**: `/api`
- **认证方式**: Bearer Token (JWT)
- **登录方式**: 仅支持微信小程序登录
- **响应格式**: JSON

### 通用响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

---

## 快速开始 - 用户使用流程

### 典型用户流程

```
1. 用户打开应用
   ↓
   GET /api/campuses (获取校区列表)
   
2. 选择"普陀校区"
   ↓
   GET /api/campuses/1/locations (获取该校区的地点列表)
   
3. 选择"图书馆"
   ↓
   GET /api/locations/1/memories (浏览该地点的记忆)
   或
   POST /api/memories (发表新记忆)
   
4. 与记忆互动
   ↓
   POST /api/memories/1/like (点赞记忆)
   POST /api/memories/1/comments (评论记忆)
   POST /api/comments/1/like (点赞评论)
```

### 认证流程

```
1. 微信小程序登录
   ↓
   wx.login() 获取 code
   ↓
   POST /api/auth/wechat/login (发送 code 到后端)
   ↓
   获取 JWT token
   ↓
   后续请求在 Header 中携带: Authorization: Bearer {token}
```

---

## 接口详细说明

### 认证相关接口

#### 1. 微信登录

- **接口**: `POST /api/auth/wechat/login`
- **描述**: 微信小程序登录，获取 JWT token
- **认证**: 无需认证
- **请求体**:

```json
{
  "code": "微信wx.login()返回的code",
  "nickname": "用户昵称（可选）",
  "avatar": "用户头像URL（可选）"
}
```

**字段说明**:
- `code` (string, 必填): 微信 wx.login() 返回的临时登录凭证
- `nickname` (string, 可选): 用户昵称
- `avatar` (string, 可选): 用户头像 URL

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "nickname": "用户昵称",
      "avatar": "https://example.com/avatar.jpg",
      "default_campus_id": 1
    }
  }
}
```

---

#### 2. 获取个人信息

- **接口**: `GET /api/auth/profile`
- **描述**: 获取当前登录用户的详细信息
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **响应**:

```json
{
  "code": 200,
  "data": {
    "id": 1,
    "nickname": "用户昵称",
    "avatar": "https://example.com/avatar.jpg",
    "default_campus_id": 1,
    "status": 1,
    "role": 0,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**字段说明**:
- `status`: 0-禁用, 1-正常
- `role`: 0-普通用户, 1-管理员

---

#### 3. 更新个人信息

- **接口**: `PUT /api/auth/profile`
- **描述**: 更新当前登录用户的信息
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "nickname": "新昵称",
  "avatar": "https://example.com/new-avatar.jpg",
  "default_campus_id": 2
}
```

**字段说明**:
- 所有字段都是可选的，只更新传入的字段
- `nickname`: 1-50 个字符
- `avatar`: 必须是有效的 URL
- `default_campus_id`: 必须大于 0

- **响应**:

```json
{
  "code": 200,
  "message": "success"
}
```

---

### 校区地点相关接口

#### 4. 获取校区列表

- **接口**: `GET /api/campuses`
- **描述**: 获取所有可用校区列表
- **认证**: 无需认证
- **使用场景**: 用户打开应用时调用
- **响应**:

```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "name": "普陀校区",
      "is_active": true,
      "sort_order": 1
    },
    {
      "id": 2,
      "name": "临港校区",
      "is_active": true,
      "sort_order": 2
    }
  ]
}
```

---

#### 5. 获取校区地点列表

- **接口**: `GET /api/campuses/:id/locations`
- **描述**: 获取指定校区的所有地点
- **认证**: 无需认证
- **使用场景**: 用户选择校区后，获取该校区下的所有地点
- **路径参数**:
  - `id` (int64): 校区 ID

- **响应**:

```json
{
  "code": 200,
  "data": {
    "campus": {
      "id": 1,
      "name": "普陀校区"
    },
    "locations": [
      {
        "id": 1,
        "name": "图书馆",
        "category": "teaching",
        "is_active": true,
        "sort_order": 1,
        "memory_count": 25
      },
      {
        "id": 2,
        "name": "第一食堂",
        "category": "dining",
        "is_active": true,
        "sort_order": 2,
        "memory_count": 18
      }
    ]
  }
}
```

**字段说明**:
- `category`: 地点类别（teaching-教学, dining-餐厅, dormitory-宿舍, scenic-景点）
- `memory_count`: 该地点的记忆数量

---

### 记忆相关接口

#### 6. 创建记忆

- **接口**: `POST /api/memories`
- **描述**: 创建新的记忆
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "title": "图书馆的美好时光",
  "content": "今天在图书馆学习，感觉很充实",
  "location_id": 1,
  "is_public": true,
  "tags": ["学习", "图书馆"],
  "image_urls": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ]
}
```

**字段说明**:
- `title` (string, 必填): 记忆标题，最大 100 字符
- `content` (string, 可选): 记忆内容，最大 5000 字符
- `location_id` (int64, 必填): 关联的地点 ID
- `is_public` (bool, 可选): 是否公开，默认 true
- `tags` (array, 可选): 标签数组
- `image_urls` (array, 可选): 图片 URL 数组

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "title": "图书馆的美好时光",
    "content": "今天在图书馆学习，感觉很充实",
    "location_name": "图书馆",
    "location_id": 1,
    "creator": {
      "id": 1,
      "nickname": "小明",
      "avatar": "https://example.com/avatar.jpg"
    },
    "like_count": 0,
    "comment_count": 0,
    "view_count": 0,
    "is_liked": false,
    "tags": ["学习", "图书馆"],
    "images": [
      {
        "id": 1,
        "url": "https://example.com/image1.jpg"
      }
    ],
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

---

#### 7. 获取记忆详情

- **接口**: `GET /api/memories/:id`
- **描述**: 获取指定记忆的详细信息
- **认证**: 需要 JWT 认证（用于判断是否已点赞）
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id` (int64): 记忆 ID

- **响应**:

```json
{
  "code": 200,
  "data": {
    "id": 1,
    "title": "图书馆的美好时光",
    "content": "今天在图书馆学习，感觉很充实",
    "location_name": "图书馆",
    "location_id": 1,
    "creator": {
      "id": 1,
      "nickname": "小明",
      "avatar": "https://example.com/avatar.jpg"
    },
    "like_count": 10,
    "comment_count": 5,
    "view_count": 100,
    "is_liked": false,
    "tags": ["学习", "图书馆"],
    "images": [
      {
        "id": 1,
        "url": "https://example.com/image1.jpg"
      },
      {
        "id": 2,
        "url": "https://example.com/image2.jpg"
      }
    ],
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

---

#### 8. 获取记忆列表

- **接口**: `GET /api/memories`
- **描述**: 获取记忆列表（支持分页和筛选）
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **查询参数**:
  - `page` (int, 可选): 页码，默认 1
  - `page_size` (int, 可选): 每页数量，默认 10，最大 100
  - `location_id` (int64, 可选): 筛选指定地点的记忆
  - `creator_id` (int64, 可选): 筛选指定用户的记忆（查看我的记忆）
  - `sort_by` (string, 可选): 排序方式（latest-最新, popular-最热门）

- **响应**:

```json
{
  "code": 200,
  "data": {
    "memories": [
      {
        "id": 1,
        "title": "图书馆的美好时光",
        "content": "今天在图书馆学习...",
        "location_name": "图书馆",
        "location_id": 1,
        "creator": {
          "id": 1,
          "nickname": "小明",
          "avatar": "https://example.com/avatar.jpg"
        },
        "like_count": 10,
        "comment_count": 5,
        "view_count": 100,
        "is_liked": false,
        "tags": ["学习", "图书馆"],
        "images": [],
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

#### 9. 更新记忆

- **接口**: `PUT /api/memories/:id`
- **描述**: 更新记忆信息（只能更新自己的记忆）
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id` (int64): 记忆 ID

- **请求体**:

```json
{
  "title": "新标题",
  "content": "新内容",
  "is_public": false,
  "tags": ["新标签"],
  "image_urls": ["https://example.com/new-image.jpg"]
}
```

**字段说明**:
- 所有字段都是可选的，只更新传入的字段
- `title`: 最大 100 字符
- `content`: 最大 5000 字符

- **响应**:

```json
{
  "code": 200,
  "message": "success"
}
```

---

#### 10. 删除记忆

- **接口**: `DELETE /api/memories/:id`
- **描述**: 删除记忆（只能删除自己的记忆）
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id` (int64): 记忆 ID

- **响应**:

```json
{
  "code": 200,
  "message": "success"
}
```

---

#### 11. 切换记忆点赞

- **接口**: `POST /api/memories/:id/like`
- **描述**: 切换记忆点赞状态（已点赞则取消，未点赞则点赞）
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id` (int64): 记忆 ID

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

**字段说明**:
- `is_liked`: 当前点赞状态（true-已点赞, false-未点赞）
- `like_count`: 点赞总数

---

### 评论相关接口

#### 12. 创建评论

- **接口**: `POST /api/memories/:id/comments`
- **描述**: 为记忆创建评论
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id` (int64): 记忆 ID

- **请求体**:

```json
{
  "content": "很棒的分享！",
  "parent_id": null,
  "reply_to_user_id": null
}
```

**字段说明**:
- `content` (string, 必填): 评论内容
- `parent_id` (int64, 可选): 父评论 ID（用于回复评论）
- `reply_to_user_id` (int64, 可选): 被回复的用户 ID

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "content": "很棒的分享！",
    "creator": {
      "id": 1,
      "nickname": "小明",
      "avatar": "https://example.com/avatar.jpg"
    },
    "like_count": 0,
    "is_liked": false,
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

---

#### 13. 获取评论列表

- **接口**: `GET /api/memories/:id/comments`
- **描述**: 获取记忆的评论列表
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id` (int64): 记忆 ID

- **查询参数**:
  - `page` (int, 可选): 页码，默认 1
  - `page_size` (int, 可选): 每页数量，默认 20，最大 100
  - `parent_id` (int64, 可选): 父评论 ID（获取子评论）

- **响应**:

```json
{
  "code": 200,
  "data": {
    "comments": [
      {
        "id": 1,
        "content": "很棒的分享！",
        "creator": {
          "id": 1,
          "nickname": "小明",
          "avatar": "https://example.com/avatar.jpg"
        },
        "like_count": 5,
        "is_liked": false,
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

---

#### 14. 删除评论

- **接口**: `DELETE /api/comments/:id`
- **描述**: 删除评论（只能删除自己的评论）
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id` (int64): 评论 ID

- **响应**:

```json
{
  "code": 200,
  "message": "success"
}
```

---

#### 15. 切换评论点赞

- **接口**: `POST /api/comments/:id/like`
- **描述**: 切换评论点赞状态（已点赞则取消，未点赞则点赞）
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id` (int64): 评论 ID

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

### 地点相关接口

#### 16. 获取地点记忆列表

- **接口**: `GET /api/locations/:id/memories`
- **描述**: 获取指定地点的所有记忆
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **使用场景**: 用户选择地点后，浏览该地点的所有记忆
- **路径参数**:
  - `id` (int64): 地点 ID

- **查询参数**:
  - `page` (int, 可选): 页码，默认 1
  - `page_size` (int, 可选): 每页数量，默认 10，最大 100

- **响应**:

```json
{
  "code": 200,
  "data": {
    "memories": [
      {
        "id": 1,
        "title": "图书馆的美好时光",
        "content": "今天在图书馆学习...",
        "location_name": "图书馆",
        "location_id": 1,
        "creator": {
          "id": 1,
          "nickname": "小明",
          "avatar": "https://example.com/avatar.jpg"
        },
        "like_count": 10,
        "comment_count": 5,
        "view_count": 100,
        "is_liked": false,
        "tags": ["学习", "图书馆"],
        "images": [],
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "page_size": 10
  }
}
```

---

## 数据模型

### 数据库关联关系

```
用户 (users)
  ↓ 创建
记忆 (memories)
  ↓ 关联
地点 (locations)
  ↓ 属于
校区 (campuses)

记忆 (memories)
  ↓ 包含
评论 (comments)
  ↓ 创建者
用户 (users)

点赞 (likes)
  ├─ 关联 → 记忆 (memories)
  ├─ 关联 → 评论 (comments)
  └─ 创建者 → 用户 (users)
```

### 核心实体

#### 用户 (User)
```
- id: 用户ID
- openid: 微信OpenID（唯一）
- nickname: 用户昵称
- avatar: 用户头像
- default_campus_id: 默认校区ID
- status: 状态（0-禁用, 1-正常）
- role: 角色（0-普通用户, 1-管理员）
```

#### 校区 (Campus)
```
- id: 校区ID
- name: 校区名称
- is_active: 是否启用
- sort_order: 显示顺序
```

#### 地点 (Location)
```
- id: 地点ID
- campus_id: 所属校区ID
- name: 地点名称
- category: 地点类别
- is_active: 是否启用
- sort_order: 显示顺序
- memory_count: 记忆数量
```

#### 记忆 (Memory)
```
- id: 记忆ID
- title: 标题
- content: 内容
- location_id: 关联地点ID
- location_name: 地点名称
- creator_id: 创建者ID
- is_public: 是否公开
- like_count: 点赞数
- comment_count: 评论数
- view_count: 浏览数
- tags: 标签（JSON数组）
- status: 状态（0-待审核, 1-已发布, 2-已下架）
```

#### 评论 (Comment)
```
- id: 评论ID
- memory_id: 记忆ID
- content: 评论内容
- creator_id: 创建者ID
- parent_id: 父评论ID（用于回复）
- reply_to_user_id: 被回复的用户ID
- like_count: 点赞数
```

#### 点赞 (Like)
```
- id: 点赞ID
- user_id: 用户ID
- target_type: 目标类型（1-记忆, 2-评论）
- target_id: 目标ID
```

---

## 错误码说明

### HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未登录或 token 无效 |
| 403 | 无权限访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

### 业务错误码

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未登录或 token 无效 |
| 403 | 无权限访问（如删除他人的记忆） |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

### 错误响应格式

```json
{
  "code": 400,
  "message": "请求参数错误: title 字段不能为空"
}
```

---

## 附录

### 字段验证规则

#### 创建记忆
- `title`: 必填，最大 100 字符
- `content`: 可选，最大 5000 字符
- `location_id`: 必填，必须是有效的地点 ID
- `is_public`: 可选，布尔值
- `tags`: 可选，字符串数组
- `image_urls`: 可选，字符串数组

#### 更新用户信息
- `nickname`: 可选，1-50 字符
- `avatar`: 可选，必须是有效的 URL
- `default_campus_id`: 可选，必须大于 0

### 分页说明

所有列表接口都支持分页：
- `page`: 页码，从 1 开始
- `page_size`: 每页数量，默认 10，最大 100

响应中包含：
- `total`: 总记录数
- `page`: 当前页码
- `page_size`: 每页数量

### 认证说明

需要认证的接口必须在请求头中携带 JWT token：

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

Token 有效期：7 天（可通过环境变量 `JWT_EXPIRE_HOURS` 配置）

---

**文档版本**: v1.0  
**最后更新**: 2024-01-01  
**维护者**: 开发团队
