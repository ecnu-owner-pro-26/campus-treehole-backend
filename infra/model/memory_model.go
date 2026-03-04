package model

import (
	"time"
)

// MemoryModel 记忆数据库模型
type MemoryModel struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title        string     `gorm:"column:title;type:text;not null" json:"title"`
	Content      string     `gorm:"column:content;type:text" json:"content"`
	LocationName string     `gorm:"column:location_name;type:text" json:"location_name"`
	LocationID   *int64     `gorm:"column:location_id;type:integer;index:idx_location" json:"location_id"` // 关联location表（可选）
	Latitude     float64    `gorm:"column:latitude;type:real;not null" json:"latitude"`                    // 纬度（必填，用于地图显示）
	Longitude    float64    `gorm:"column:longitude;type:real;not null" json:"longitude"`                  // 经度（必填，用于地图显示）
	IsPublic     bool       `gorm:"column:is_public;type:integer;default:1;not null" json:"is_public"`
	ViewCount    int64      `gorm:"column:view_count;type:integer;default:0" json:"view_count"`
	LikeCount    int64      `gorm:"column:like_count;type:integer;default:0" json:"like_count"`
	CommentCount int64      `gorm:"column:comment_count;type:integer;default:0" json:"comment_count"`
	Tags         string     `gorm:"column:tags;type:text" json:"tags"`                           // JSON格式存储标签
	Status       int8       `gorm:"column:status;type:integer;default:1;not null" json:"status"` // 0-待审核 1-已发布 2-已下架
	CreatorID    int64      `gorm:"column:creator_id;type:integer;not null;index:idx_creator" json:"creator_id"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;index:idx_created" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;type:datetime;index:idx_deleted" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (MemoryModel) TableName() string {
	return "memories"
}
