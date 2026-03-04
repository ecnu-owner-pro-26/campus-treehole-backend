package dto

// CommentDTO 留言数据传输对象（已废弃，使用 CommentResponse）
type CommentDTO struct {
	ID        int64          `json:"id"`
	Content   string         `json:"content"`
	UserID    UserSimpleInfo `json:"creator"`
	MemoryID  int64          `json:"memoryId"`
	ParentID  int64          `json:"parentId"`
	CreatedAt string         `json:"createdAt"`
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	MemoryID      int64  `json:"memoryId" binding:"required"`      // 记忆ID
	Content       string `json:"content" binding:"required"`       // 评论内容
	ParentID      *int64 `json:"parentId"`                         // 父评论ID，用于回复
	ReplyToUserID *int64 `json:"replyToUserId"`                    // 被回复的用户ID
}

// CommentListRequest 获取评论列表请求
type CommentListRequest struct {
	MemoryID int64  `form:"memoryId" binding:"required"`
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"pageSize" binding:"min=1,max=100"`
	SortBy   string `form:"sortBy"`
}

// CommentListResponse 留言列表响应
type CommentListResponse struct {
	Comments []*CommentResponse `json:"comments"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
}

// ListRepliesRequest 获取回复列表请求
type ListRepliesRequest struct {
	ParentID int64  `uri:"parentId" binding:"required,min=1"`           // 父评论ID（路径参数）
	Page     int    `form:"page,default=1" binding:"min=1"`             // 页码
	PageSize int    `form:"pageSize,default=20" binding:"min=1,max=100"` // 每页数量
	SortBy   string `form:"sortBy"`                                     // 排序方式（可选）
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID          int64           `json:"id"`
	MemoryID    int64           `json:"memoryId"`
	Content     string          `json:"content"`
	ParentID    *int64          `json:"parentId,omitempty"`
	Creator     UserSimpleInfo  `json:"creator"`           // 改为 creator 与文档一致
	ReplyToUser *UserSimpleInfo `json:"replyToUser,omitempty"`
	LikeCount   int64           `json:"likeCount"`
	ReplyCount  int64           `json:"replyCount"`
	IsLiked     bool            `json:"isLiked"`
	CreatedAt   string          `json:"createdAt"`
}
