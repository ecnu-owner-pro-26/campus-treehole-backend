package service

import (
	"campus-memory/application/assembler"
	"campus-memory/application/dto"
	"campus-memory/infra/model"
	"campus-memory/infra/repo"
	"campus-memory/types/errno"
	"context"

	"gorm.io/gorm"
)

// CommentService 留言服务
type CommentService struct {
	commentRepo *repo.CommentRepo
	assembler   *assembler.CommentAssembler
	userRepo    *repo.UserRepo
	likeRepo    *repo.LikeRepo
	memoryRepo  *repo.MemoryRepo
}

func NewCommentService(
	commentRepo *repo.CommentRepo,
	userRepo *repo.UserRepo,
	likeRepo *repo.LikeRepo,
	assembler *assembler.CommentAssembler,
	memoryRepo *repo.MemoryRepo,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		userRepo:    userRepo,
		likeRepo:    likeRepo,
		assembler:   assembler,
		memoryRepo:  memoryRepo,
	}
}

// CreateComment 创建评论
func (s *CommentService) CreateComment(ctx context.Context, req *dto.CreateCommentRequest, userID int64) (*dto.CommentResponse, error) {
	// 验证记忆是否存在且公开
	memory, err := s.memoryRepo.GetByID(req.MemoryID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrMemoryNotFound
		}
		return nil, err
	}
	if !memory.IsPublic {
		return nil, errno.ErrMemoryNotPublic
	}

	// 验证父评论（如果是回复）
	var replyToUserID *int64
	if req.ParentID != nil && *req.ParentID > 0 {
		parentComment, err := s.commentRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errno.ErrParentCommentNotFound
			}
			return nil, err
		}

		// 验证父评论状态
		if parentComment == nil {
			return nil, errno.ErrParentCommentNotFound
		}

		// 设置回复目标用户ID
		replyToUserID = &parentComment.UserID
	}

	// 转换为 Model
	comment := s.assembler.ToDomain(req, userID)
	comment.ReplyToUserID = replyToUserID

	// 保存到数据库
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, errno.ErrCommentCreateFail
	}

	// 如果是回复，更新父评论的回复数
	if req.ParentID != nil && *req.ParentID > 0 {
		// 忽略更新失败，不影响主流程
		_ = s.commentRepo.IncrementReplyCount(ctx, *req.ParentID)
	}

	// 异步更新记忆的评论计数
	commentCount := int64(1)
	go func() {
		_ = s.memoryRepo.UpdateCounts(req.MemoryID, nil, &commentCount)
	}()

	// 获取评论者信息
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		// 如果用户不存在，返回适当的错误码；如果是其他错误，返回系统错误
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrUserNotFound
		}
		return nil, err
	}

	// 获取被回复用户信息
	var replyToUser *model.UserModel
	if replyToUserID != nil {
		replyToUser, _ = s.userRepo.GetUserByID(ctx, *replyToUserID)
	}

	// 获取回复数量（如果是顶级评论）
	var replyCount int64
	if req.ParentID == nil {
		replyCount, _ = s.commentRepo.GetReplyCount(ctx, comment.ID)
	}

	// 组装响应
	return s.assembler.ToCommentResponse(comment, user, replyToUser, false, replyCount), nil
}

// ListComments 获取评论列表（只查询顶级评论）
func (s *CommentService) ListComments(ctx context.Context, req *dto.CommentListRequest, currentUserID *int64) (*dto.CommentListResponse, error) {
	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.SortBy == "" {
		req.SortBy = "latest"
	}
	// 查询评论列表
	comments, total, err := s.commentRepo.ListByMemoryID(ctx, req.MemoryID, req.Page, req.PageSize, req.SortBy)
	if err != nil {
		return nil, err
	}
	// 组装响应列表
	commentItems := make([]*dto.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		// 获取评论用户
		user, _ := s.userRepo.GetUserByID(ctx, comment.UserID)
		if user == nil {
			continue
		}

		// 获取回复目标用户（如果是回复）
		var replyToUser *model.UserModel
		if comment.ReplyToUserID != nil {
			replyToUser, _ = s.userRepo.GetUserByID(ctx, *comment.ReplyToUserID)
		}

		// 检查是否已点赞
		isLiked := false
		if currentUserID != nil {
			isLiked, _ = s.likeRepo.CheckLiked(*currentUserID, comment.ID, 2)
		}

		// 获取回复数量
		replyCount, _ := s.commentRepo.GetReplyCount(ctx, comment.ID)

		commentItems = append(commentItems, s.assembler.ToCommentResponse(comment, user, replyToUser, isLiked, replyCount))
	}

	return &dto.CommentListResponse{
		Comments: commentItems,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// ListReplies 获取评论的回复列表
func (s *CommentService) ListReplies(ctx context.Context, parentID int64, page, pageSize int, currentUserID *int64) (*dto.CommentListResponse, error) {
	// 设置默认值
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	if pageSize > 50 {
		pageSize = 50
	}

	// 验证父评论是否存在
	parentComment, err := s.commentRepo.GetByID(ctx, parentID)
	if err != nil {
		return nil, err
	}
	if parentComment == nil {
		return nil, errno.ErrCommentNotFound
	}

	// 查询回复列表
	replies, total, err := s.commentRepo.ListRepliesByParentID(ctx, parentID, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 组装响应列表
	replyItems := make([]*dto.CommentResponse, 0, len(replies))
	for _, reply := range replies {
		// 获取评论用户
		user, err := s.userRepo.GetUserByID(ctx, reply.UserID)
		if err != nil || user == nil {
			continue
		}

		// 获取回复目标用户（如果是回复）
		var replyToUser *model.UserModel
		if reply.ReplyToUserID != nil {
			replyToUser, _ = s.userRepo.GetUserByID(ctx, *reply.ReplyToUserID)
		}

		// 检查是否已点赞
		isLiked := false
		if currentUserID != nil {
			isLiked, _ = s.likeRepo.CheckLiked(*currentUserID, reply.ID, 2)
		}

		// 获取回复数量（回复的回复）
		replyCount, _ := s.commentRepo.GetReplyCount(ctx, reply.ID)

		replyItems = append(replyItems, s.assembler.ToCommentResponse(reply, user, replyToUser, isLiked, replyCount))
	}

	return &dto.CommentListResponse{
		Comments: replyItems,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(ctx context.Context, id int64, userID int64) error {
	//获取评论
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查评论是否存在
	if comment == nil {
		return errno.ErrCommentNotFound
	}

	// 权限检查（只能删除自己的评论）
	if comment.UserID != userID {
		return errno.ErrForbidden
	}

	// 软删除
	if err := s.commentRepo.Delete(ctx, id); err != nil {
		return errno.ErrCommentDeleteFail
	}

	// 异步更新记忆的评论计数
	go func() {
		// 传递增量 -1
		delta := int64(-1)
		_ = s.memoryRepo.UpdateCounts(comment.MemoryID, nil, &delta)
	}()

	// 如果是回复，更新父评论的回复数
	if comment.ParentID != nil && *comment.ParentID > 0 {
		go func() {
			newCtx := context.Background()
			_ = s.commentRepo.DecrementReplyCount(newCtx, *comment.ParentID)
		}()
	}
	return nil
}
