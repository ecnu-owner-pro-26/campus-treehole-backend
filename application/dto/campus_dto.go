package dto

// CampusResponse 校区响应
type CampusResponse struct {
	// TODO: 定义校区响应字段
	// ID          int64  `json:"id"`
	// Name        string `json:"name"`        // 校区名称：普陀校区
	// Description string `json:"description"` // 校区描述
}

// LocationResponse 地点响应
type LocationResponse struct {
	// TODO: 定义地点响应字段
	// ID          int64  `json:"id"`
	// Name        string `json:"name"`         // 地点名称：图书馆
	// Category    string `json:"category"`     // 地点类别：library
	// Address     string `json:"address"`      // 详细地址
	// Description string `json:"description"`  // 地点描述
	// Icon        string `json:"icon"`         // 地点图标
	// MemoryCount int64  `json:"memory_count"` // 记忆数量
}

// CampusLocationsResponse 校区地点列表响应
type CampusLocationsResponse struct {
	// TODO: 定义校区地点列表响应字段
	// Campus    CampusResponse     `json:"campus"`
	// Locations []LocationResponse `json:"locations"`
}

// NavTreeResponse 导航树响应
type NavTreeResponse struct {
	// TODO: 定义导航树响应字段
	// Campuses []CampusWithLocations `json:"campuses"`
}

// CampusWithLocations 带地点的校区
type CampusWithLocations struct {
	// TODO: 定义带地点的校区结构
	// Campus    CampusResponse     `json:"campus"`
	// Locations []LocationResponse `json:"locations"`
}