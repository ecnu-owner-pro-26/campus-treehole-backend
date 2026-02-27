package handler

import (
	"campus-memory/application/dto"
	"campus-memory/application/service"
	"campus-memory/infra/util"
	"campus-memory/types/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

// QuickNavHandler 快速导航处理器
type QuickNavHandler struct {
	quicknavService service.QuickNavService
	campusService   service.CampusService
}

// NewQuickNavHandler 创建处理器实例
func NewQuickNavHandler(quicknaveService service.QuickNavService) *QuickNavHandler {
	return &QuickNavHandler{
		quicknavService: quicknaveService,
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

	// 验证校区ID
	if req.CampusID <= 0 {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "校区ID不能为空")
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

	// 参数验证
	if req.CampusID <= 0 {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "校区ID不能为空")
		return
	}
	if req.Category == "" {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "类别不能为空")
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
	var campusID int64
	if req.Campus != "" {
		var err error
		campusID, err = strconv.ParseInt(req.Campus, 10, 64)
		if err != nil {
			util.ErrorResponse(c, errno.ErrBadRequest.Code, "校区ID格式错误")
			return
		}
	}

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
	// 解析校区ID
	campusID, err := strconv.ParseInt(c.Query("campus_id"), 10, 64)
	if err != nil || campusID <= 0 {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的校区ID")
		return
	}

	// 解析并限制数量
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if limit <= 0 || limit > 50 {
		limit = 10
	}

	// 调用服务层
	locations, err := h.quicknavService.GetPopularLocations(campusID, limit)
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
