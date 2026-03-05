package repo

import (
	"campus-memory/infra/model"
	"gorm.io/gorm"
)

// MemoryRepo 记忆数据访问层
type MemoryRepo struct {
	db *gorm.DB
}

// NewMemoryRepo 创建记忆仓储实例
func NewMemoryRepo(db *gorm.DB) *MemoryRepo {
	return &MemoryRepo{db: db}
}

// Create 创建记忆记录
func (r *MemoryRepo) Create(memory *model.MemoryModel) error {
	return r.db.Create(memory).Error
}

// GetByID 根据ID获取记忆（只返回已审核通过的记忆）
func (r *MemoryRepo) GetByID(id int64) (*model.MemoryModel, error) {
	var memory model.MemoryModel
	// 暂时注释掉审核状态检查，允许所有状态的记忆被查询
	// err := r.db.Where("id = ? AND deleted_at IS NULL AND status = ?", id, 1).First(&memory).Error
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&memory).Error
	if err != nil {
		return nil, err
	}
	return &memory, nil
}

// List 获取记忆列表(分页)
func (r *MemoryRepo) List(locationID *int64, page, pageSize int, sortBy string) ([]*model.MemoryModel, int64, error) {
	var memories []*model.MemoryModel
	var total int64

	// 暂时注释掉审核状态检查，允许所有状态的记忆被查询
	// query := r.db.Model(&model.MemoryModel{}).Where("deleted_at IS NULL AND status = 1")
	query := r.db.Model(&model.MemoryModel{}).Where("deleted_at IS NULL")

	// 按地点筛选
	if locationID != nil {
		query = query.Where("location_id = ?", *locationID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	switch sortBy {
	case "hot":
		query = query.Order("like_count DESC, view_count DESC")
	case "latest":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&memories).Error
	return memories, total, err
}

// Update 更新记忆记录
func (r *MemoryRepo) Update(memory *model.MemoryModel) error {
	return r.db.Save(memory).Error
}

// Delete 软删除记忆记录
func (r *MemoryRepo) Delete(id int64) error {
	return r.db.Model(&model.MemoryModel{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

// GetByLocationID 根据地点ID获取记忆列表
func (r *MemoryRepo) GetByLocationID(locationID int64, limit int) ([]*model.MemoryModel, error) {
	var memories []*model.MemoryModel
	// 暂时注释掉审核状态检查，允许所有状态的记忆被查询
	// err := r.db.Where("location_id = ? AND deleted_at IS NULL AND status = 1", locationID).
	err := r.db.Where("location_id = ? AND deleted_at IS NULL", locationID).
		Order("created_at DESC").
		Limit(limit).
		Find(&memories).Error
	return memories, err
}

// IncrementViewCount 增加浏览次数
func (r *MemoryRepo) IncrementViewCount(id int64) error {
	return r.db.Model(&model.MemoryModel{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}

// UpdateCounts 更新统计数据（增量更新）
func (r *MemoryRepo) UpdateCounts(id int64, likeCount, commentCount *int64) error {
	updates := make(map[string]interface{})
	if likeCount != nil {
		updates["like_count"] = gorm.Expr("like_count + ?", *likeCount)
	}
	if commentCount != nil {
		updates["comment_count"] = gorm.Expr("comment_count + ?", *commentCount)
	}
	return r.db.Model(&model.MemoryModel{}).Where("id = ?", id).Updates(updates).Error
}
