package handler

import (
	"campus-memory/application/service"
	"campus-memory/infra/util"
	"campus-memory/types/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CampusHandler 校区处理器
type CampusHandler struct {
	campusService *service.CampusService
}

// NewCampusHandler 创建校区处理器实例
func NewCampusHandler(campusService *service.CampusService) *CampusHandler {
	return &CampusHandler{
		campusService: campusService,
	}
}

// ListCampuses 获取所有校区列表
func (h *CampusHandler) ListCampuses(c *gin.Context) {
	// 调用服务层
	result, err := h.campusService.ListCampuses()
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
	campus, err := h.campusService.GetCampus(id)
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
	result, err := h.campusService.GetCampusWithLocations(id)
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
