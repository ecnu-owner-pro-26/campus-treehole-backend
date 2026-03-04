package handler

import (
	"campus-memory/application/dto"
	"campus-memory/infra/util"
	"campus-memory/types/errno"
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MemoryServiceInterface 定义记忆服务需要实现的方法
type MemoryServiceInterface interface {
	CreateMemory(ctx context.Context, req *dto.CreateMemoryRequest, creatorID int64) (*dto.MemoryResponse, error)
	GetMemory(ctx context.Context, id int64, currentUserID *int64) (*dto.MemoryResponse, error)
	ListMemories(ctx context.Context, req *dto.MemoryListRequest, currentUserID *int64) (*dto.MemoryListResponse, error)
	UpdateMemory(ctx context.Context, id int64, req *dto.UpdateMemoryRequest, userID int64) error
	DeleteMemory(ctx context.Context, id int64, userID int64) error
}

// MemoryHandler 记忆处理器
type MemoryHandler struct {
	memoryService MemoryServiceInterface
}

// NewMemoryHandler 创建记忆处理器实例
func NewMemoryHandler(memoryService MemoryServiceInterface) *MemoryHandler {
	return &MemoryHandler{
		memoryService: memoryService,
	}
}

// CreateMemory 创建记忆
func (h *MemoryHandler) CreateMemory(c *gin.Context) {
	// 1. 获取当前用户ID(从JWT中间件获取)
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, errno.ErrUnauthorized.Code, errno.ErrUnauthorized.Message)
		return
	}

	// 2. 绑定请求参数
	var req dto.CreateMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "参数错误: "+err.Error())
		return
	}

	// 3. 调用服务层
	memory, err := h.memoryService.CreateMemory(c.Request.Context(), &req, userID.(int64))
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 4. 返回成功响应
	util.SuccessResponse(c, memory)
}

// GetMemory 获取记忆详情
func (h *MemoryHandler) GetMemory(c *gin.Context) {
	// 1. 获取记忆ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的记忆ID")
		return
	}

	// 2. 获取当前用户ID(可选)
	var currentUserID *int64
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(int64)
		currentUserID = &uid
	}

	// 3. 调用服务层
	memory, err := h.memoryService.GetMemory(c.Request.Context(), id, currentUserID)
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 4. 返回成功响应
	util.SuccessResponse(c, memory)
}

// ListMemories 获取记忆列表
func (h *MemoryHandler) ListMemories(c *gin.Context) {
	// 1. 绑定查询参数
	var req dto.MemoryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "参数错误: "+err.Error())
		return
	}

	// 2. 获取当前用户ID(可选)
	var currentUserID *int64
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(int64)
		currentUserID = &uid
	}

	// 3. 调用服务层
	result, err := h.memoryService.ListMemories(c.Request.Context(), &req, currentUserID)
	if err != nil {
		util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		return
	}

	// 4. 返回成功响应
	util.SuccessResponse(c, result)
}

// UpdateMemory 更新记忆
func (h *MemoryHandler) UpdateMemory(c *gin.Context) {
	// 1. 获取记忆ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的记忆ID")
		return
	}

	// 2. 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, errno.ErrUnauthorized.Code, errno.ErrUnauthorized.Message)
		return
	}

	// 3. 绑定请求参数
	var req dto.UpdateMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "参数错误: "+err.Error())
		return
	}

	// 4. 调用服务层
	if err := h.memoryService.UpdateMemory(c.Request.Context(), id, &req, userID.(int64)); err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 5. 返回成功响应
	util.SuccessResponse(c, nil)
}

// DeleteMemory 删除记忆
func (h *MemoryHandler) DeleteMemory(c *gin.Context) {
	// 1. 获取记忆ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的记忆ID")
		return
	}

	// 2. 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, errno.ErrUnauthorized.Code, errno.ErrUnauthorized.Message)
		return
	}

	// 3. 调用服务层
	if err := h.memoryService.DeleteMemory(c.Request.Context(), id, userID.(int64)); err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 4. 返回成功响应
	util.SuccessResponse(c, nil)
}
