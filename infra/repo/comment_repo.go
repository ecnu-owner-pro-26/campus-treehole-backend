package repo

import (
	"campus-memory/types/errno"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"campus-memory/infra/model"
)

// CommentRepo 评论数据访问层
type CommentRepo struct {
	db *gorm.DB
}

// NewCommentRepo 创建评论仓库实例
func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

// Create 创建评论记录
func (r *CommentRepo) Create(ctx context.Context, comment *model.CommentModel) error {
	if comment == nil {
		return errors.New("comment is nil")
	}

	return r.db.WithContext(ctx).Create(comment).Error
}

// ListByMemoryID 根据记忆ID获取评论列表(只查询顶级评论）
func (r *CommentRepo) ListByMemoryID(ctx context.Context, memoryID int64, page, pageSize int, sortBy string) ([]*model.CommentModel, int64, error) {
	var comments []*model.CommentModel
	var total int64

	// 构建查询
	// 暂时注释掉审核状态检查，允许所有状态的评论被查询
	// query := r.db.WithContext(ctx).Model(&model.CommentModel{}).
	// 	Where("memory_id = ? AND status = ? AND deleted_at IS NULL", memoryID, 1).
	// 	Where("parent_id IS NULL") // 只查询顶级评论
	query := r.db.WithContext(ctx).Model(&model.CommentModel{}).
		Where("memory_id = ? AND deleted_at IS NULL", memoryID).
		Where("parent_id IS NULL") // 只查询顶级评论

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	switch sortBy {
	case "hot":
		query = query.Order("like_count DESC, created_at DESC")
	default: // "latest"
		query = query.Order("created_at DESC")
	}

	// 分页
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// ListRepliesByParentID 根据父评论ID获取回复列表
func (r *CommentRepo) ListRepliesByParentID(ctx context.Context, parentID int64, page, pageSize int) ([]*model.CommentModel, int64, error) {
	var replies []*model.CommentModel
	var total int64

	// 暂时注释掉审核状态检查，允许所有状态的评论被查询
	// query := r.db.WithContext(ctx).Model(&model.CommentModel{}).
	// 	Where("parent_id = ? AND status = ? AND deleted_at IS NULL", parentID, 1)
	query := r.db.WithContext(ctx).Model(&model.CommentModel{}).
		Where("parent_id = ? AND deleted_at IS NULL", parentID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.CommentModel{}, 0, nil
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at ASC").Find(&replies).Error; err != nil {
		return nil, 0, err
	}

	return replies, total, nil
}

// Delete 删除评论记录（软删除）
func (r *CommentRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.CommentModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
		}).Error
}

// GetByID 根据ID获取评论
func (r *CommentRepo) GetByID(ctx context.Context, id int64) (*model.CommentModel, error) {
	var comment model.CommentModel
	// 暂时注释掉审核状态检查，允许所有状态的评论被查询
	// err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL AND status = ?", id, 1).First(&comment).Error
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&comment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

// GetReplyCount 获取评论的回复数量
func (r *CommentRepo) GetReplyCount(ctx context.Context, commentID int64) (int64, error) {
	var count int64
	// 暂时注释掉审核状态检查，允许所有状态的评论被统计
	// err := r.db.WithContext(ctx).Model(&model.CommentModel{}).
	// 	Where("parent_id = ? AND status = ?", commentID, 1).
	// 	Count(&count).Error
	err := r.db.WithContext(ctx).Model(&model.CommentModel{}).
		Where("parent_id = ?", commentID).
		Count(&count).Error
	return count, err
}

// BatchGetReplyCount 批量获取多个评论的回复数量
func (r *CommentRepo) BatchGetReplyCount(ctx context.Context, commentIDs []int64) (map[int64]int64, error) {
	type Result struct {
		ParentID int64
		Count    int64
	}

	var results []Result
	// 暂时注释掉审核状态检查，允许所有状态的评论被统计
	// err := r.db.WithContext(ctx).Model(&model.CommentModel{}).
	// 	Select("parent_id, COUNT(*) as count").
	// 	Where("parent_id IN (?) AND status = ?", commentIDs, 1).
	// 	Group("parent_id").
	// 	Scan(&results).Error
	err := r.db.WithContext(ctx).Model(&model.CommentModel{}).
		Select("parent_id, COUNT(*) as count").
		Where("parent_id IN (?)", commentIDs).
		Group("parent_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	countMap := make(map[int64]int64)
	for _, r := range results {
		countMap[r.ParentID] = r.Count
	}
	return countMap, nil
}

// IncrementReplyCount 增加评论回复数
func (r *CommentRepo) IncrementReplyCount(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&model.CommentModel{}).
		Where("id = ?", id).
		UpdateColumn("reply_count", gorm.Expr("reply_count + ?", 1))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		// 没有找到对应记录，返回自定义错误
		return errno.ErrCommentNotFound
	}
	return nil
}

// DecrementReplyCount 减少评论回复数
func (r *CommentRepo) DecrementReplyCount(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&model.CommentModel{}).
		Where("id = ? AND reply_count > ?", id, 0).
		UpdateColumn("reply_count", gorm.Expr("reply_count - ?", 1))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		// 可能记录不存在或 reply_count 已为0
		return errno.ErrCommentNotFound
	}
	return nil
}

// IncrementLikeCount 增加评论点赞数
func (r *CommentRepo) IncrementLikeCount(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&model.CommentModel{}).
		Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errno.ErrCommentNotFound
	}
	return nil
}

// DecrementLikeCount 减少评论点赞数
func (r *CommentRepo) DecrementLikeCount(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Model(&model.CommentModel{}).
		Where("id = ? AND like_count > ?", id, 0).
		UpdateColumn("like_count", gorm.Expr("like_count - ?", 1))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errno.ErrCommentNotFound
	}
	return nil
}
