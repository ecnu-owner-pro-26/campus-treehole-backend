package handler

import (
	"campus-memory/application/dto"
	"campus-memory/infra/util"
	"campus-memory/types/errno"
	"context"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommentServiceInterface 定义评论服务需要实现的方法
type CommentServiceInterface interface {
	CreateComment(ctx context.Context, req *dto.CreateCommentRequest, userID int64) (*dto.CommentResponse, error)
	ListComments(ctx context.Context, req *dto.CommentListRequest, currentUserID *int64) (*dto.CommentListResponse, error)
	ListReplies(ctx context.Context, parentID int64, page, pageSize int, currentUserID *int64) (*dto.CommentListResponse, error)
	DeleteComment(ctx context.Context, id int64, userID int64) error
}

// CommentHandler 评论处理器
type CommentHandler struct {
	commentService CommentServiceInterface
}

// NewCommentHandler 创建评论处理器实例
func NewCommentHandler(commentService CommentServiceInterface) *CommentHandler {
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
	var uri dto.ParentUri
	// 绑定路径参数
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Printf("ShouldBindUri 错误: %v", err)
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的父评论ID")
		return
	}

	var req dto.ListRepliesRequest
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
	resp, err := h.commentService.ListReplies(ctx, uri.ParentID, req.Page, req.PageSize, currentUserID)
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
