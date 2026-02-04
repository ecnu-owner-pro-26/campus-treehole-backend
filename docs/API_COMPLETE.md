# 校园树洞完整API文档

## 用户使用流程API

### 1. 快速导航 - 选择校区和地点

#### 获取校区列表
- **接口**: `GET /api/campuses`
- **描述**: 用户打开应用时调用，获取所有可用校区
- **认证**: 无需认证
- **响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "普陀校区",
      "description": "上海海洋大学普陀校区"
    }
  ]
}
```

#### 获取校区地点列表
- **接口**: `GET /api/campuses/:id/locations`
- **描述**: 用户选择校区后，获取该校区下的所有地点
- **认证**: 无需认证
- **响应示例**:
```json
{
  "code": 200,
  "message": "success", 
  "data": [
    {
      "id": 1,
      "name": "图书馆",
      "category": "library",
      "address": "普陀校区图书馆大楼",
      "description": "普陀校区图书馆",
      "icon": "library_books",
      "memory_count": 25
    }
  ]
}
```

### 2. 记忆管理 - 发表和浏览记忆

#### 创建记忆
- **接口**: `POST /api/memories`
- **描述**: 用户在选定地点发表记忆
- **认证**: 需要JWT认证
- **请求参数**:
```json
{
  "title": "图书馆的美好时光",
  "content": "今天在图书馆学习，感觉很充实",
  "location_id": 1,
  "is_public": true,
  "tags": ["学习", "图书馆"]
}
```

#### 获取地点记忆列表
- **接口**: `GET /api/locations/:id/memories`
- **描述**: 浏览指定地点的所有记忆
- **认证**: 可选（影响是否显示私密记忆）
- **响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "title": "图书馆的美好时光",
      "content": "今天在图书馆学习...",
      "location_name": "图书馆",
      "creator": {
        "id": 1,
        "nickname": "小明"
      },
      "like_count": 5,
      "comment_count": 3,
      "created_at": "2024-01-01T10:00:00Z"
    }
  ]
}
```

#### 获取记忆详情
- **接口**: `GET /api/memories/:id`
- **描述**: 查看记忆详细内容
- **认证**: 可选

### 3. 点赞功能

#### 点赞记忆
- **接口**: `POST /api/memories/:id/like`
- **描述**: 给记忆点赞
- **认证**: 需要JWT认证

#### 取消点赞记忆
- **接口**: `DELETE /api/memories/:id/like`
- **描述**: 取消记忆点赞
- **认证**: 需要JWT认证

#### 点赞评论
- **接口**: `POST /api/comments/:id/like`
- **描述**: 给别人的评论点赞
- **认证**: 需要JWT认证

#### 取消点赞评论
- **接口**: `DELETE /api/comments/:id/like`
- **描述**: 取消评论点赞
- **认证**: 需要JWT认证

### 4. 评论功能

#### 创建评论
- **接口**: `POST /api/memories/:id/comments`
- **描述**: 给记忆添加评论
- **认证**: 需要JWT认证
- **请求参数**:
```json
{
  "content": "很棒的分享！"
}
```

#### 获取记忆评论列表
- **接口**: `GET /api/memories/:id/comments`
- **描述**: 获取记忆的所有评论
- **认证**: 可选

#### 删除评论
- **接口**: `DELETE /api/comments/:id`
- **描述**: 删除自己的评论
- **认证**: 需要JWT认证

### 5. 用户认证

#### 用户注册
- **接口**: `POST /api/auth/register`
- **描述**: 用户注册账号
- **认证**: 无需认证

#### 用户登录
- **接口**: `POST /api/auth/login`
- **描述**: 用户登录获取JWT token
- **认证**: 无需认证

#### 获取个人信息
- **接口**: `GET /api/auth/profile`
- **描述**: 获取当前用户信息
- **认证**: 需要JWT认证

## 完整用户流程示例

```
1. 用户打开应用
   GET /api/campuses
   
2. 选择"普陀校区"
   GET /api/campuses/1/locations
   
3. 选择"图书馆"
   GET /api/locations/1/memories (浏览记忆)
   或
   POST /api/memories (发表记忆)
   
4. 与记忆互动
   POST /api/memories/1/like (点赞记忆)
   POST /api/memories/1/comments (评论记忆)
   POST /api/comments/1/like (点赞评论)
```

## 数据库关联关系

- 用户(users) → 记忆(memories) → 地点(locations) → 校区(campuses)
- 记忆(memories) → 评论(comments) → 用户(users)
- 点赞(likes) → 记忆/评论 + 用户