package model

// AudioModel 音频数据库模型
type AudioModel struct {
	// TODO: 定义数据库表结构
	// ID, MemoryID, URL, Duration, CreatedAt
}

// TableName 指定表名
func (AudioModel) TableName() string {
	return "audios"
}
