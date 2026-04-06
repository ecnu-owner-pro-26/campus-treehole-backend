package repo

import (
	"campus-memory/infra/model"
	"context"

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
func (r *MemoryRepo) Create(ctx context.Context, memory *model.MemoryModel) error {
	return r.db.WithContext(ctx).Create(memory).Error
}

// GetByID 根据ID获取记忆
func (r *MemoryRepo) GetByID(ctx context.Context, id int64) (*model.MemoryModel, error) {
	var memory model.MemoryModel
	// 暂时注释掉审核状态检查，允许所有状态的记忆被查询
	// err := r.db.Where("id = ? AND deleted_at IS NULL AND status = ?", id, 1).First(&memory).Error
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&memory).Error
	if err != nil {
		return nil, err
	}
	return &memory, nil
}

// List 获取记忆列表(分页)
func (r *MemoryRepo) List(ctx context.Context, locationID *int64, tagsMask int64,
	minLat, maxLat, minLng, maxLng *float64, page, pageSize int, sortBy string) ([]*model.MemoryModel, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.MemoryModel{}).Where("deleted_at IS NULL").Debug()
	var memories []*model.MemoryModel
	var total int64

	// 暂时注释掉审核状态检查，允许所有状态的记忆被查询
	// query := r.db.Model(&model.MemoryModel{}).Where("deleted_at IS NULL AND status = 1")

	// 按地点筛选
	if locationID != nil && *locationID != 0 {
		query = query.Where("location_id = ?", *locationID)
	}

	// 标签筛选（位掩码）
	if tagsMask != 0 {
		// (tags_mask & tagsMask) = tagsMask 表示包含全部标签
		query = query.Where("(tags_mask & ?) = ?", tagsMask, tagsMask)
	}

	// 坐标范围筛选
	if minLat != nil && maxLat != nil && minLng != nil && maxLng != nil {
		// 只有当范围有效（min < max）且非全零时才添加条件
		if *minLat < *maxLat && *minLng < *maxLng && (*minLat != 0 || *maxLat != 0 || *minLng != 0 || *maxLng != 0) {
			query = query.Where("latitude BETWEEN ? AND ?", minLat, maxLat).
				Where("longitude BETWEEN ? AND ?", minLng, maxLng)
		}
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
func (r *MemoryRepo) Update(ctx context.Context, memory *model.MemoryModel) error {
	return r.db.WithContext(ctx).Save(memory).Error
}

// Delete 软删除记忆记录
func (r *MemoryRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.MemoryModel{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

// GetByLocationID 根据地点ID获取记忆列表
func (r *MemoryRepo) GetByLocationID(ctx context.Context, locationID int64, limit int) ([]*model.MemoryModel, error) {
	var memories []*model.MemoryModel
	// 暂时注释掉审核状态检查，允许所有状态的记忆被查询
	// err := r.db.Where("location_id = ? AND deleted_at IS NULL AND status = 1", locationID).
	err := r.db.WithContext(ctx).Where("location_id = ? AND deleted_at IS NULL", locationID).
		Order("created_at DESC").
		Limit(limit).
		Find(&memories).Error
	return memories, err
}

// IncrementViewCount 增加浏览次数
func (r *MemoryRepo) IncrementViewCount(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.MemoryModel{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}

// UpdateCounts 更新统计数据（增量更新）
func (r *MemoryRepo) UpdateCounts(ctx context.Context, id int64, likeCount, commentCount *int64) error {
	updates := make(map[string]interface{})
	if likeCount != nil {
		updates["like_count"] = gorm.Expr("like_count + ?", *likeCount)
	}
	if commentCount != nil {
		updates["comment_count"] = gorm.Expr("comment_count + ?", *commentCount)
	}
	return r.db.WithContext(ctx).Model(&model.MemoryModel{}).Where("id = ?", id).Updates(updates).Error
}

// Search 模糊搜索（LIKE）
func (r *MemoryRepo) Search(ctx context.Context, keyword string, tagsMask int64, tagMode string, page, pageSize int, sortBy string, currentUserID *int64) ([]*model.MemoryModel, int64, error) {
	// 基础查询
	query := r.db.WithContext(ctx).Table("memories m")

	// 权限过滤
	if currentUserID != nil {
		query = query.Where("m.is_public = ? OR m.creator_id = ?", true, *currentUserID)
	} else {
		query = query.Where("m.is_public = ?", true)
	}

	// 标签筛选（位掩码）
	if tagsMask != 0 {
		if tagMode == "and" {
			query = query.Where("(m.tags_mask & ?) = ?", tagsMask, tagsMask)
		} else {
			query = query.Where("(m.tags_mask & ?) != 0", tagsMask)
		}
	}

	// 关键词搜索：使用 LIKE 模糊匹配
	if keyword != "" {
		likePattern := "%" + keyword + "%"
		query = query.Where("m.title LIKE ? OR m.content LIKE ?", likePattern, likePattern)
	}

	// 排序：按创建时间倒序
	query = query.Order("m.created_at DESC")

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	var memories []*model.MemoryModel
	err := query.Select("m.*").Offset(offset).Limit(pageSize).Find(&memories).Error
	return memories, total, err
}
