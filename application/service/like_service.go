package service

import (
	"campus-memory/infra/repo"
	"campus-memory/types/errno"
)

// LikeService 点赞服务
type LikeService struct {
	likeRepo *repo.LikeRepo
}

// NewLikeService 创建点赞服务实例
func NewLikeService(likeRepo *repo.LikeRepo) *LikeService {
	return &LikeService{
		likeRepo: likeRepo,
	}
}

// ToggleLike 切换点赞状态（翻转点赞）
func (s *LikeService) ToggleLike(userID, targetID int64, targetType int8) (bool, error) {
	// 检查是否已点赞
	exists, err := s.likeRepo.CheckExists(userID, targetID, targetType)
	if err != nil {
		return false, err
	}

	if exists {
		// 已点赞，取消点赞
		err = s.likeRepo.Delete(userID, targetID, targetType)
		if err != nil {
			return false, err
		}
		return false, nil // 返回false表示取消点赞
	} else {
		// 未点赞，添加点赞
		err = s.likeRepo.Create(userID, targetID, targetType)
		if err != nil {
			return false, err
		}
		return true, nil // 返回true表示点赞成功
	}
}

// GetLikeStatus 获取点赞状态和数量
func (s *LikeService) GetLikeStatus(userID, targetID int64, targetType int8) (bool, int64, error) {
	// 检查用户是否已点赞
	isLiked, err := s.likeRepo.CheckExists(userID, targetID, targetType)
	if err != nil {
		return false, 0, err
	}

	// 获取总点赞数
	count, err := s.likeRepo.GetLikeCount(targetID, targetType)
	if err != nil {
		return false, 0, err
	}

	return isLiked, count, nil
}

// LikeMemory 点赞记忆（兼容旧接口）
func (s *LikeService) LikeMemory(userID, memoryID int64) error {
	// 检查是否已点赞
	exists, err := s.likeRepo.CheckExists(userID, memoryID, 1)
	if err != nil {
		return err
	}
	if exists {
		return errno.New(4001, "已经点赞过了")
	}

	return s.likeRepo.Create(userID, memoryID, 1)
}

// UnlikeMemory 取消点赞记忆（兼容旧接口）
func (s *LikeService) UnlikeMemory(userID, memoryID int64) error {
	// 检查是否已点赞
	exists, err := s.likeRepo.CheckExists(userID, memoryID, 1)
	if err != nil {
		return err
	}
	if !exists {
		return errno.New(4002, "还没有点赞")
	}

	return s.likeRepo.Delete(userID, memoryID, 1)
}

// LikeComment 点赞留言（兼容旧接口）
func (s *LikeService) LikeComment(userID, commentID int64) error {
	// 检查是否已点赞
	exists, err := s.likeRepo.CheckExists(userID, commentID, 2)
	if err != nil {
		return err
	}
	if exists {
		return errno.New(4001, "已经点赞过了")
	}

	return s.likeRepo.Create(userID, commentID, 2)
}

// UnlikeComment 取消点赞留言（兼容旧接口）
func (s *LikeService) UnlikeComment(userID, commentID int64) error {
	// 检查是否已点赞
	exists, err := s.likeRepo.CheckExists(userID, commentID, 2)
	if err != nil {
		return err
	}
	if !exists {
		return errno.New(4002, "还没有点赞")
	}

	return s.likeRepo.Delete(userID, commentID, 2)
}