package handler

import (
	"campus-memory/application/dto"
	"campus-memory/infra/util"
	"campus-memory/types/errno"
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LocationServiceInterface 定义地点服务需要实现的方法
type LocationServiceInterface interface {
	GetLocation(ctx context.Context, id int64) (*dto.LocationResponse, error)
	ListLocations(ctx context.Context, req *dto.LocationListRequest) (*dto.LocationListResponse, error)
	CreateLocation(ctx context.Context, req *dto.CreateLocationRequest) (*dto.LocationResponse, error)
	UpdateLocation(ctx context.Context, id int64, req *dto.UpdateLocationRequest) error
	DeleteLocation(ctx context.Context, id int64) error
	SearchLocations(ctx context.Context, keyword string, campusID *int64) ([]dto.LocationResponse, error)
}

// LocationHandler 地点处理器
type LocationHandler struct {
	locationService LocationServiceInterface
}

// NewLocationHandler 创建地点处理器实例
func NewLocationHandler(locationService LocationServiceInterface) *LocationHandler {
	return &LocationHandler{
		locationService: locationService,
	}
}

// GetLocation 获取地点详情
func (h *LocationHandler) GetLocation(c *gin.Context) {
	// 获取地点ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的地点ID")
		return
	}

	// 调用服务层
	location, err := h.locationService.GetLocation(c.Request.Context(), id)
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 返回成功响应
	util.SuccessResponse(c, location)
}

// ListLocations 获取地点列表
func (h *LocationHandler) ListLocations(c *gin.Context) {
	// 绑定查询参数
	var req dto.LocationListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "参数错误: "+err.Error())
		return
	}

	// 调用服务层
	result, err := h.locationService.ListLocations(c.Request.Context(), &req)
	if err != nil {
		util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		return
	}

	// 返回成功响应
	util.SuccessResponse(c, result)
}

// CreateLocation 创建地点
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	// 绑定请求参数
	var req dto.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "参数错误: "+err.Error())
		return
	}

	// 调用服务层
	location, err := h.locationService.CreateLocation(c.Request.Context(), &req)
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 返回成功响应
	util.SuccessResponse(c, location)
}

// UpdateLocation 更新地点
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	// 获取地点ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的地点ID")
		return
	}

	// 绑定请求参数
	var req dto.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "参数错误: "+err.Error())
		return
	}

	// 调用服务层
	if err := h.locationService.UpdateLocation(c.Request.Context(), id, &req); err != nil {
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

// DeleteLocation 删除地点
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	// 获取地点ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的地点ID")
		return
	}

	// 调用服务层
	if err := h.locationService.DeleteLocation(c.Request.Context(), id); err != nil {
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

// SearchLocations 搜索地点
func (h *LocationHandler) SearchLocations(c *gin.Context) {
	// 获取查询参数
	keyword := c.Query("keyword")
	if keyword == "" {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "搜索关键词不能为空")
		return
	}

	// 获取可选的校区ID
	var campusID *int64
	if campusIDStr := c.Query("campus_id"); campusIDStr != "" {
		id, err := strconv.ParseInt(campusIDStr, 10, 64)
		if err == nil {
			campusID = &id
		}
	}

	// 调用服务层
	locations, err := h.locationService.SearchLocations(c.Request.Context(), keyword, campusID)
	if err != nil {
		util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		return
	}

	// 返回成功响应
	util.SuccessResponse(c, locations)
}
