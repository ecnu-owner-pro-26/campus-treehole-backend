package dto

// CreateMemoryRequest 创建记忆请求
type CreateMemoryRequest struct {
	Title      string   `json:"title" binding:"required,max=100"`
	Content    string   `json:"content" binding:"max=5000"`
	LocationID *int64   `json:"location_id" binding:"required"`
	IsPublic   bool     `json:"is_public"`
	Tags       []string `json:"tags"`
	ImageURLs  []string `json:"image_urls"`
}

// UpdateMemoryRequest 更新记忆请求
type UpdateMemoryRequest struct {
	Title      *string  `json:"title" binding:"omitempty,max=100"`
	Content    *string  `json:"content" binding:"omitempty,max=5000"`
	LocationID *int64   `json:"location_id"`
	IsPublic   *bool    `json:"is_public"`
	Tags       []string `json:"tags"`
	ImageURLs  []string `json:"image_urls"`
}

// MemoryResponse 记忆响应
type MemoryResponse struct {
	ID           int64          `json:"id"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	LocationName string         `json:"location_name"`
	LocationID   *int64         `json:"location_id"`
	Creator      UserSimpleInfo `json:"creator"`
	LikeCount    int64          `json:"like_count"`
	CommentCount int64          `json:"comment_count"`
	ViewCount    int64          `json:"view_count"`
	IsLiked      bool           `json:"is_liked"`
	Tags         []string       `json:"tags"`
	Images       []ImageInfo    `json:"images"`
	CreatedAt    string         `json:"created_at"`
}

// MemoryListRequest 记忆列表请求
type MemoryListRequest struct {
	LocationID *int64 `form:"location_id"`
	Page       int    `form:"page" binding:"min=1"`
	PageSize   int    `form:"page_size" binding:"min=1,max=100"`
	SortBy     string `form:"sort_by"`
}

// MemoryListResponse 记忆列表响应
type MemoryListResponse struct {
	Memories []MemoryResponse `json:"memories"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

// UserSimpleInfo 用户简单信息
type UserSimpleInfo struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// ImageInfo 图片信息
type ImageInfo struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}