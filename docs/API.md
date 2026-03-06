# 校园记忆 API 文档

## 📖 目录

1. [基础信息](#基础信息)
2. [快速开始 - 用户使用流程](#快速开始---用户使用流程)
3. [接口详细说明](#接口详细说明)
   - [认证相关接口](#认证相关接口)
   - [校区相关接口](#校区相关接口)
   - [地点相关接口](#地点相关接口)
   - [快速导航相关接口](#快速导航相关接口)
   - [记忆相关接口](#记忆相关接口)
   - [评论相关接口](#评论相关接口)
   - [点赞相关接口](#点赞相关接口)
   - [图片相关接口](#图片相关接口)
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
   或
   GET /api/quicknav/tree?campusId=1 (获取导航树)
   
3. 选择"图书馆"
   ↓
   GET /api/memories?locationId=1 (浏览该地点的记忆)
   或
   POST /api/memories (发表新记忆)
   
4. 与记忆互动
   ↓
   POST /api/memories/1/like (点赞记忆)
   POST /api/comments (评论记忆)
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
      "defaultCampusId": 1
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
  "message": "success",
  "data": {
    "id": 1,
    "nickname": "用户昵称",
    "avatar": "https://example.com/avatar.jpg",
    "defaultCampusId": 1,
    "status": 1,
    "role": 0,
    "createdAt": "2024-01-01T00:00:00Z"
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
  "defaultCampusId": 2
}
```

**字段说明**:
- 所有字段都是可选的，只更新传入的字段
- `nickname`: 1-50 个字符
- `avatar`: 必须是有效的 URL
- `defaultCampusId`: 必须大于 0

- **响应**:

```json
{
  "code": 200,
  "message": "success"
}
```

---

### 校区相关接口

#### 4. 获取校区列表

- **接口**: `GET /api/campuses`
- **描述**: 获取所有可用校区列表
- **认证**: 无需认证
- **使用场景**: 用户打开应用时调用
- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "campuses": [
      {
        "id": 1,
        "name": "普陀校区",
        "created_at": "2024-01-01T00:00:00Z"
      },
      {
        "id": 2,
        "name": "临港校区",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 2
  }
}
```

---

#### 5. 获取校区详情

- **接口**: `GET /api/campuses/:id`
- **描述**: 获取指定校区的详细信息
- **认证**: 无需认证
- **路径参数**:
  - `id` (int64): 校区 ID

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "普陀校区",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

#### 6. 获取校区地点列表

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
  "message": "success",
  "data": {
    "campus": {
      "id": 1,
      "name": "普陀校区",
      "created_at": "2024-01-01T00:00:00Z"
    },
    "locations": [
      {
        "id": 1,
        "name": "图书馆",
        "campus_id": 1,
        "category": "teaching",
        "latitude": 31.2304,
        "longitude": 121.4737,
        "memory_count": 25,
        "created_at": "2024-01-01T00:00:00Z"
      },
      {
        "id": 2,
        "name": "第一食堂",
        "campus_id": 1,
        "category": "dining",
        "latitude": 31.2305,
        "longitude": 121.4738,
        "memory_count": 18,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

**字段说明**:
- `category`: 地点类别（teaching-教学, dining-餐厅, dormitory-宿舍, scenic-景点）
- `latitude`: 纬度
- `longitude`: 经度
- `memory_count`: 该地点的记忆数量

---

### 地点相关接口

#### 7. 获取地点列表

- **接口**: `GET /api/locations`
- **描述**: 获取地点列表（支持分页和筛选）
- **认证**: 无需认证
- **查询参数**:
  - `campus_id` (int64, 可选): 筛选指定校区的地点
  - `category` (string, 可选): 筛选指定类别的地点
  - `page` (int, 必填): 页码，最小值 1
  - `page_size` (int, 必填): 每页数量，最小值 1，最大值 100

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "locations": [
      {
        "id": 1,
        "name": "图书馆",
        "campus_id": 1,
        "category": "teaching",
        "latitude": 31.2304,
        "longitude": 121.4737,
        "memory_count": 25,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

---

#### 8. 获取地点详情

- **接口**: `GET /api/locations/:id`
- **描述**: 获取指定地点的详细信息
- **认证**: 无需认证
- **路径参数**:
  - `id` (int64): 地点 ID

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "图书馆",
    "campus_id": 1,
    "category": "teaching",
    "latitude": 31.2304,
    "longitude": 121.4737,
    "memory_count": 25,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

#### 9. 搜索地点

- **接口**: `GET /api/locations/search`
- **描述**: 根据关键词搜索地点
- **认证**: 无需认证
- **查询参数**:
  - `keyword` (string, 必填): 搜索关键词
  - `campus_id` (int64, 可选): 限定校区

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "图书馆",
      "campus_id": 1,
      "category": "teaching",
      "latitude": 31.2304,
      "longitude": 121.4737,
      "memory_count": 25,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

#### 10. 创建地点

- **接口**: `POST /api/locations`
- **描述**: 创建新地点（管理功能）
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "campus_id": 1,
  "name": "新图书馆",
  "category": "teaching",
  "latitude": 31.2304,
  "longitude": 121.4737,
  "sort_order": 10
}
```

**字段说明**:
- `campus_id` (int64, 必填): 所属校区 ID
- `name` (string, 必填): 地点名称，最大 100 字符
- `category` (string, 可选): 地点类别，最大 50 字符
- `latitude` (float64, 必填): 纬度，范围 -90 到 90
- `longitude` (float64, 必填): 经度，范围 -180 到 180
- `sort_order` (int, 可选): 显示顺序

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 10,
    "name": "新图书馆",
    "campus_id": 1,
    "category": "teaching",
    "latitude": 31.2304,
    "longitude": 121.4737,
    "memory_count": 0,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

#### 11. 更新地点

- **接口**: `PUT /api/locations/:id`
- **描述**: 更新地点信息（管理功能）
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id` (int64): 地点 ID

- **请求体**:

```json
{
  "name": "新名称",
  "category": "scenic",
  "latitude": 31.2306,
  "longitude": 121.4739,
  "is_active": true,
  "sort_order": 5
}
```

**字段说明**:
- 所有字段都是可选的，只更新传入的字段
- `name`: 最大 100 字符
- `category`: 最大 50 字符
- `latitude`: 纬度，范围 -90 到 90
- `longitude`: 经度，范围 -180 到 180

- **响应**:

```json
{
  "code": 200,
  "message": "success"
}
```

---

#### 12. 删除地点

- **接口**: `DELETE /api/locations/:id`
- **描述**: 删除地点（管理功能）
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id` (int64): 地点 ID

- **响应**:

```json
{
  "code": 200,
  "message": "success"
}
```

---

### 快速导航相关接口

#### 13. 获取导航树

- **接口**: `GET /api/quicknav/tree`
- **描述**: 获取校区的分类导航树
- **认证**: 无需认证
- **查询参数**:
  - `campusId` (int64, 必填): 校区 ID

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "campusId": 1,
    "campusName": "普陀校区",
    "categories": [
      {
        "category": "teaching",
        "locations": [
          {
            "id": 1,
            "name": "图书馆",
            "campus_id": 1,
            "category": "teaching",
            "latitude": 31.2304,
            "longitude": 121.4737,
            "is_active": true,
            "sort_order": 0,
            "memory_count": 25,
            "created_at": "2024-01-01T00:00:00Z",
            "updated_at": "2024-01-01T00:00:00Z"
          }
        ],
        "count": 5
      }
    ]
  }
}
```

---

#### 14. 搜索地点（快速导航）

- **接口**: `GET /api/quicknav/search`
- **描述**: 在快速导航中搜索地点
- **认证**: 无需认证
- **查询参数**:
  - `keyword` (string, 必填): 搜索关键词，1-50 字符
  - `campusId` (int64, 必填): 校区 ID
  - `category` (string, 可选): 限定类别
  - `lat` (float64, 可选): 当前位置纬度（-90 到 90）
  - `lng` (float64, 可选): 当前位置经度（-180 到 180）
  - `page` (int, 可选): 页码，默认 1
  - `pageSize` (int, 可选): 每页数量，默认 20

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "图书馆",
        "campus_id": 1,
        "category": "teaching",
        "latitude": 31.2304,
        "longitude": 121.4737,
        "memory_count": 25
      }
    ],
    "total": 10,
    "page": 1,
    "size": 20
  }
}
```

---

#### 15. 获取热门地点

- **接口**: `GET /api/quicknav/popular`
- **描述**: 获取热门地点列表
- **认证**: 无需认证
- **查询参数**:
  - `campusId` (int64, 必填): 校区 ID
  - `limit` (int, 可选): 返回数量，默认 10，范围 1-50

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "图书馆",
      "campus_id": 1,
      "category": "teaching",
      "latitude": 31.2304,
      "longitude": 121.4737,
      "memory_count": 100
    }
  ]
}
```

---

#### 16. 根据类别获取地点

- **接口**: `GET /api/quicknav/category`
- **描述**: 获取指定类别的地点列表
- **认证**: 无需认证
- **查询参数**:
  - `campusId` (int64, 必填): 校区 ID
  - `category` (string, 必填): 地点类别
  - `page` (int, 可选): 页码，默认 1
  - `pageSize` (int, 可选): 每页数量，默认 20

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "图书馆",
        "campus_id": 1,
        "category": "teaching",
        "latitude": 31.2304,
        "longitude": 121.4737,
        "memory_count": 25
      }
    ],
    "total": 15,
    "page": 1,
    "size": 20
  }
}
```

---

### 记忆相关接口

#### 17. 创建记忆

- **接口**: `POST /api/memories`
- **描述**: 创建新的记忆
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "title": "图书馆的美好时光",
  "content": "今天在图书馆学习，感觉很充实",
  "locationId": 1,
  "latitude": 31.2304,
  "longitude": 121.4737,
  "isPublic": true,
  "tags": ["学习", "图书馆"],
  "imageUrls": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ]
}
```

**字段说明**:
- `title` (string, 必填): 记忆标题，最大 100 字符
- `content` (string, 可选): 记忆内容，最大 5000 字符
- `locationId` (int64, 可选): 关联的地点 ID
- `latitude` (float64, 必填): 纬度，范围 -90 到 90
- `longitude` (float64, 必填): 经度，范围 -180 到 180
- `isPublic` (bool, 可选): 是否公开，默认 true
- `tags` (array, 可选): 标签数组
- `imageUrls` (array, 可选): 图片 URL 数组

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "title": "图书馆的美好时光",
    "content": "今天在图书馆学习，感觉很充实",
    "locationName": "图书馆",
    "locationId": 1,
    "latitude": 31.2304,
    "longitude": 121.4737,
    "creator": {
      "id": 1,
      "nickname": "小明",
      "avatar": "https://example.com/avatar.jpg"
    },
    "likeCount": 0,
    "commentCount": 0,
    "viewCount": 0,
    "isLiked": false,
    "tags": ["学习", "图书馆"],
    "images": [
      {
        "id": 1,
        "url": "https://example.com/image1.jpg"
      }
    ],
    "createdAt": "2024-01-01T10:00:00Z"
  }
}
```

---

#### 18. 获取记忆详情

- **接口**: `GET /api/memories/:id`
- **描述**: 获取指定记忆的详细信息
- **认证**: 可选认证（用于判断是否已点赞）
- **请求头**: `Authorization: Bearer {token}` (可选)
- **路径参数**:
  - `id` (int64): 记忆 ID

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "title": "图书馆的美好时光",
    "content": "今天在图书馆学习，感觉很充实",
    "locationName": "图书馆",
    "locationId": 1,
    "latitude": 31.2304,
    "longitude": 121.4737,
    "creator": {
      "id": 1,
      "nickname": "小明",
      "avatar": "https://example.com/avatar.jpg"
    },
    "likeCount": 10,
    "commentCount": 5,
    "viewCount": 100,
    "isLiked": false,
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
    "createdAt": "2024-01-01T10:00:00Z"
  }
}
```

---

#### 19. 获取记忆列表

- **接口**: `GET /api/memories`
- **描述**: 获取记忆列表（支持分页和筛选）
- **认证**: 可选认证（用于判断是否已点赞）
- **请求头**: `Authorization: Bearer {token}` (可选)
- **查询参数**:
  - `page` (int, 必填): 页码，最小值 1
  - `pageSize` (int, 必填): 每页数量，最小值 1，最大值 100
  - `locationId` (int64, 可选): 筛选指定地点的记忆
  - `sortBy` (string, 可选): 排序方式（latest-最新, popular-最热门）

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "memories": [
      {
        "id": 1,
        "title": "图书馆的美好时光",
        "content": "今天在图书馆学习...",
        "locationName": "图书馆",
        "locationId": 1,
        "latitude": 31.2304,
        "longitude": 121.4737,
        "creator": {
          "id": 1,
          "nickname": "小明",
          "avatar": "https://example.com/avatar.jpg"
        },
        "likeCount": 10,
        "commentCount": 5,
        "viewCount": 100,
        "isLiked": false,
        "tags": ["学习", "图书馆"],
        "images": [],
        "createdAt": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "pageSize": 10
  }
}
```

---

#### 20. 更新记忆

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
  "locationId": 2,
  "latitude": 31.2305,
  "longitude": 121.4738,
  "isPublic": false,
  "tags": ["新标签"],
  "imageUrls": ["https://example.com/new-image.jpg"]
}
```

**字段说明**:
- 所有字段都是可选的，只更新传入的字段
- `title`: 最大 100 字符
- `content`: 最大 5000 字符
- `locationId`: 地点 ID
- `latitude`: 纬度，范围 -90 到 90
- `longitude`: 经度，范围 -180 到 180

- **响应**:

```json
{
  "code": 200,
  "message": "success"
}
```

---

#### 21. 删除记忆

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

### 评论相关接口

#### 22. 创建评论

- **接口**: `POST /api/comments`
- **描述**: 创建评论或回复
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **请求体**:

```json
{
  "memoryId": 1,
  "content": "很棒的分享！",
  "parentId": null,
  "replyToUserId": null
}
```

**字段说明**:
- `memoryId` (int64, 必填): 记忆 ID
- `content` (string, 必填): 评论内容
- `parentId` (int64, 可选): 父评论 ID（用于回复评论）
- `replyToUserId` (int64, 可选): 被回复的用户 ID

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "memoryId": 1,
    "content": "很棒的分享！",
    "parentId": null,
    "creator": {
      "id": 1,
      "nickname": "小明",
      "avatar": "https://example.com/avatar.jpg"
    },
    "replyToUser": null,
    "likeCount": 0,
    "replyCount": 0,
    "isLiked": false,
    "createdAt": "2024-01-01T10:00:00Z"
  }
}
```

---

#### 23. 获取评论列表

- **接口**: `GET /api/comments`
- **描述**: 获取记忆的评论列表（只获取顶级评论）
- **认证**: 可选认证（用于判断是否已点赞）
- **请求头**: `Authorization: Bearer {token}` (可选)
- **查询参数**:
  - `memoryId` (int64, 必填): 记忆 ID
  - `page` (int, 必填): 页码，最小值 1
  - `pageSize` (int, 必填): 每页数量，最小值 1，最大值 100
  - `sortBy` (string, 可选): 排序方式

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "comments": [
      {
        "id": 1,
        "memoryId": 1,
        "content": "很棒的分享！",
        "parentId": null,
        "creator": {
          "id": 1,
          "nickname": "小明",
          "avatar": "https://example.com/avatar.jpg"
        },
        "replyToUser": null,
        "likeCount": 5,
        "replyCount": 3,
        "isLiked": false,
        "createdAt": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "pageSize": 20
  }
}
```

---

#### 24. 获取回复列表

- **接口**: `GET /api/comments/:parentId/replies`
- **描述**: 获取评论的回复列表
- **认证**: 可选认证（用于判断是否已点赞）
- **请求头**: `Authorization: Bearer {token}` (可选)
- **路径参数**:
  - `parentId` (int64): 父评论 ID

- **查询参数**:
  - `page` (int, 可选): 页码，默认 1
  - `pageSize` (int, 可选): 每页数量，默认 20，最大 100
  - `sortBy` (string, 可选): 排序方式

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "comments": [
      {
        "id": 2,
        "memoryId": 1,
        "content": "我也这么觉得",
        "parentId": 1,
        "creator": {
          "id": 2,
          "nickname": "小红",
          "avatar": "https://example.com/avatar2.jpg"
        },
        "replyToUser": {
          "id": 1,
          "nickname": "小明",
          "avatar": "https://example.com/avatar.jpg"
        },
        "likeCount": 2,
        "replyCount": 0,
        "isLiked": false,
        "createdAt": "2024-01-01T10:05:00Z"
      }
    ],
    "total": 3,
    "page": 1,
    "pageSize": 20
  }
}
```

---

#### 25. 删除评论

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

### 点赞相关接口

#### 26. 切换记忆点赞

- **接口**: `POST /api/memories/:id/like`
- **描述**: 切换记忆点赞状态（已点赞则取消，未点赞则点赞）
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id` (int64): 记忆 ID

- **请求体**: 无

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "isLiked": true,
    "likeCount": 11
  }
}
```

**字段说明**:
- `isLiked`: 当前点赞状态（true-已点赞, false-未点赞）
- `likeCount`: 点赞总数

---

#### 27. 切换评论点赞

- **接口**: `POST /api/comments/:id/like`
- **描述**: 切换评论点赞状态（已点赞则取消，未点赞则点赞）
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id` (int64): 评论 ID

