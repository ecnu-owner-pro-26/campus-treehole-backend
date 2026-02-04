package dto

// ToggleLikeRequest 切换点赞请求
type ToggleLikeRequest struct {
	TargetID   int64 `json:"target_id" binding:"required"`   // 目标ID
	TargetType int8  `json:"target_type" binding:"required"` // 目标类型：1-记忆 2-留言
}

// ToggleLikeResponse 切换点赞响应
type ToggleLikeResponse struct {
	IsLiked   bool  `json:"is_liked"`   // 是否已点赞
	LikeCount int64 `json:"like_count"` // 点赞总数
}

// LikeStatusResponse 点赞状态响应
type LikeStatusResponse struct {
	IsLiked   bool  `json:"is_liked"`   // 是否已点赞
	LikeCount int64 `json:"like_count"` // 点赞总数
}