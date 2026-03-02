package repo

import (
	"campus-memory/infra/model"
	"gorm.io/gorm"
)

// ImageRepo 图片数据访问层
type ImageRepo struct {
	db *gorm.DB
}

// NewImageRepo 创建图片仓储实例
func NewImageRepo(db *gorm.DB) *ImageRepo {
	return &ImageRepo{db: db}
}

// Create 创建图片记录
func (r *ImageRepo) Create(image *model.ImageModel) error {
	return r.db.Create(image).Error
}

// GetByID 根据ID获取图片
func (r *ImageRepo) GetByID(id int64) (*model.ImageModel, error) {
	var image model.ImageModel
	err := r.db.Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetByMemoryID 根据记忆ID获取图片列表
func (r *ImageRepo) GetByMemoryID(memoryID int64) ([]*model.ImageModel, error) {
	var images []*model.ImageModel
	err := r.db.Where("memory_id = ?", memoryID).
		Order("sort_order ASC").
		Find(&images).Error
	return images, err
}

// Delete 删除图片
func (r *ImageRepo) Delete(id int64) error {
	return r.db.Delete(&model.ImageModel{}, id).Error
}

// DeleteByMemoryID 删除记忆的所有图片
func (r *ImageRepo) DeleteByMemoryID(memoryID int64) error {
	return r.db.Where("memory_id = ?", memoryID).Delete(&model.ImageModel{}).Error
}
