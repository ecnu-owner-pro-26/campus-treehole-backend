package repo

import (
	"campus-memory/infra/model"
	"context"

	"gorm.io/gorm"
)

// CampusRepo 校区数据访问层
type CampusRepo struct {
	db *gorm.DB
}

// NewCampusRepo 创建校区仓储实例
func NewCampusRepo(db *gorm.DB) *CampusRepo {
	return &CampusRepo{db: db}
}

// GetAll 获取所有校区
func (r *CampusRepo) GetAll(ctx context.Context) ([]model.CampusModel, error) {
	var campuses []model.CampusModel
	// SQLite 中 bool 存储为 0/1，使用 1 而不是 true
	err := r.db.WithContext(ctx).Where("is_active = ?", 1).
		Order("sort_order ASC, id ASC").
		Find(&campuses).Error
	return campuses, err
}

// GetByID 根据ID获取校区
func (r *CampusRepo) GetByID(ctx context.Context, id int64) (*model.CampusModel, error) {
	var campus model.CampusModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&campus).Error
	if err != nil {
		return nil, err
	}
	return &campus, nil
}

// Create 创建校区
func (r *CampusRepo) Create(ctx context.Context, campus *model.CampusModel) error {
	return r.db.WithContext(ctx).Create(campus).Error
}

// Update 更新校区
func (r *CampusRepo) Update(ctx context.Context, campus *model.CampusModel) error {
	return r.db.WithContext(ctx).Save(campus).Error
}

// Delete 删除校区
func (r *CampusRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.CampusModel{}, id).Error
}
