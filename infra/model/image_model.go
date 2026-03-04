package model

import (
	"time"
)

// ImageModel 图片数据库模型（简化版）
type ImageModel struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MemoryID  int64     `gorm:"column:memory_id;type:integer;not null;index:idx_image_memory" json:"memory_id"`
	URL       string    `gorm:"column:url;type:text;not null" json:"url"`
	Size      int64     `gorm:"column:size;type:integer" json:"size"`      // 文件大小（字节）
	SortOrder int       `gorm:"column:sort_order;type:integer;default:0" json:"sort_order"` // 显示顺序
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName 指定表名
func (ImageModel) TableName() string {
	return "images"
}
