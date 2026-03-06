package dto

import "campus-memory/infra/model"

// GetNavTreeRequest 导航树请求参数
type GetNavTreeRequest struct {
	CampusID int64 `form:"campusId" binding:"required,min=1"` // 校区ID，必填且为正整数
}

// CategoryDTO 类别传输对象
type CategoryDTO struct {
	Category  string                `json:"category"`
	Locations []model.LocationModel `json:"locations"`
	Count     int                   `json:"count"`
}

// CampusNavDTO 校区导航传输对象
type CampusNavDTO struct {
	CampusID   int64         `json:"campusId"`
	CampusName string        `json:"campusName"`
	Categories []CategoryDTO `json:"categories"`
}

// GetLocationsByCategoryRequest 按类别查询请求
type GetLocationsByCategoryRequest struct {
	CampusID int64  `form:"campusId" binding:"required,min=1"` // 校区ID，必填且为正整数
	Category string `form:"category" binding:"required"`       // 类别
	Page     int    `form:"page,default=1"`                    // 页码
	PageSize int    `form:"pageSize,default=20"`               // 每页数量
}

// LocationSearchRequest 地点搜索请求
type LocationSearchRequest struct {
	Keyword  string  `form:"keyword" binding:"required,min=1,max=50"`      // 搜索关键词
	CampusID int64   `form:"campusId" binding:"required,min=1"`            // 校区ID，必填且为正整数
	Category string  `form:"category"`                                     // 可选，限定类别
	Lat      float64 `form:"lat" binding:"omitempty,min=-90,max=90"`       // 可选，当前位置纬度
	Lng      float64 `form:"lng" binding:"omitempty,min=-180,max=180"`     // 可选，当前位置经度
	Page     int     `form:"page,default=1"`                               // 页码
	PageSize int     `form:"pageSize,default=20"`                          // 每页数量
}

// GetPopularLocationsRequest 获取热门地点请求参数
type GetPopularLocationsRequest struct {
	CampusID int64 `form:"campusId" binding:"required,min=1"`   // 校区ID，必填且为正整数
	Limit    int   `form:"limit,default=10" binding:"min=1,max=50"` // 返回数量，默认10，范围1-50
}

// CategoryNodeResult 类别节点结果
type CategoryNodeResult struct {
	Category  string                `json:"category"`
	Locations []model.LocationModel `json:"locations"`
	Count     int                   `json:"count"`
}
