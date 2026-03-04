package repo

import (
	"context"

	"campus-memory/infra/model"

	"gorm.io/gorm"
)

// LikeRepo 点赞数据访问层
type LikeRepo struct {
	db *gorm.DB
}

// NewLikeRepo 创建点赞数据访问层
func NewLikeRepo(db *gorm.DB) *LikeRepo {
	return &LikeRepo{db: db}
}

// CheckLiked 检查用户是否已点赞
func (r *LikeRepo) CheckLiked(ctx context.Context, userID, targetID int64, targetType int8) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.LikeModel{}).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Count(&count).Error
	return count > 0, err
}

// CreateLike 创建点赞记录
func (r *LikeRepo) CreateLike(ctx context.Context, userID, targetID int64, targetType int8) error {
	like := &model.LikeModel{
		UserID:     userID,
		TargetID:   targetID,
		TargetType: targetType,
	}
	return r.db.WithContext(ctx).Create(like).Error
}

// DeleteLike 删除点赞记录
func (r *LikeRepo) DeleteLike(ctx context.Context, userID, targetID int64, targetType int8) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Delete(&model.LikeModel{}).Error
}

// GetLikeCount 获取点赞数量
func (r *LikeRepo) GetLikeCount(ctx context.Context, targetID int64, targetType int8) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.LikeModel{}).
		Where("target_id = ? AND target_type = ?", targetID, targetType).
		Count(&count).Error
	return count, err
}

// BatchCheckLiked 批量查询是否点赞
func (r *LikeRepo) BatchCheckLiked(ctx context.Context, userID int64, targetIDs []int64, targetType int) (map[int64]bool, error) {
	if len(targetIDs) == 0 {
		return make(map[int64]bool), nil
	}
	var likes []*model.LikeModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND target_id IN ? AND target_type = ?", userID, targetIDs, targetType).
		Find(&likes).Error
	if err != nil {
		return nil, err
	}
	likedMap := make(map[int64]bool, len(likes))
	for _, like := range likes {
		likedMap[like.TargetID] = true
	}
	return likedMap, nil
}
