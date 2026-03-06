package model

import (
	"time"
)

// UserModel 用户数据库模型（微信登录）
type UserModel struct {
	ID              int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OpenID          string     `gorm:"column:openid;type:text;not null;uniqueIndex:idx_openid" json:"-"`
	UnionID         string     `gorm:"column:unionid;type:text;uniqueIndex:idx_unionid" json:"-"`
	Nickname        string     `gorm:"column:nickname;type:text;not null" json:"nickname"`
	Avatar          string     `gorm:"column:avatar;type:text" json:"avatar"`
	DefaultCampusID *int64     `gorm:"column:default_campus_id;type:integer" json:"default_campus_id"`
	Status          int8       `gorm:"column:status;type:integer;default:1;not null" json:"status"`
	Role            int8       `gorm:"column:role;type:integer;default:0;not null" json:"role"`
	CreatedAt       time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "users"
}

