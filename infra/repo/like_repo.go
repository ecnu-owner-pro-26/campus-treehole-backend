package repo

import (
	"campus-memory/infra/model"
	"gorm.io/gorm"
)

// LikeRepo 点赞数据访问层
type LikeRepo struct {
	db *gorm.DB
}

// NewLikeRepo 创建点赞仓库实例
func NewLikeRepo(db *gorm.DB) *LikeRepo {
	return &LikeRepo{db: db}
}

// Create 创建点赞记录
func (r *LikeRepo) Create(userID, targetID int64, targetType int8) error {
	like := &model.LikeModel{
		UserID:     userID,
		TargetID:   targetID,
		TargetType: targetType,
	}
	return r.db.Create(like).Error
}

// Delete 删除点赞记录
func (r *LikeRepo) Delete(userID, targetID int64, targetType int8) error {
	return r.db.Where("user_id = ? AND target_id = ? AND target_type = ?", 
		userID, targetID, targetType).Delete(&model.LikeModel{}).Error
}

// CheckExists 检查点赞是否存在
func (r *LikeRepo) CheckExists(userID, targetID int64, targetType int8) (bool, error) {
	var count int64
	err := r.db.Model(&model.LikeModel{}).
		Where("user_id = ? AND target_id = ? AND target_type = ?", 
			userID, targetID, targetType).Count(&count).Error
	return count > 0, err
}

// GetLikeCount 获取点赞数量
func (r *LikeRepo) GetLikeCount(targetID int64, targetType int8) (int64, error) {
	var count int64
	err := r.db.Model(&model.LikeModel{}).
		Where("target_id = ? AND target_type = ?", targetID, targetType).Count(&count).Error
	return count, err
}