package dto

// LocationResponse 地点响应
type LocationResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	CampusID    int64  `json:"campus_id"`
	Category    string `json:"category"`
	MemoryCount int64  `json:"memory_count"`
	CreatedAt   string `json:"created_at"`
}

// LocationListRequest 地点列表请求
type LocationListRequest struct {
	CampusID *int64 `form:"campus_id"`
	Category string `form:"category"`
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
}

// LocationListResponse 地点列表响应
type LocationListResponse struct {
	Locations []LocationResponse `json:"locations"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
}

// CreateLocationRequest 创建地点请求
type CreateLocationRequest struct {
	CampusID  int64  `json:"campus_id" binding:"required"`
	Name      string `json:"name" binding:"required,max=100"`
	Category  string `json:"category" binding:"max=50"`
	SortOrder int    `json:"sort_order"`
}

// UpdateLocationRequest 更新地点请求
type UpdateLocationRequest struct {
	Name      *string `json:"name" binding:"omitempty,max=100"`
	Category  *string `json:"category" binding:"omitempty,max=50"`
	IsActive  *bool   `json:"is_active"`
	SortOrder *int    `json:"sort_order"`
}
