package service

import (
	"campus-memory/application/dto"
	"campus-memory/infra/repo"
	"campus-memory/types/errno"
	"context"
	"gorm.io/gorm"
)

// LikeService 点赞服务
type LikeService struct {
	likeRepo    *repo.LikeRepo
	memoryRepo  *repo.MemoryRepo
	commentRepo *repo.CommentRepo
}

// NewLikeService 创建点赞服务
func NewLikeService(
	likeRepo *repo.LikeRepo,
	memoryRepo *repo.MemoryRepo,
	commentRepo *repo.CommentRepo,
) *LikeService {
	return &LikeService{
		likeRepo:    likeRepo,
		memoryRepo:  memoryRepo,
		commentRepo: commentRepo,
	}
}

// ToggleLike 切换点赞状态（统一方法，通过targetType区分记忆/评论）
func (s *LikeService) ToggleLike(userID, targetID int64, targetType int8) (*dto.ToggleLikeResponse, error) {
	// 1. 根据targetType检查目标是否存在
	if targetType == dto.LikeTargetTypeMemory {
		_, err := s.memoryRepo.GetByID(targetID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errno.ErrMemoryNotFound
			}
			return nil, err
		}
	} else if targetType == dto.LikeTargetTypeComment {
		ctx := context.Background()
		_, err := s.commentRepo.GetByID(ctx, targetID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errno.ErrCommentNotFound
			}
			return nil, err
		}
	} else {
		return nil, errno.ErrBadRequest
	}

	// 2. 检查是否已点赞
	isLiked, err := s.likeRepo.CheckLiked(userID, targetID, targetType)
	if err != nil {
		return nil, err
	}

	// 3. 翻转逻辑
	var delta int64
	if isLiked {
		// 已点赞 → 取消点赞
		if err := s.likeRepo.DeleteLike(userID, targetID, targetType); err != nil {
			return nil, errno.ErrLikeDeleteFail
		}
		delta = -1
	} else {
		// 未点赞 → 点赞
		if err := s.likeRepo.CreateLike(userID, targetID, targetType); err != nil {
			return nil, errno.ErrLikeCreateFail
		}
		delta = 1
	}

	// 4. 更新目标的点赞计数
	if targetType == dto.LikeTargetTypeMemory {
		_ = s.memoryRepo.UpdateCounts(targetID, &delta, nil)
	} else if targetType == dto.LikeTargetTypeComment {
		ctx := context.Background()
		if delta > 0 {
			_ = s.commentRepo.IncrementLikeCount(ctx, targetID)
		} else {
			_ = s.commentRepo.DecrementLikeCount(ctx, targetID)
		}
	}

	// 5. 获取最新点赞数
	likeCount, err := s.likeRepo.GetLikeCount(targetID, targetType)
	if err != nil {
		return nil, err
	}

	// 6. 返回响应
	return &dto.ToggleLikeResponse{
		IsLiked:   !isLiked,
		LikeCount: likeCount,
	}, nil
}
