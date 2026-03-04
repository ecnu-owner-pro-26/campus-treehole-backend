package service

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"campus-memory/application/assembler"
	"campus-memory/application/dto"
	"campus-memory/infra/repo"
	"campus-memory/types/errno"
)

// TestCommentService_CreateComment 测试创建评论
func TestCommentService_CreateComment(t *testing.T) {
	db, mock := setupMockDB(t)
	commentRepo := repo.NewCommentRepo(db)
	userRepo := repo.NewUserRepo(db)
	likeRepo := repo.NewLikeRepo(db)
	memoryRepo := repo.NewMemoryRepo(db)
	assembler := assembler.NewCommentAssembler()
	service := NewCommentService(commentRepo, userRepo, likeRepo, assembler, memoryRepo)

	ctx := context.Background()
	userID := int64(100)
	memoryID := int64(1)

	t.Run("成功创建顶级评论", func(t *testing.T) {
		req := &dto.CreateCommentRequest{
			MemoryID: memoryID,
			Content:  "这是一条评论",
		}

		// 1. 查询记忆是否存在且公开
		memoryRows := sqlmock.NewRows([]string{"id", "is_public"}).AddRow(memoryID, true)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(memoryID, 1).
			WillReturnRows(memoryRows)

		// 2. 创建评论（INSERT）
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `comments`").
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // deleted_at
				memoryID,
				userID,
				req.Content,
				nil, // ParentID
				nil, // ReplyToUserID
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// 3. 查询用户信息（用于组装响应）
		userRows := sqlmock.NewRows([]string{"id", "nickname", "avatar"}).
			AddRow(userID, "测试用户", "avatar.jpg")
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs(userID, 1).
			WillReturnRows(userRows)

		// 异步更新记忆评论数（不验证，因为 goroutine 可能延迟执行）
		// 如果希望验证，可以添加期望，但需要确保 goroutine 在测试结束前执行

		resp, err := service.CreateComment(ctx, req, userID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int64(1), resp.ID)
		assert.Equal(t, req.Content, resp.Content)
		assert.Equal(t, userID, resp.User.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功创建回复", func(t *testing.T) {
		parentID := int64(10)
		parentUserID := int64(200)
		req := &dto.CreateCommentRequest{
			MemoryID: memoryID,
			Content:  "回复内容",
			ParentID: &parentID,
		}

		// 1. 查询记忆
		memoryRows := sqlmock.NewRows([]string{"id", "is_public"}).AddRow(memoryID, true)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(memoryID, 1).
			WillReturnRows(memoryRows)

		// 2. 查询父评论
		parentRows := sqlmock.NewRows([]string{"id", "user_id"}).AddRow(parentID, parentUserID)
		mock.ExpectQuery("SELECT .* FROM `comments` WHERE id = ? AND `comments`.`deleted_at` IS NULL ORDER BY `comments`.`id` LIMIT ?").
			WithArgs(parentID, 1).
			WillReturnRows(parentRows)

		// 3. 创建回复
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `comments`").
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				memoryID,
				userID,
				req.Content,
				parentID,
				parentUserID, // ReplyToUserID
			).WillReturnResult(sqlmock.NewResult(2, 1))
		mock.ExpectCommit()

		// 4. 更新父评论回复数（异步，但可以添加期望，因为主流程中调用了 IncrementReplyCount，是同步的？实际上是同步调用，不是 goroutine）
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `comments` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), parentID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// 5. 查询用户信息
		userRows := sqlmock.NewRows([]string{"id", "nickname", "avatar"}).
			AddRow(userID, "测试用户", "avatar.jpg")
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs(userID, 1).
			WillReturnRows(userRows)

		// 6. 查询被回复用户（parentUser）
		parentUserRows := sqlmock.NewRows([]string{"id", "nickname", "avatar"}).
			AddRow(parentUserID, "父用户", "parent.jpg")
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs(parentUserID, 1).
			WillReturnRows(parentUserRows)

		resp, err := service.CreateComment(ctx, req, userID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int64(2), resp.ID)
		assert.Equal(t, req.Content, resp.Content)
		assert.Equal(t, userID, resp.User.ID)
		assert.NotNil(t, resp.ReplyToUser)
		assert.Equal(t, parentUserID, resp.ReplyToUser.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("记忆不存在", func(t *testing.T) {
		req := &dto.CreateCommentRequest{MemoryID: 999, Content: "test"}
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		resp, err := service.CreateComment(ctx, req, userID)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrMemoryNotFound, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("记忆不公开", func(t *testing.T) {
		req := &dto.CreateCommentRequest{MemoryID: memoryID, Content: "test"}
		memoryRows := sqlmock.NewRows([]string{"id", "is_public"}).AddRow(memoryID, false)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(memoryID, 1).
			WillReturnRows(memoryRows)

		resp, err := service.CreateComment(ctx, req, userID)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrMemoryNotPublic, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("父评论不存在", func(t *testing.T) {
		parentID := int64(999)
		req := &dto.CreateCommentRequest{MemoryID: memoryID, Content: "test", ParentID: &parentID}
		memoryRows := sqlmock.NewRows([]string{"id", "is_public"}).AddRow(memoryID, true)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(memoryID, 1).
			WillReturnRows(memoryRows)
		mock.ExpectQuery("SELECT .* FROM `comments` WHERE id = ? AND `comments`.`deleted_at` IS NULL ORDER BY `comments`.`id` LIMIT ?").
			WithArgs(parentID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		resp, err := service.CreateComment(ctx, req, userID)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrParentCommentNotFound, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCommentService_ListComments 测试获取评论列表
func TestCommentService_ListComments(t *testing.T) {
	db, mock := setupMockDB(t)
	commentRepo := repo.NewCommentRepo(db)
	userRepo := repo.NewUserRepo(db)
	likeRepo := repo.NewLikeRepo(db)
	memoryRepo := repo.NewMemoryRepo(db)
	assembler := assembler.NewCommentAssembler()
	service := NewCommentService(commentRepo, userRepo, likeRepo, assembler, memoryRepo)

	ctx := context.Background()
	memoryID := int64(1)
	currentUserID := int64(100)

	t.Run("成功获取评论列表（无点赞状态）", func(t *testing.T) {
		req := &dto.CommentListRequest{
			MemoryID: memoryID,
			Page:     1,
			PageSize: 10,
			SortBy:   "latest",
		}

		// 1. 查询评论列表
		commentRows := sqlmock.NewRows([]string{"id", "memory_id", "user_id", "content", "parent_id", "reply_to_user_id", "created_at"}).
			AddRow(1, memoryID, 101, "评论1", nil, nil, time.Now()).
			AddRow(2, memoryID, 102, "评论2", nil, nil, time.Now())
		mock.ExpectQuery("SELECT .* FROM `comments` WHERE memory_id = ? AND parent_id IS NULL AND `comments`.`deleted_at` IS NULL ORDER BY .* LIMIT 10 OFFSET 0").
			WithArgs(memoryID).
			WillReturnRows(commentRows)

		// 2. 查询总数
		mock.ExpectQuery("SELECT count\\(.*\\) FROM `comments` WHERE memory_id = \\? AND parent_id IS NULL AND `comments`.`deleted_at` IS NULL").
			WithArgs(memoryID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		// 3. 批量查询用户（userID 101, 102）
		userRows := sqlmock.NewRows([]string{"id", "nickname", "avatar"}).
			AddRow(101, "用户1", "a.jpg").
			AddRow(102, "用户2", "b.jpg")
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id IN \\(\\?,\\?\\) AND `users`.`deleted_at` IS NULL").
			WithArgs(101, 102).
			WillReturnRows(userRows)

		// 4. 批量查询点赞状态（由于 currentUserID 为 nil，不会调用 likeRepo）

		// 5. 批量查询回复数
		mock.ExpectQuery("SELECT comment_id, count\\(.*\\) FROM `comments` WHERE parent_id IN \\(\\?,\\?\\) AND `comments`.`deleted_at` IS NULL GROUP BY parent_id").
			WithArgs(1, 2).
			WillReturnRows(sqlmock.NewRows([]string{"parent_id", "count"}).
				AddRow(1, 3).
				AddRow(2, 0))

		resp, err := service.ListComments(ctx, req, nil)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 2, len(resp.Comments))
		assert.Equal(t, int64(2), resp.Total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功获取评论列表（带点赞状态）", func(t *testing.T) {
		req := &dto.CommentListRequest{
			MemoryID: memoryID,
			Page:     1,
			PageSize: 10,
			SortBy:   "latest",
		}

		commentRows := sqlmock.NewRows([]string{"id", "memory_id", "user_id", "content", "parent_id", "reply_to_user_id", "created_at"}).
			AddRow(1, memoryID, 101, "评论1", nil, nil, time.Now())
		mock.ExpectQuery("SELECT .* FROM `comments` WHERE memory_id = ? AND parent_id IS NULL AND `comments`.`deleted_at` IS NULL ORDER BY .* LIMIT 10 OFFSET 0").
			WithArgs(memoryID).
			WillReturnRows(commentRows)

		mock.ExpectQuery("SELECT count\\(.*\\) FROM `comments` WHERE memory_id = \\? AND parent_id IS NULL AND `comments`.`deleted_at` IS NULL").
			WithArgs(memoryID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		userRows := sqlmock.NewRows([]string{"id", "nickname", "avatar"}).
			AddRow(101, "用户1", "a.jpg")
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id IN \\(\\?\\) AND `users`.`deleted_at` IS NULL").
			WithArgs(101).
			WillReturnRows(userRows)

		// 批量查询点赞状态
		mock.ExpectQuery("SELECT .* FROM `likes` WHERE user_id = \\? AND target_id IN \\(\\?\\) AND target_type = \\?").
			WithArgs(currentUserID, 1, 2).
			WillReturnRows(sqlmock.NewRows([]string{"target_id"}).AddRow(1)) // 假设评论 1 被点赞

		mock.ExpectQuery("SELECT comment_id, count\\(.*\\) FROM `comments` WHERE parent_id IN \\(\\?\\) AND `comments`.`deleted_at` IS NULL GROUP BY parent_id").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"parent_id", "count"}).AddRow(1, 2))

		resp, err := service.ListComments(ctx, req, &currentUserID)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Comments))
		assert.True(t, resp.Comments[0].IsLiked)
		assert.Equal(t, int64(2), resp.Comments[0].ReplyCount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCommentService_ListReplies 测试获取回复列表
func TestCommentService_ListReplies(t *testing.T) {
	db, mock := setupMockDB(t)
	commentRepo := repo.NewCommentRepo(db)
	userRepo := repo.NewUserRepo(db)
	likeRepo := repo.NewLikeRepo(db)
	memoryRepo := repo.NewMemoryRepo(db)
	assembler := assembler.NewCommentAssembler()
	service := NewCommentService(commentRepo, userRepo, likeRepo, assembler, memoryRepo)

	ctx := context.Background()
	parentID := int64(10)
	currentUserID := int64(100)

	t.Run("成功获取回复列表", func(t *testing.T) {
		// 1. 验证父评论存在
		parentRows := sqlmock.NewRows([]string{"id"}).AddRow(parentID)
		mock.ExpectQuery("SELECT .* FROM `comments` WHERE id = ? AND `comments`.`deleted_at` IS NULL ORDER BY `comments`.`id` LIMIT ?").
			WithArgs(parentID, 1).
			WillReturnRows(parentRows)

		// 2. 查询回复列表
		replyRows := sqlmock.NewRows([]string{"id", "user_id", "content", "reply_to_user_id", "created_at"}).
			AddRow(21, 101, "回复1", int64(102), time.Now()).
			AddRow(22, 103, "回复2", nil, time.Now())
		mock.ExpectQuery("SELECT .* FROM `comments` WHERE parent_id = \\? AND `comments`.`deleted_at` IS NULL ORDER BY created_at ASC LIMIT 20 OFFSET 0").
			WithArgs(parentID).
			WillReturnRows(replyRows)

		// 3. 查询总数
		mock.ExpectQuery("SELECT count\\(.*\\) FROM `comments` WHERE parent_id = \\? AND `comments`.`deleted_at` IS NULL").
			WithArgs(parentID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		// 4. 批量查询用户（101, 103, 102）
		userRows := sqlmock.NewRows([]string{"id", "nickname", "avatar"}).
			AddRow(101, "用户A", "a.jpg").
			AddRow(102, "用户B", "b.jpg").
			AddRow(103, "用户C", "c.jpg")
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id IN \\(\\?,\\?,\\?\\) AND `users`.`deleted_at` IS NULL").
			WithArgs(101, 103, 102).
			WillReturnRows(userRows)

		// 5. 批量查询点赞状态
		mock.ExpectQuery("SELECT .* FROM `likes` WHERE user_id = \\? AND target_id IN \\(\\?,\\?\\) AND target_type = \\?").
			WithArgs(currentUserID, 21, 22, 2).
			WillReturnRows(sqlmock.NewRows([]string{"target_id"}).AddRow(21)) // 回复21被点赞

		// 6. 批量查询回复数（回复的回复）
		mock.ExpectQuery("SELECT comment_id, count\\(.*\\) FROM `comments` WHERE parent_id IN \\(\\?,\\?\\) AND `comments`.`deleted_at` IS NULL GROUP BY parent_id").
			WithArgs(21, 22).
			WillReturnRows(sqlmock.NewRows([]string{"parent_id", "count"}).
				AddRow(21, 1).
				AddRow(22, 0))

		resp, err := service.ListReplies(ctx, parentID, 1, 20, &currentUserID)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(resp.Comments))
		assert.Equal(t, int64(2), resp.Total)
		assert.True(t, resp.Comments[0].IsLiked)
		assert.Equal(t, int64(1), resp.Comments[0].ReplyCount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("父评论不存在", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `comments` WHERE id = ? AND `comments`.`deleted_at` IS NULL ORDER BY `comments`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		resp, err := service.ListReplies(ctx, 999, 1, 20, nil)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCommentService_DeleteComment 测试删除评论
func TestCommentService_DeleteComment(t *testing.T) {
	db, mock := setupMockDB(t)
	commentRepo := repo.NewCommentRepo(db)
	userRepo := repo.NewUserRepo(db)
	likeRepo := repo.NewLikeRepo(db)
	memoryRepo := repo.NewMemoryRepo(db)
	assembler := assembler.NewCommentAssembler()
	service := NewCommentService(commentRepo, userRepo, likeRepo, assembler, memoryRepo)

	ctx := context.Background()
	userID := int64(100)
	commentID := int64(1)

	t.Run("成功删除自己的顶级评论", func(t *testing.T) {
		// 1. 查询评论
		commentRows := sqlmock.NewRows([]string{"id", "user_id", "memory_id", "parent_id"}).
			AddRow(commentID, userID, 10, nil)
		mock.ExpectQuery("SELECT .* FROM `comments` WHERE id = ? AND `comments`.`deleted_at` IS NULL ORDER BY `comments`.`id` LIMIT ?").
			WithArgs(commentID, 1).
			WillReturnRows(commentRows)

		// 2. 软删除评论
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `comments` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), commentID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// 异步更新记忆评论数（不验证）

		err := service.DeleteComment(ctx, commentID, userID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功删除自己的回复", func(t *testing.T) {
		parentID := int64(5)
		commentRows := sqlmock.NewRows([]string{"id", "user_id", "memory_id", "parent_id"}).
			AddRow(commentID, userID, 10, &parentID)
		mock.ExpectQuery("SELECT .* FROM `comments` WHERE id = ? AND `comments`.`deleted_at` IS NULL ORDER BY `comments`.`id` LIMIT ?").
			WithArgs(commentID, 1).
			WillReturnRows(commentRows)

		// 删除评论
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `comments` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), commentID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// 异步更新父评论回复数（不验证）

		err := service.DeleteComment(ctx, commentID, userID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("评论不存在", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `comments` WHERE id = ? AND `comments`.`deleted_at` IS NULL ORDER BY `comments`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		err := service.DeleteComment(ctx, 999, userID)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err) // 原代码返回 err 而不是 errno
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("无权删除他人评论", func(t *testing.T) {
		otherUserID := int64(200)
		commentRows := sqlmock.NewRows([]string{"id", "user_id"}).AddRow(commentID, otherUserID)
		mock.ExpectQuery("SELECT .* FROM `comments` WHERE id = ? AND `comments`.`deleted_at` IS NULL ORDER BY `comments`.`id` LIMIT ?").
			WithArgs(commentID, 1).
			WillReturnRows(commentRows)

		err := service.DeleteComment(ctx, commentID, userID)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrForbidden, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
