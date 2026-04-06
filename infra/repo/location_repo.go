package repo

import (
	"campus-memory/infra/model"
	"context"

	"gorm.io/gorm"
)

// LocationRepo 地点数据访问层
type LocationRepo struct {
	db *gorm.DB
}

// NewLocationRepo 创建地点仓储实例
func NewLocationRepo(db *gorm.DB) *LocationRepo {
	return &LocationRepo{db: db}
}

// GetByID 根据ID获取地点
func (r *LocationRepo) GetByID(ctx context.Context, id int64) (*model.LocationModel, error) {
	var location model.LocationModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&location).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

// GetByCampusID 根据校区ID获取地点列表
func (r *LocationRepo) GetByCampusID(ctx context.Context, campusID int64) ([]model.LocationModel, error) {
	var locations []model.LocationModel
	err := r.db.WithContext(ctx).Where("campus_id = ? AND is_active = ?", campusID, true).
		Order("sort_order ASC, id ASC").
		Find(&locations).Error
	return locations, err
}

// List 获取地点列表(分页)
func (r *LocationRepo) List(ctx context.Context, campusID *int64, category string, page, pageSize int) ([]model.LocationModel, int64, error) {
	var locations []model.LocationModel
	var total int64

	query := r.db.WithContext(ctx).Model(&model.LocationModel{}).Where("is_active = ?", true)

	// 按校区筛选
	if campusID != nil {
		query = query.Where("campus_id = ?", *campusID)
	}

	// 按类别筛选
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("sort_order ASC, id ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&locations).Error
	return locations, total, err
}

// GetByCategory 根据类别获取地点
func (r *LocationRepo) GetByCategory(ctx context.Context, category string) ([]model.LocationModel, error) {
	var locations []model.LocationModel
	err := r.db.WithContext(ctx).Where("category = ? AND is_active = ?", category, true).
		Order("sort_order ASC, id ASC").
		Find(&locations).Error
	return locations, err
}

// Search 搜索地点
func (r *LocationRepo) Search(ctx context.Context, keyword string, campusID *int64) ([]model.LocationModel, error) {
	var locations []model.LocationModel
	query := r.db.WithContext(ctx).Where("is_active = ? AND name LIKE ?", true, "%"+keyword+"%")

	if campusID != nil {
		query = query.Where("campus_id = ?", *campusID)
	}

	err := query.Order("memory_count DESC, sort_order ASC").
		Limit(20).
		Find(&locations).Error
	return locations, err
}

// GetPopular 获取热门地点
func (r *LocationRepo) GetPopular(ctx context.Context, limit int) ([]model.LocationModel, error) {
	var locations []model.LocationModel
	err := r.db.WithContext(ctx).Where("is_active = ?", true).
		Order("memory_count DESC").
		Limit(limit).
		Find(&locations).Error
	return locations, err
}

// Create 创建地点
func (r *LocationRepo) Create(ctx context.Context, location *model.LocationModel) error {
	return r.db.WithContext(ctx).Create(location).Error
}

// Update 更新地点
func (r *LocationRepo) Update(ctx context.Context, location *model.LocationModel) error {
	return r.db.WithContext(ctx).Save(location).Error
}

// Delete 删除地点
func (r *LocationRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.LocationModel{}, id).Error
}

// UpdateMemoryCount 更新记忆数量
func (r *LocationRepo) UpdateMemoryCount(ctx context.Context, id int64, delta int64) error {
	return r.db.WithContext(ctx).Model(&model.LocationModel{}).Where("id = ?", id).
		Update("memory_count", gorm.Expr("memory_count + ?", delta)).Error
}
