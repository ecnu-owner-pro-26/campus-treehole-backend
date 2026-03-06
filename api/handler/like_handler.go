package handler

import (
	"campus-memory/application/dto"
	"campus-memory/application/service"
	"campus-memory/infra/util"
	"campus-memory/types/errno"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// LikeHandler 点赞处理器
type LikeHandler struct {
	likeService *service.LikeService
}

// NewLikeHandler 创建点赞处理器
func NewLikeHandler(likeService *service.LikeService) *LikeHandler {
	return &LikeHandler{
		likeService: likeService,
	}
}

// ToggleLike 切换点赞状态（统一接口，通过路由路径判断targetType）
// 路由: POST /memories/:id/like 或 POST /comments/:id/like
func (h *LikeHandler) ToggleLike(c *gin.Context) {
	// 1. 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, errno.ErrUnauthorized.Code, errno.ErrUnauthorized.Message)
		return
	}

	// 2. 获取目标ID
	idStr := c.Param("id")
	targetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的ID")
		return
	}

	// 3. 从路由路径判断targetType
	var targetType int8
	path := c.Request.URL.Path
	if strings.Contains(path, "/memories/") {
		targetType = dto.LikeTargetTypeMemory
	} else if strings.Contains(path, "/comments/") {
		targetType = dto.LikeTargetTypeComment
	} else {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的点赞目标")
		return
	}

	// 4. 调用服务层
	result, err := h.likeService.ToggleLike(userID.(int64), targetID, targetType)
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 5. 返回响应
	util.SuccessResponse(c, result)
}
