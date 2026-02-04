package model

import (
	"time"
)

// UserModel 用户数据库模型
type UserModel struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username  string     `gorm:"column:username;type:text;not null;uniqueIndex:idx_username" json:"username"`
	Password  string     `gorm:"column:password;type:text;not null" json:"-"` // 密码不返回给前端
	Phone     string     `gorm:"column:phone;type:text;uniqueIndex:idx_phone" json:"phone"`
	Email     string     `gorm:"column:email;type:text;uniqueIndex:idx_email" json:"email"`
	Nickname  string     `gorm:"column:nickname;type:text" json:"nickname"`
	Avatar    string     `gorm:"column:avatar;type:text" json:"avatar"`
	Status    int8       `gorm:"column:status;type:integer;default:1;not null" json:"status"` // 0-禁用 1-正常
	Role      int8       `gorm:"column:role;type:integer;default:0;not null" json:"role"`     // 0-普通用户 1-管理
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "users"
}
