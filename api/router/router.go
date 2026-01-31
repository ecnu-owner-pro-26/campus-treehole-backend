package router

// Router 路由管理
type Router struct {
	// TODO: 注入handler依赖
}

// NewRouter 创建路由实例
func NewRouter() *Router {
	// TODO: 初始化路由
	return &Router{}
}

// SetupRoutes 设置路由
func (r *Router) SetupRoutes() {
	// TODO: 定义所有API路由
	// POST   /api/memories          - 创建记忆
	// GET    /api/memories/:id      - 获取记忆详情
	// GET    /api/memories          - 获取记忆列表
	// PUT    /api/memories/:id      - 更新记忆
	// DELETE /api/memories/:id      - 删除记忆
	// GET    /api/locations/:lat/:lng/memories - 根据位置获取记忆
	// POST   /api/memories/:id/comments - 创建留言
	// GET    /api/memories/:id/comments - 获取留言列表
	// DELETE /api/comments/:id      - 删除留言
}
