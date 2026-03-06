package repo

import (
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
func (r *LikeRepo) CheckLiked(userID, targetID int64, targetType int8) (bool, error) {
	var count int64
	err := r.db.Model(&model.LikeModel{}).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Count(&count).Error
	return count > 0, err
}

// CreateLike 创建点赞记录
func (r *LikeRepo) CreateLike(userID, targetID int64, targetType int8) error {
	like := &model.LikeModel{
		UserID:     userID,
		TargetID:   targetID,
		TargetType: targetType,
	}
	return r.db.Create(like).Error
}

// DeleteLike 删除点赞记录
func (r *LikeRepo) DeleteLike(userID, targetID int64, targetType int8) error {
	return r.db.Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Delete(&model.LikeModel{}).Error
}

// GetLikeCount 获取点赞数量
func (r *LikeRepo) GetLikeCount(targetID int64, targetType int8) (int64, error) {
	var count int64
	err := r.db.Model(&model.LikeModel{}).
		Where("target_id = ? AND target_type = ?", targetID, targetType).
		Count(&count).Error
	return count, err
}
