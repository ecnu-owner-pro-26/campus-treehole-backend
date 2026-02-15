package service

// LikeService 点赞服务
type LikeService struct {
	// TODO: 注入依赖（likeRepo, memoryRepo, commentRepo）
}

// NewLikeService 创建点赞服务
func NewLikeService() *LikeService {
	// TODO: 注入依赖
	return &LikeService{}
}

// ToggleLike 切换点赞状态（统一方法，通过targetType区分记忆/评论）
func (s *LikeService) ToggleLike() error {
	// TODO: 实现统一的点赞业务逻辑
	// 1. 根据targetType检查目标是否存在
	//    - targetType=1 → 调用 memoryRepo.ExistsByID(targetID)
	//    - targetType=2 → 调用 commentRepo.ExistsByID(targetID)
	// 2. 调用 likeRepo.CheckLiked(userID, targetID, targetType)
	// 3. 翻转逻辑：
	//    - 如果已点赞 → likeRepo.DeleteLike(userID, targetID, targetType)
	//    - 如果未点赞 → likeRepo.CreateLike(userID, targetID, targetType)
	// 4. 调用 likeRepo.GetLikeCount(targetID, targetType)
	// 5. 返回 ToggleLikeResponse{IsLiked: !原状态, LikeCount: 数量}
	return nil
}
