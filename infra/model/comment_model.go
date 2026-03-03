package model

import (
	"time"
)

// CommentModel 留言数据库模型
type CommentModel struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MemoryID      int64      `gorm:"column:memory_id;type:integer;not null;index:idx_comment_memory" json:"memory_id"`
	UserID        int64      `gorm:"column:user_id;type:integer;not null;index:idx_comment_user" json:"user_id"`
	Content       string     `gorm:"column:content;type:text;not null" json:"content"`
	ParentID      *int64     `gorm:"column:parent_id;type:integer;index:idx_comment_parent" json:"parent_id,omitempty"` // NULL表示一级评论
	ReplyToUserID *int64     `gorm:"column:reply_to_user_id;type:integer" json:"reply_to_user_id,omitempty"`
	LikeCount     int64      `gorm:"column:like_count;type:integer;default:0" json:"like_count"`
	ReplyCount    int64      `gorm:"column:reply_count;type:integer;default:0" json:"reply_count"`
	Status        int8       `gorm:"column:status;type:integer;default:1;not null" json:"status"` // 0-待审核 1-已发布
	CreatedAt     time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (CommentModel) TableName() string {
	return "comments"
}
