package assembler

import (
	"campus-memory/application/dto"
	"campus-memory/infra/model"
	"time"
)

// CommentAssembler 评论数据组装器
type CommentAssembler struct{}

// NewCommentAssembler 创建评论组装器
func NewCommentAssembler() *CommentAssembler {
	return &CommentAssembler{}
}

// ToDomain 将DTO转换为评论Model
func (a *CommentAssembler) ToDomain(req *dto.CreateCommentRequest, userID int64) *model.CommentModel {
	if req == nil {
		return nil
	}

	now := time.Now()
	return &model.CommentModel{
		MemoryID:      req.MemoryID,
		UserID:        userID,
		Content:       req.Content,
		ParentID:      req.ParentID,
		ReplyToUserID: req.ReplyToUserID,
		LikeCount:     0,
		Status:        1, // 暂时设置为已发布状态，跳过审核
		CreatedAt:     now,
	}
}

// ToCommentResponse 将 Model 转换为 DTO Response
func (a *CommentAssembler) ToCommentResponse(
	comment *model.CommentModel,
	user *model.UserModel,
	replyToUser *model.UserModel,
	isLiked bool,
	replyCount int64,
) *dto.CommentResponse {
	// 组装响应
	return &dto.CommentResponse{
		ID:       comment.ID,
		MemoryID: comment.MemoryID,
		Content:  comment.Content,
		ParentID: comment.ParentID,
		Creator: dto.UserSimpleInfo{
			ID:       user.ID,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		},
		ReplyToUser: func() *dto.UserSimpleInfo {
			if replyToUser != nil {
				return &dto.UserSimpleInfo{
					ID:       replyToUser.ID,
					Nickname: replyToUser.Nickname,
					Avatar:   replyToUser.Avatar,
				}
			}
			return nil
		}(),
		LikeCount:  comment.LikeCount,
		ReplyCount: replyCount,
		IsLiked:    isLiked,
		CreatedAt:  comment.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
