package model

import "time"

// CampusModel 校区数据库模型
type CampusModel struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:text;not null" json:"name"`                 // 校区名称：普陀校区、临港校区、闵行校区
	IsActive  bool      `gorm:"column:is_active;type:integer;default:1" json:"is_active"`   // 是否启用
	SortOrder int       `gorm:"column:sort_order;type:integer;default:0" json:"sort_order"` // 显示顺序
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (CampusModel) TableName() string {
	return "campuses"
}
