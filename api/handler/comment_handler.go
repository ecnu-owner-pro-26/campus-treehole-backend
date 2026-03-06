package handler

import (
	"campus-memory/application/dto"
	"campus-memory/application/service"
	"campus-memory/infra/util"
	"campus-memory/types/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	commentService *service.CommentService
}

// NewCommentHandler 创建评论处理器实例
func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateComment 创建评论
func (h *CommentHandler) CreateComment(c *gin.Context) {
	// 从context获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, errno.ErrUnauthorized.Code, errno.ErrUnauthorized.Message)
		return
	}

	// 绑定请求参数
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "参数错误: "+err.Error())
		return
	}

	// 调用服务层
	ctx := c.Request.Context()
	resp, err := h.commentService.CreateComment(ctx, &req, userID.(int64))
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 返回成功响应
	util.SuccessResponse(c, resp)
}

// ListComments 获取评论列表（只获取顶级评论）
func (h *CommentHandler) ListComments(c *gin.Context) {
	// 绑定查询参数
	var req dto.CommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "查询参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID(可选，用于判断是否点赞等)
	var currentUserID *int64
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(int64)
		currentUserID = &uid
	}

	// 调用服务层
	resp, err := h.commentService.ListComments(c.Request.Context(), &req, currentUserID)
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 返回成功响应
	util.SuccessResponse(c, resp)
}

// ListReplies 获取评论的回复列表
func (h *CommentHandler) ListReplies(c *gin.Context) {
	var req dto.ListRepliesRequest
	// 绑定路径参数
	if err := c.ShouldBindUri(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的父评论ID")
		return
	}

	// 绑定查询参数
	if err := c.ShouldBindQuery(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "分页参数错误: "+err.Error())
		return
	}

	// 获取当前用户ID（可选）
	var currentUserID *int64
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(int64); ok {
			currentUserID = &uid
		}
	}

	// 调用service层
	ctx := c.Request.Context()
	resp, err := h.commentService.ListReplies(ctx, req.ParentID, req.Page, req.PageSize, currentUserID)
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, "获取回复列表失败")
		}
		return
	}

	// 返回成功响应
	util.SuccessResponse(c, resp)
}

// DeleteComment 删除评论
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	// 获取评论ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, errno.ErrBadRequest.Message)
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, errno.ErrUnauthorized.Code, errno.ErrUnauthorized.Message)
		return
	}

	// 调用服务层
	if err := h.commentService.DeleteComment(c.Request.Context(), id, userID.(int64)); err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 返回成功响应
	util.SuccessResponse(c, nil)
}
