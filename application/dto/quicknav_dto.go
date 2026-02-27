package dto

import "campus-memory/infra/model"

// GetNavTreeRequest 导航树请求参数
type GetNavTreeRequest struct {
	CampusID int64 `form:"campusID"` // 可选，指定校区
}

// GetLocationsByCategoryRequest 按类别查询请求
type GetLocationsByCategoryRequest struct {
	CampusID int64  `form:"campus_id" binding:"required"` // 校区ID
	Category string `form:"category" binding:"required"`  // 类别
	Page     int    `form:"page,default=1"`               // 页码
	PageSize int    `form:"page_size,default=20"`         // 每页数量
}

// LocationSearchRequest 地点搜索请求
type LocationSearchRequest struct {
	Keyword  string  `form:"keyword" binding:"required"` // 搜索关键词
	Campus   string  `form:"campus"`                     // 可选，限定校区
	Category string  `form:"category"`                   // 可选，限定类别
	Lat      float64 `form:"lat"`                        // 可选，当前位置纬度
	Lng      float64 `form:"lng"`                        // 可选，当前位置经度
	Page     int     `form:"page,default=1"`             // 页码
	PageSize int     `form:"page_size,default=20"`       // 每页数量
}

// CategoryNodeResult 类别节点结果
type CategoryNodeResult struct {
	Category  string                `json:"category"`
	Locations []model.LocationModel `json:"locations"`
	Count     int                   `json:"count"`
}
