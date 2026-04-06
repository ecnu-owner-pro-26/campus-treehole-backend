package handler

import (
	"campus-memory/application/dto"
	"campus-memory/infra/util"
	"campus-memory/types/errno"
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CampusServiceInterface 定义校区服务需要实现的方法
type CampusServiceInterface interface {
	ListCampuses(ctx context.Context) (*dto.CampusListResponse, error)
	GetCampus(ctx context.Context, id int64) (*dto.CampusResponse, error)
	GetCampusWithLocations(ctx context.Context, campusID int64) (*dto.CampusLocationsResponse, error)
}

// CampusHandler 校区处理器
type CampusHandler struct {
	campusService CampusServiceInterface
}

// NewCampusHandler 创建校区处理器实例
func NewCampusHandler(campusService CampusServiceInterface) *CampusHandler {
	return &CampusHandler{
		campusService: campusService,
	}
}

// ListCampuses 获取所有校区列表
func (h *CampusHandler) ListCampuses(c *gin.Context) {
	// 调用服务层
	result, err := h.campusService.ListCampuses(c.Request.Context())
	if err != nil {
		util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		return
	}

	// 返回成功响应
	util.SuccessResponse(c, result)
}

// GetCampus 获取校区详情
func (h *CampusHandler) GetCampus(c *gin.Context) {
	// 获取校区ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的校区ID")
		return
	}

	// 调用服务层
	campus, err := h.campusService.GetCampus(c.Request.Context(), id)
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 返回成功响应
	util.SuccessResponse(c, campus)
}

// GetCampusWithLocations 获取校区及其地点列表
func (h *CampusHandler) GetCampusWithLocations(c *gin.Context) {
	// 获取校区ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的校区ID")
		return
	}

	// 调用服务层
	result, err := h.campusService.GetCampusWithLocations(c.Request.Context(), id)
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 返回成功响应
	util.SuccessResponse(c, result)
}
