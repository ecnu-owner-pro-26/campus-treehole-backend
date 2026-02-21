package repo

// CommentRepo 评论数据访问层
type CommentRepo struct {
	// TODO: 注入数据库连接
}

// Create 创建评论记录
func (r *CommentRepo) Create() error {
	// TODO: 实现数据库插入逻辑
	return nil
}

// ListByMemoryID 根据记忆ID获取评论列表
func (r *CommentRepo) ListByMemoryID() error {
	// TODO: 实现数据库查询逻辑
	return nil
}

// Delete 删除评论记录
func (r *CommentRepo) Delete() error {
	// TODO: 实现数据库删除逻辑
	return nil
}

