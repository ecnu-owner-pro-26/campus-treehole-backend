package dto

// CommentDTO 留言数据传输对象
type CommentDTO struct {
	ID        int64          `json:"id"`
	Content   string         `json:"content"`
	UserID    UserSimpleInfo `json:"creator"`
	MemoryID  int64          `json:"memory_id"`
	ParentID  int64          `json:"parent_id"`
	CreatedAt string         `json:"created_at"`
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	MemoryID int64  `json:"memoryId" binding:"required"` // 记忆ID
	Content  string `json:"content" binding:"required"`  // 评论内容
	ParentID *int64 `json:"parentId"`                    // 父评论ID，用于回复
}

// CommentListRequest 获取评论列表请求
type CommentListRequest struct {
	MemoryID int64  `form:"memoryId" binding:"required"`
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	SortBy   string `form:"sort_by"`
}

// CommentListResponse 留言列表响应
type CommentListResponse struct {
	Comments []*CommentResponse `json:"comments"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// ListRepliesRequest 获取回复列表请求
type ListRepliesRequest struct {
	ParentID int64  `uri:"parent_id" binding:"required,min=1"`            // 父评论ID（路径参数）
	Page     int    `form:"page,default=1" binding:"min=1"`               // 页码
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"` // 每页数量
	SortBy   string `form:"sort_by"`                                      // 排序方式（可选）
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID          int64           `json:"id"`
	MemoryID    int64           `json:"memoryId"`
	Content     string          `json:"content"`
	ParentID    *int64          `json:"parentId,omitempty"`
	User        UserSimpleInfo  `json:"user"`
	ReplyToUser *UserSimpleInfo `json:"replyToUser,omitempty"`
	LikeCount   int64           `json:"likeCount"`
	ReplyCount  int64           `json:"replyCount"`
	IsLiked     bool            `json:"isLiked"`
	CreatedAt   string          `json:"createdAt"`
}
