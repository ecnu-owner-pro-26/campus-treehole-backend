# 点赞API文档

## 概述

本项目实现了独立的点赞API，支持翻转式点赞操作。用户可以通过一个接口来切换点赞状态，无需分别调用点赞和取消点赞接口。

## 特性

- **翻转式点赞**：一个接口自动判断当前状态并切换
- **统一接口**：支持记忆和留言的点赞操作
- **便捷接口**：为记忆和留言提供专门的便捷接口
- **状态查询**：可以查询点赞状态和点赞数量

## API接口

### 1. 通用翻转点赞接口

**POST** `/api/likes/toggle`

切换指定目标的点赞状态。如果已点赞则取消，如果未点赞则点赞。

**请求体：**
```json
{
  "target_id": 123,
  "target_type": 1
}
```

**参数说明：**
- `target_id`: 目标ID（记忆ID或留言ID）
- `target_type`: 目标类型（1-记忆，2-留言）

**响应：**
```json
{
  "code": 200,
  "message": "点赞成功",
  "data": {
    "is_liked": true,
    "like_count": 15
  }
}
```

### 2. 获取点赞状态接口

**GET** `/api/likes/status?target_id=123&target_type=1`

获取指定目标的点赞状态和点赞数量。

**查询参数：**
- `target_id`: 目标ID
- `target_type`: 目标类型（1-记忆，2-留言）

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "is_liked": false,
    "like_count": 14
  }
}
```

### 3. 记忆点赞便捷接口

**POST** `/api/memories/{id}/like`

翻转指定记忆的点赞状态。

**路径参数：**
- `id`: 记忆ID

**响应：**
```json
{
  "code": 200,
  "message": "取消点赞成功",
  "data": {
    "is_liked": false,
    "like_count": 13
  }
}
```

### 4. 留言点赞便捷接口

**POST** `/api/comments/{id}/like`

翻转指定留言的点赞状态。

**路径参数：**
- `id`: 留言ID

**响应：**
```json
{
  "code": 200,
  "message": "点赞成功",
  "data": {
    "is_liked": true,
    "like_count": 8
  }
}
```

## 使用示例

### JavaScript/前端示例

```javascript
// 翻转记忆点赞状态
async function toggleMemoryLike(memoryId) {
  try {
    const response = await fetch(`/api/memories/${memoryId}/like`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
    
    const result = await response.json();
    if (result.code === 200) {
      console.log(result.message); // "点赞成功" 或 "取消点赞成功"
      console.log('点赞状态:', result.data.is_liked);
      console.log('点赞数量:', result.data.like_count);
    }
  } catch (error) {
    console.error('操作失败:', error);
  }
}

// 使用通用接口翻转点赞
async function toggleLike(targetId, targetType) {
  try {
    const response = await fetch('/api/likes/toggle', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        target_id: targetId,
        target_type: targetType
      })
    });
    
    const result = await response.json();
    return result.data;
  } catch (error) {
    console.error('操作失败:', error);
  }
}

// 获取点赞状态
async function getLikeStatus(targetId, targetType) {
  try {
    const response = await fetch(`/api/likes/status?target_id=${targetId}&target_type=${targetType}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    const result = await response.json();
    return result.data;
  } catch (error) {
    console.error('获取状态失败:', error);
  }
}
```

### Go客户端示例

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type ToggleLikeRequest struct {
    TargetID   int64 `json:"target_id"`
    TargetType int8  `json:"target_type"`
}

type ToggleLikeResponse struct {
    IsLiked   bool  `json:"is_liked"`
    LikeCount int64 `json:"like_count"`
}

func toggleLike(targetID int64, targetType int8, token string) (*ToggleLikeResponse, error) {
    req := ToggleLikeRequest{
        TargetID:   targetID,
        TargetType: targetType,
    }
    
    jsonData, _ := json.Marshal(req)
    
    httpReq, _ := http.NewRequest("POST", "/api/likes/toggle", bytes.NewBuffer(jsonData))
    httpReq.Header.Set("Authorization", "Bearer "+token)
    httpReq.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    resp, err := client.Do(httpReq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result struct {
        Code    int                 `json:"code"`
        Message string              `json:"message"`
        Data    ToggleLikeResponse  `json:"data"`
    }
    
    json.NewDecoder(resp.Body).Decode(&result)
    
    if result.Code == 200 {
        return &result.Data, nil
    }
    
    return nil, fmt.Errorf(result.Message)
}
```

## 数据库表结构

```sql
CREATE TABLE likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    target_id INTEGER NOT NULL,
    target_type INTEGER NOT NULL, -- 1-记忆 2-留言
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user (user_id),
    INDEX idx_target (target_id),
    UNIQUE KEY unique_like (user_id, target_id, target_type)
);
```

## 错误码说明

- `200`: 操作成功
- `400`: 参数错误
- `401`: 用户未登录
- `500`: 服务器内部错误
- `4001`: 已经点赞过了（兼容旧接口）
- `4002`: 还没有点赞（兼容旧接口）

## 注意事项

1. 所有点赞接口都需要用户登录认证
2. 翻转接口会自动判断当前状态，无需客户端判断
3. 目标类型：1表示记忆，2表示留言
4. 便捷接口内部调用通用接口，功能完全一致
5. 支持并发操作，数据库层面有唯一约束防止重复点赞