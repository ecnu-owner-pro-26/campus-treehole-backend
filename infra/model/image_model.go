package model

// ImageModel 图片数据库模型
type ImageModel struct {
	// TODO: 定义数据库表结构
	// ID, MemoryID, URL, CreatedAt
}

// TableName 指定表名
func (ImageModel) TableName() string {
	return "images"
}
