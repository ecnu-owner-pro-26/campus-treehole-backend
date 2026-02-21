package handler

import (
	"github.com/gin-gonic/gin"
)

// LikeHandler 点赞处理器
type LikeHandler struct {
	// TODO: 注入依赖（likeService）
}

// NewLikeHandler 创建点赞处理器
func NewLikeHandler() *LikeHandler {
	// TODO: 注入依赖
	return &LikeHandler{}
}

// ToggleLike 切换点赞状态（统一接口，通过路由路径判断targetType）
// 路由: POST /memories/:id/like 或 POST /comments/:id/like
func (h *LikeHandler) ToggleLike(c *gin.Context) {
	// TODO: 实现统一点赞逻辑
	// 1. 获取当前用户ID（从JWT token的context中获取）
	// 2. 从路由参数获取目标ID（c.Param("id")）
	// 3. 从路由路径判断targetType
	//    - 如果路径包含 "/memories/" → targetType = 1
	//    - 如果路径包含 "/comments/" → targetType = 2
	// 4. 调用 service.ToggleLike(userID, targetID, targetType)
	// 5. 返回响应（ToggleLikeResponse）
}
