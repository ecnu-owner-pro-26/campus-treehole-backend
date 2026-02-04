package dto

// CreateMemoryRequest 创建记忆请求
type CreateMemoryRequest struct {
	// TODO: 定义创建记忆请求字段
	// Title      string   `json:"title" binding:"required"`       // 记忆标题
	// Content    string   `json:"content"`                        // 记忆内容
	// LocationID *int64   `json:"location_id" binding:"required"` // 地点ID（来自快速导航选择）
	// IsPublic   bool     `json:"is_public"`                      // 是否公开
	// Tags       []string `json:"tags"`                           // 标签列表
}

// MemoryResponse 记忆响应
type MemoryResponse struct {
	// TODO: 定义记忆响应字段
	// ID           int64            `json:"id"`
	// Title        string           `json:"title"`
	// Content      string           `json:"content"`
	// LocationName string           `json:"location_name"`
	// LocationID   *int64           `json:"location_id"`
	// Creator      UserSimpleInfo   `json:"creator"`
	// LikeCount    int64            `json:"like_count"`
	// CommentCount int64            `json:"comment_count"`
	// ViewCount    int64            `json:"view_count"`
	// IsLiked      bool             `json:"is_liked"`      // 当前用户是否已点赞
	// Tags         []string         `json:"tags"`
	// CreatedAt    string           `json:"created_at"`
}

// MemoryListResponse 记忆列表响应
type MemoryListResponse struct {
	// TODO: 定义记忆列表响应字段
	// Memories []MemoryResponse `json:"memories"`
	// Total    int64            `json:"total"`
	// Page     int              `json:"page"`
	// PageSize int              `json:"page_size"`
}

// UserSimpleInfo 用户简单信息
type UserSimpleInfo struct {
	// TODO: 定义用户简单信息字段
	// ID       int64  `json:"id"`
	// Nickname string `json:"nickname"`
	// Avatar   string `json:"avatar"`
}