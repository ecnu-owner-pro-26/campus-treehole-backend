package repo

// CommentRepo 留言仓储层
type CommentRepo struct {
	// TODO: 注入数据库连接
}

// NewCommentRepo 创建留言仓储
func NewCommentRepo() *CommentRepo {
	// TODO: 初始化仓储
	return &CommentRepo{}
}

// Create 创建留言记录
func (r *CommentRepo) Create() error {
	// TODO: 实现数据库插入逻辑
	return nil
}

// ListByMemoryID 根据记忆ID获取留言列表
func (r *CommentRepo) ListByMemoryID() error {
	// TODO: 实现数据库查询逻辑
	return nil
}

// Delete 删除留言记录
func (r *CommentRepo) Delete() error {
	// TODO: 实现数据库删除逻辑
	return nil
}