- **请求体**: 无

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "isLiked": true,
    "likeCount": 6
  }
}
```

**字段说明**:
- `isLiked`: 当前点赞状态（true-已点赞, false-未点赞）
- `likeCount`: 点赞总数

---

### 图片相关接口

#### 28. 上传图片

- **接口**: `POST /api/images/upload`
- **描述**: 上传图片到服务器
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **Content-Type**: `multipart/form-data`
- **请求体**:
  - `file` (file, 必填): 图片文件
  - `memory_id` (int64, 必填): 关联的记忆 ID

**文件要求**:
- 支持格式: jpg, jpeg, png, gif
- 最大大小: 5MB

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "url": "/uploads/images/1234567890_image.jpg"
  }
}
```

---

#### 29. 删除图片

- **接口**: `DELETE /api/images/:id`
- **描述**: 删除图片（只能删除自己记忆的图片）
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `id` (int64): 图片 ID

- **响应**:

```json
{
  "code": 200,
  "message": "success"
}
```

---

#### 30. 获取记忆的所有图片

- **接口**: `GET /api/images/memory/:memory_id`
- **描述**: 获取指定记忆的所有图片
- **认证**: 需要 JWT 认证
- **请求头**: `Authorization: Bearer {token}`
- **路径参数**:
  - `memory_id` (int64): 记忆 ID

