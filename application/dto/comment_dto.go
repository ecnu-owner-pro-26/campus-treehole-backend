package dto

// CommentDTO 留言数据传输对象
type CommentDTO struct {
	// TODO: 定义用于API传输的留言数据结构
	// ID        int64          `json:"id"`
	// Content   string         `json:"content"`
	// Creator   UserSimpleInfo `json:"creator"`
	// LikeCount int64          `json:"like_count"`
	// IsLiked   bool           `json:"is_liked"`    // 当前用户是否已点赞
	// CreatedAt string         `json:"created_at"`
}

// CreateCommentRequest 创建留言请求
type CreateCommentRequest struct {
	// TODO: 定义创建留言的请求参数
	// Content string `json:"content" binding:"required"` // 留言内容
}

// CommentListResponse 留言列表响应
type CommentListResponse struct {
	// TODO: 定义留言列表的响应结构
	// Comments []CommentDTO `json:"comments"`
	// Total    int64        `json:"total"`
	// Page     int          `json:"page"`
	// PageSize int          `json:"page_size"`
}
