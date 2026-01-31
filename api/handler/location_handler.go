package handler

// LocationHandler 处理地图位置相关的HTTP请求
type LocationHandler struct {
	// TODO: 注入application层依赖
}

// NewLocationHandler 创建位置处理器
func NewLocationHandler() *LocationHandler {
	// TODO: 初始化handler
	return &LocationHandler{}
}

// GetMemoriesByLocation 根据地图位置获取记忆
func (h *LocationHandler) GetMemoriesByLocation() {
	// TODO: 实现根据位置获取记忆的逻辑
}

// GetNearbyMemories 获取附近的记忆
func (h *LocationHandler) GetNearbyMemories() {
	// TODO: 实现获取附近记忆的逻辑
}
