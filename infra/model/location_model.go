package model

import "time"

// LocationModel 地点数据库模型
type LocationModel struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CampusID    int64     `gorm:"column:campus_id;type:integer;not null;index:idx_campus" json:"campus_id"` // 所属校区
	Name        string    `gorm:"column:name;type:text;not null" json:"name"`                               // 地点名称：图书馆、教学楼A
	Category    string    `gorm:"column:category;type:text" json:"category"`                                // 地点类型：教学、宿舍、餐厅、景点
	Latitude    float64   `gorm:"column:latitude;type:real;not null" json:"latitude"`                       // 纬度（必填，用于地图传送）
	Longitude   float64   `gorm:"column:longitude;type:real;not null" json:"longitude"`                     // 经度（必填，用于地图传送）
	IsActive    bool      `gorm:"column:is_active;type:integer;default:1" json:"is_active"`                 // 是否启用
	SortOrder   int       `gorm:"column:sort_order;type:integer;default:0" json:"sort_order"`               // 显示顺序
	MemoryCount int64     `gorm:"column:memory_count;type:integer;default:0" json:"memory_count"`           // 记忆数量（冗余字段）
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName 指定表名
func (LocationModel) TableName() string {
	return "locations"
}
