package handler

import (
	"campus-memory/application/dto"
	"campus-memory/application/service"
	"campus-memory/infra/util"
	"campus-memory/types/errno"

	"github.com/gin-gonic/gin"
)

// QuickNavHandler 快速导航处理器
type QuickNavHandler struct {
	quicknavService *service.QuickNavService
}

// NewQuickNavHandler 创建处理器实例
func NewQuickNavHandler(quicknavService service.QuickNavService) *QuickNavHandler {
	return &QuickNavHandler{
		quicknavService: &quicknavService,
	}
}

// GetNavTree 获取导航树
func (h *QuickNavHandler) GetNavTree(c *gin.Context) {
	// 接收请求参数
	var req dto.GetNavTreeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "参数错误: "+err.Error())
		return
	}

	// 调用Service层
	result, err := h.quicknavService.BuildNavTree(req.CampusID)
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

// GetLocationsByCategory 根据类别获取地点
func (h *QuickNavHandler) GetLocationsByCategory(c *gin.Context) {
	// 接收请求参数
	var req dto.GetLocationsByCategoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "参数错误: "+err.Error())
		return
	}

	// 设置默认分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 调用Service
	locations, total, err := h.quicknavService.GetLocationsByCategory(req.CampusID, req.Category, req.Page, req.PageSize)
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 统一的分页响应格式
	util.SuccessResponse(c, gin.H{
		"list":  locations,
		"total": total,
		"page":  req.Page,
		"size":  req.PageSize,
	})
}

// SearchLocations 搜索地点
func (h *QuickNavHandler) SearchLocations(c *gin.Context) {
	// 使用 DTO 接收请求参数
	var req dto.LocationSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "参数错误: "+err.Error())
		return
	}

	// 设置默认分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 解析校区ID
	campusID := req.CampusID

	// 调用service层
	locations, total, err := h.quicknavService.SearchLocations(req.Keyword, campusID, req.Page, req.PageSize)
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 统一的分页响应格式
	util.SuccessResponse(c, gin.H{
		"list":  locations,
		"total": total,
		"page":  req.Page,
		"size":  req.PageSize,
	})
}

// GetPopularLocations 获取热门地点
func (h *QuickNavHandler) GetPopularLocations(c *gin.Context) {
	var req dto.GetPopularLocationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "参数错误: "+err.Error())
		return
	}

	// 调用服务层
	locations, err := h.quicknavService.GetPopularLocations(req.CampusID, req.Limit)
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 成功响应
	util.SuccessResponse(c, locations)
}
