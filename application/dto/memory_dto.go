package dto

// CreateMemoryRequest 创建记忆请求
type CreateMemoryRequest struct {
	Title      string   `json:"title" binding:"required,max=100"`
	Content    string   `json:"content" binding:"max=5000"`
	LocationID *int64   `json:"locationId"`
	Latitude   float64  `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude  float64  `json:"longitude" binding:"required,min=-180,max=180"`
	IsPublic   bool     `json:"isPublic"`
	Tags       []string `json:"tags"`
	ImageURLs  []string `json:"imageUrls"`
}

// UpdateMemoryRequest 更新记忆请求
type UpdateMemoryRequest struct {
	Title      *string  `json:"title" binding:"omitempty,max=100"`
	Content    *string  `json:"content" binding:"omitempty,max=5000"`
	LocationID *int64   `json:"locationId"`
	Latitude   *float64 `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude  *float64 `json:"longitude" binding:"omitempty,min=-180,max=180"`
	IsPublic   *bool    `json:"isPublic"`
	Tags       []string `json:"tags"`
	ImageURLs  []string `json:"imageUrls"`
}

// MemoryResponse 记忆响应
type MemoryResponse struct {
	ID           int64          `json:"id"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	LocationName string         `json:"locationName"`
	LocationID   *int64         `json:"locationId"`
	Latitude     float64        `json:"latitude"`
	Longitude    float64        `json:"longitude"`
	Creator      UserSimpleInfo `json:"creator"`
	LikeCount    int64          `json:"likeCount"`
	CommentCount int64          `json:"commentCount"`
	ViewCount    int64          `json:"viewCount"`
	IsLiked      bool           `json:"isLiked"`
	Tags         []string       `json:"tags"`
	Images       []ImageInfo    `json:"images"`
	CreatedAt    string         `json:"createdAt"`
}

// MemoryListRequest 记忆列表请求
type MemoryListRequest struct {
	LocationID *int64   `form:"locationId"`
	Page       int      `form:"page" binding:"min=1"`
	PageSize   int      `form:"pageSize" binding:"min=1,max=100"`
	SortBy     string   `form:"sortBy"`
	MinLat     *float64 `form:"min_lat"`
	MaxLat     *float64 `form:"max_lat"`
	MinLng     *float64 `form:"min_lng"`
	MaxLng     *float64 `form:"max_lng"`
	Tags       string   `form:"tags"`
}

// MemoryListResponse 记忆列表响应
type MemoryListResponse struct {
	Memories []MemoryResponse `json:"memories"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}

// UserSimpleInfo 用户简单信息
type UserSimpleInfo struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// SearchMemoriesRequest 搜索记忆请求
type SearchMemoriesRequest struct {
	Keyword  string `form:"keyword" binding:"required,max=100"` // 搜索关键词
	Page     int    `form:"page"`                               // 页码，默认1
	PageSize int    `form:"pageSize"`                           // 每页数量，默认20，最大100
	SortBy   string `form:"sortBy"`                             // 排序方式：latest（最新）、relevance（相关性，暂按时间）
	Tags     string `form:"tags"`                               // 标签筛选，逗号分隔
	TagMode  string `form:"tagMode"`                            // 标签匹配模式：or / and，默认 or
}

// SearchMemoriesResponse 搜索记忆响应
type SearchMemoriesResponse struct {
	List     []MemoryResponse `json:"list"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}
