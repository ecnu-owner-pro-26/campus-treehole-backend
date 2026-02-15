package repo

// LikeRepo 点赞数据访问层
type LikeRepo struct {
	// TODO: 注入数据库连接
}

// NewLikeRepo 创建点赞数据访问层
func NewLikeRepo() *LikeRepo {
	// TODO: 注入数据库连接
	return &LikeRepo{}
}

// CheckLiked 检查用户是否已点赞（通过targetType区分记忆/评论）
func (r *LikeRepo) CheckLiked() (bool, error) {
	// TODO: 实现检查点赞状态
	return false, nil
}

// CreateLike 创建点赞记录（通过targetType区分记忆/评论）
func (r *LikeRepo) CreateLike() error {
	// TODO: 实现创建点赞
	return nil
}

// DeleteLike 删除点赞记录（通过targetType区分记忆/评论）
func (r *LikeRepo) DeleteLike() error {
	// TODO: 实现删除点赞（硬删除）
	return nil
}

// GetLikeCount 获取点赞数量（通过targetType区分记忆/评论）
func (r *LikeRepo) GetLikeCount() (int64, error) {
	// TODO: 实现获取点赞数量
	return 0, nil
}
