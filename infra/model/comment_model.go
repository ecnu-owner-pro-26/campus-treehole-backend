package model

// CommentModel 留言数据库模型
type CommentModel struct {
	// TODO: 定义数据库表结构
	// ID, MemoryID, Content, UserID, CreatedAt
}

// TableName 指定表名
func (CommentModel) TableName() string {
	return "comments"
}