- **响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "url": "/uploads/images/1234567890_image1.jpg"
    },
    {
      "id": 2,
      "url": "/uploads/images/1234567891_image2.jpg"
    }
  ]
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

记忆 (memories)
  ↓ 包含
图片 (images)

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
- defaultCampusId: 默认校区ID
- status: 状态（0-禁用, 1-正常）
- role: 角色（0-普通用户, 1-管理员）
```

#### 校区 (Campus)
```
- id: 校区ID
- name: 校区名称
- isActive: 是否启用
- sortOrder: 显示顺序
- createdAt: 创建时间
- updatedAt: 更新时间
```

#### 地点 (Location)
```
- id: 地点ID
- campusId: 所属校区ID
- name: 地点名称
- category: 地点类别
- latitude: 纬度
- longitude: 经度
- isActive: 是否启用
- sortOrder: 显示顺序
- memoryCount: 记忆数量
- createdAt: 创建时间
- updatedAt: 更新时间
```

#### 记忆 (Memory)
```
- id: 记忆ID
- title: 标题
- content: 内容
- locationId: 关联地点ID
- locationName: 地点名称
- latitude: 纬度
- longitude: 经度
- creatorId: 创建者ID
- isPublic: 是否公开
- likeCount: 点赞数
- commentCount: 评论数
- viewCount: 浏览数
- tags: 标签（JSON数组）
- status: 状态（0-待审核, 1-已发布, 2-已下架）
- createdAt: 创建时间
- updatedAt: 更新时间
- deletedAt: 删除时间（软删除）
```

#### 评论 (Comment)
```
- id: 评论ID
- memoryId: 记忆ID
- content: 评论内容
- creatorId: 创建者ID
- parentId: 父评论ID（用于回复）
- replyToUserId: 被回复的用户ID
- likeCount: 点赞数
- replyCount: 回复数
- createdAt: 创建时间
- updatedAt: 更新时间
- deletedAt: 删除时间（软删除）
```

#### 点赞 (Like)
```
- id: 点赞ID
- userId: 用户ID
- targetType: 目标类型（1-记忆, 2-评论）
- targetId: 目标ID
- createdAt: 创建时间
```

#### 图片 (Image)
```
- id: 图片ID
- memoryId: 关联记忆ID
- url: 图片URL
- size: 文件大小
- createdAt: 创建时间
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
| 10001 | 请求参数错误 |
| 10002 | 未登录或 token 无效 |
| 10003 | 无权限访问（如删除他人的记忆） |
| 10004 | 资源不存在 |
| 10005 | 服务器内部错误 |
| 11xxx | 记忆相关错误 |
| 12xxx | 评论相关错误 |
| 14xxx | 点赞相关错误 |
| 15xxx | 用户相关错误 |
| 16xxx | 校区地点相关错误 |

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
- `locationId`: 必填，必须是有效的地点 ID
- `isPublic`: 可选，布尔值
- `tags`: 可选，字符串数组
- `imageUrls`: 可选，字符串数组

#### 更新用户信息
- `nickname`: 可选，1-50 字符
- `avatar`: 可选，必须是有效的 URL
- `defaultCampusId`: 可选，必须大于 0

#### 创建评论
- `memoryId`: 必填，必须是有效的记忆 ID
- `content`: 必填，评论内容
- `parentId`: 可选，父评论 ID
- `replyToUserId`: 可选，被回复的用户 ID

#### 上传图片
- 文件格式: jpg, jpeg, png, gif
- 文件大小: 最大 5MB
- `memory_id`: 必填，关联的记忆 ID

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
