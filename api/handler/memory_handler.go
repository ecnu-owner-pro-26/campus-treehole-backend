package handler

import (
	"campus-memory/application/dto"
	"campus-memory/application/service"
	"campus-memory/infra/util"
	"campus-memory/types/consts"
	"campus-memory/types/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MemoryHandler 记忆处理器
type MemoryHandler struct {
	memoryService *service.MemoryService
}

// NewMemoryHandler 创建记忆处理器实例
func NewMemoryHandler(memoryService *service.MemoryService) *MemoryHandler {
	return &MemoryHandler{
		memoryService: memoryService,
	}
}

// CreateMemory 创建记忆
// @Summary 创建记忆
// @Tags Memory
// @Accept json
// @Produce json
// @Param request body dto.CreateMemoryRequest true "创建记忆请求"
// @Success 200 {object} util.Response{data=dto.MemoryResponse}
// @Router /api/memories [post]
func (h *MemoryHandler) CreateMemory(c *gin.Context) {
	// 1. 获取当前用户ID(从JWT中间件获取)
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, consts.ErrCodeUnauthorized, "未登录")
		return
	}

	// 2. 绑定请求参数
	var req dto.CreateMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, consts.ErrCodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 3. 调用服务层
	memory, err := h.memoryService.CreateMemory(&req, userID.(int64))
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, consts.ErrCodeServerError, "创建失败")
		}
		return
	}

	// 4. 返回成功响应
	util.SuccessResponse(c, memory)
}

// GetMemory 获取记忆详情
// @Summary 获取记忆详情
// @Tags Memory
// @Produce json
// @Param id path int true "记忆ID"
// @Success 200 {object} util.Response{data=dto.MemoryResponse}
// @Router /api/memories/{id} [get]
func (h *MemoryHandler) GetMemory(c *gin.Context) {
	// 1. 获取记忆ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, consts.ErrCodeBadRequest, "无效的记忆ID")
		return
	}

	// 2. 获取当前用户ID(可选)
	var currentUserID *int64
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(int64)
		currentUserID = &uid
	}

	// 3. 调用服务层
	memory, err := h.memoryService.GetMemory(id, currentUserID)
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, consts.ErrCodeServerError, "获取失败")
		}
		return
	}

	// 4. 返回成功响应
	util.SuccessResponse(c, memory)
}

// ListMemories 获取记忆列表
// @Summary 获取记忆列表
// @Tags Memory
// @Produce json
// @Param location_id query int false "地点ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序方式(latest/hot)" default(latest)
// @Success 200 {object} util.Response{data=dto.MemoryListResponse}
// @Router /api/memories [get]
func (h *MemoryHandler) ListMemories(c *gin.Context) {
	// 1. 绑定查询参数
	var req dto.MemoryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		util.ErrorResponse(c, consts.ErrCodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 2. 获取当前用户ID(可选)
	var currentUserID *int64
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(int64)
		currentUserID = &uid
	}

	// 3. 调用服务层
	result, err := h.memoryService.ListMemories(&req, currentUserID)
	if err != nil {
		util.ErrorResponse(c, consts.ErrCodeServerError, "获取列表失败")
		return
	}

	// 4. 返回成功响应
	util.SuccessResponse(c, result)
}

// UpdateMemory 更新记忆
// @Summary 更新记忆
// @Tags Memory
// @Accept json
// @Produce json
// @Param id path int true "记忆ID"
// @Param request body dto.UpdateMemoryRequest true "更新记忆请求"
// @Success 200 {object} util.Response
// @Router /api/memories/{id} [put]
func (h *MemoryHandler) UpdateMemory(c *gin.Context) {
	// 1. 获取记忆ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, consts.ErrCodeBadRequest, "无效的记忆ID")
		return
	}

	// 2. 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, consts.ErrCodeUnauthorized, "未登录")
		return
	}

	// 3. 绑定请求参数
	var req dto.UpdateMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, consts.ErrCodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 4. 调用服务层
	if err := h.memoryService.UpdateMemory(id, &req, userID.(int64)); err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, consts.ErrCodeServerError, "更新失败")
		}
		return
	}

	// 5. 返回成功响应
	util.SuccessResponse(c, nil)
}

// DeleteMemory 删除记忆
// @Summary 删除记忆
// @Tags Memory
// @Produce json
// @Param id path int true "记忆ID"
// @Success 200 {object} util.Response
// @Router /api/memories/{id} [delete]
func (h *MemoryHandler) DeleteMemory(c *gin.Context) {
	// 1. 获取记忆ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, consts.ErrCodeBadRequest, "无效的记忆ID")
		return
	}

	// 2. 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, consts.ErrCodeUnauthorized, "未登录")
		return
	}

	// 3. 调用服务层
	if err := h.memoryService.DeleteMemory(id, userID.(int64)); err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, consts.ErrCodeServerError, "删除失败")
		}
		return
	}

	// 4. 返回成功响应
	util.SuccessResponse(c, nil)
}

