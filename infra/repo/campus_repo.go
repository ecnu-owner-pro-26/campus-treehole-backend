package repo

import (
	"campus-memory/infra/model"
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
func (r *CampusRepo) GetAll() ([]model.CampusModel, error) {
	var campuses []model.CampusModel
	err := r.db.Where("is_active = ?", true).
		Order("sort_order ASC, id ASC").
		Find(&campuses).Error
	return campuses, err
}

// GetByID 根据ID获取校区
func (r *CampusRepo) GetByID(id int64) (*model.CampusModel, error) {
	var campus model.CampusModel
	err := r.db.Where("id = ?", id).First(&campus).Error
	if err != nil {
		return nil, err
	}
	return &campus, nil
}

// Create 创建校区
func (r *CampusRepo) Create(campus *model.CampusModel) error {
	return r.db.Create(campus).Error
}

// Update 更新校区
func (r *CampusRepo) Update(campus *model.CampusModel) error {
	return r.db.Save(campus).Error
}

// Delete 删除校区
func (r *CampusRepo) Delete(id int64) error {
	return r.db.Delete(&model.CampusModel{}, id).Error
}
