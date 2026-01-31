package types

// Memory 记忆数据类型
type Memory struct {
	// TODO: 定义记忆的数据结构
	// ID, 标题, 内容, 位置坐标, 创建时间, 更新时间
	// 是否公开, 是否允许留言, 创建者ID
	// 关联的图片列表, 音频列表
}

// MemoryImage 记忆关联的图片
type MemoryImage struct {
	// TODO: 定义图片的数据结构
	// ID, MemoryID, 图片URL, 创建时间
}

// MemoryAudio 记忆关联的音频
type MemoryAudio struct {
	// TODO: 定义音频的数据结构
	// ID, MemoryID, 音频URL, 时长, 创建时间
}
