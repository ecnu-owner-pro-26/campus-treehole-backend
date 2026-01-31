package model

// MemoryModel 记忆数据库模型
type MemoryModel struct {
	// TODO: 定义数据库表结构
	// ID, Title, Content, Latitude, Longitude
	// IsPublic, AllowComments, CreatorID
	// CreatedAt, UpdatedAt
}

// TableName 指定表名
func (MemoryModel) TableName() string {
	return "memories"
}
