package model

import (
	"time"
)

// LikeModel 点赞数据库模型
type LikeModel struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"column:user_id;type:integer;not null;index:idx_like_user" json:"user_id"`
	TargetID   int64     `gorm:"column:target_id;type:integer;not null;index:idx_like_target" json:"target_id"`
	TargetType int8      `gorm:"column:target_type;type:integer;not null" json:"target_type"` // 1-记忆 2-留言
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName 指定表名
func (LikeModel) TableName() string {
	return "likes"
}
