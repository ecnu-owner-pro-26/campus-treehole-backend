package repo

// LocationRepo 地点数据访问层
type LocationRepo struct {
	// TODO: 注入数据库连接
}

// GetByCampusID 根据校区ID获取地点列表
func (r *LocationRepo) GetByCampusID() error {
	// TODO: 实现根据校区ID获取地点列表
	return nil
}

// GetByCategory 根据类别获取地点
func (r *LocationRepo) GetByCategory() error {
	// TODO: 实现根据类别获取地点
	return nil
}

// Search 搜索地点
func (r *LocationRepo) Search() error {
	// TODO: 实现搜索地点
	return nil
}

// GetPopular 获取热门地点
func (r *LocationRepo) GetPopular() error {
	// TODO: 实现获取热门地点
	return nil
}

// UpdateMemoryCount 更新记忆数量
func (r *LocationRepo) UpdateMemoryCount() error {
	// TODO: 实现更新记忆数量
	return nil
}