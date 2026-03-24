package model

import (
	"time"
)

// UserModel 用户数据库模型（微信登录）
type UserModel struct {
	// 身份相关
	ID      int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"` // 用户唯一ID，主键，自动递增
	OpenID  string `gorm:"column:openid;type:text;not null;uniqueIndex:idx_openid" json:"-"`  // 微信OpenID，不返回给前端
	UnionID string `gorm:"column:unionid;type:text;uniqueIndex:idx_unionid" json:"-"`         // 微信UnionID，不返回给前端

	// 用户信息
	Nickname        string `gorm:"column:nickname;type:text;not null" json:"nickname"`                // 昵称
	Avatar          string `gorm:"column:avatar;type:text" json:"avatar"`                             // 头像URL
	BackgroundImage string `gorm:"column:background_image;type:text" json:"backgroundImage"`          // 背景图片URL
	Bio             string `gorm:"column:bio;type:text" json:"bio"`                                   // 个人简介
	DefaultCampusID *int64 `gorm:"column:default_campus_id;type:integer" json:"default_campus_id"`   // 默认校区ID（可为空）

	// 系统字段
	Status int8 `gorm:"column:status;type:integer;default:1;not null" json:"status"` // 用户状态：0=禁用 1=正常
	Role   int8 `gorm:"column:role;type:integer;default:0;not null" json:"role"`     // 用户角色：0=普通用户 1=管理员

	// 时间字段
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`    // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`    // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at,omitempty"`                    // 删除时间（软删除，可为空）
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "users"
}

