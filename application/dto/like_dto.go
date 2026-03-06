package dto

// 点赞目标类型常量
const (
	LikeTargetTypeMemory  int8 = 1 // 记忆
	LikeTargetTypeComment int8 = 2 // 评论
)

// ToggleLikeResponse 切换点赞响应
type ToggleLikeResponse struct {
	IsLiked   bool  `json:"isLiked"`   // 当前点赞状态（true-已点赞 false-未点赞）
	LikeCount int64 `json:"likeCount"` // 点赞总数
}
