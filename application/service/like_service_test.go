package service

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"campus-memory/application/dto"
	"campus-memory/infra/repo"
	"campus-memory/types/errno"
)

// TestLikeService_ToggleLike 测试切换点赞状态
func TestLikeService_ToggleLike(t *testing.T) {
	db, mock := setupMockDB(t)
	likeRepo := repo.NewLikeRepo(db)
	memoryRepo := repo.NewMemoryRepo(db)
	commentRepo := repo.NewCommentRepo(db)
	service := NewLikeService(likeRepo, memoryRepo, commentRepo)

	ctx := context.Background()
	userID := int64(100)
	targetID := int64(1)

	// 辅助函数：设置 GetByID 成功（记忆）
	expectMemoryExists := func() {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(targetID)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(targetID, 1).
			WillReturnRows(rows)
	}

	// 辅助函数：设置 GetByID 成功（评论）
	expectCommentExists := func() {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(targetID)
		mock.ExpectQuery("SELECT .* FROM `comments` WHERE id = ? AND `comments`.`deleted_at` IS NULL ORDER BY `comments`.`id` LIMIT ?").
			WithArgs(targetID, 1).
			WillReturnRows(rows)
	}

	// 辅助函数：设置 CheckLiked 返回指定状态
	expectCheckLiked := func(isLiked bool) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(map[bool]int{true: 1, false: 0}[isLiked])
		mock.ExpectQuery("SELECT count\\(.*\\) FROM `likes` WHERE user_id = \\? AND target_id = \\? AND target_type = \\? AND `likes`.`deleted_at` IS NULL").
			WithArgs(userID, targetID, sqlmock.AnyArg()).
			WillReturnRows(rows)
	}

	// 辅助函数：设置 DeleteLike 成功
	expectDeleteLike := func() {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `likes` SET .* WHERE user_id = \\? AND target_id = \\? AND target_type = \\?").
			WithArgs(sqlmock.AnyArg(), userID, targetID, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
	}

	// 辅助函数：设置 CreateLike 成功
	expectCreateLike := func() {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `likes`").
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // deleted_at
				userID,
				targetID,
				sqlmock.AnyArg(), // targetType
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	// 辅助函数：设置 UpdateCounts（记忆）成功
	expectUpdateMemoryCounts := func(delta int64) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `memories` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), delta, targetID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
	}

	// 辅助函数：设置 UpdateLikeCount（评论）成功
	expectUpdateCommentLikeCount := func(delta int64) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `comments` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), delta, targetID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
	}

	// 辅助函数：设置 GetLikeCount 返回指定数量
	expectGetLikeCount := func(count int64) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(count)
		mock.ExpectQuery("SELECT count\\(.*\\) FROM `likes` WHERE target_id = \\? AND target_type = \\? AND `likes`.`deleted_at` IS NULL").
			WithArgs(targetID, sqlmock.AnyArg()).
			WillReturnRows(rows)
	}

	t.Run("成功点赞记忆", func(t *testing.T) {
		expectMemoryExists()
		expectCheckLiked(false) // 未点赞
		expectCreateLike()
		expectUpdateMemoryCounts(1)
		expectGetLikeCount(1) // 点赞后总数

		resp, err := service.ToggleLike(ctx, userID, targetID, dto.LikeTargetTypeMemory)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.IsLiked)
		assert.Equal(t, int64(1), resp.LikeCount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功取消点赞记忆", func(t *testing.T) {
		expectMemoryExists()
		expectCheckLiked(true) // 已点赞
		expectDeleteLike()
		expectUpdateMemoryCounts(-1)
		expectGetLikeCount(0)

		resp, err := service.ToggleLike(ctx, userID, targetID, dto.LikeTargetTypeMemory)

		assert.NoError(t, err)
		assert.False(t, resp.IsLiked)
		assert.Equal(t, int64(0), resp.LikeCount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功点赞评论", func(t *testing.T) {
		expectCommentExists()
		expectCheckLiked(false)
		expectCreateLike()
		expectUpdateCommentLikeCount(1)
		expectGetLikeCount(1)

		resp, err := service.ToggleLike(ctx, userID, targetID, dto.LikeTargetTypeComment)

		assert.NoError(t, err)
		assert.True(t, resp.IsLiked)
		assert.Equal(t, int64(1), resp.LikeCount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功取消点赞评论", func(t *testing.T) {
		expectCommentExists()
		expectCheckLiked(true)
		expectDeleteLike()
		expectUpdateCommentLikeCount(-1)
		expectGetLikeCount(0)

		resp, err := service.ToggleLike(ctx, userID, targetID, dto.LikeTargetTypeComment)

		assert.NoError(t, err)
		assert.False(t, resp.IsLiked)
		assert.Equal(t, int64(0), resp.LikeCount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("目标记忆不存在", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(targetID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		resp, err := service.ToggleLike(ctx, userID, targetID, dto.LikeTargetTypeMemory)

		assert.Error(t, err)
		assert.Equal(t, errno.ErrMemoryNotFound, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("目标评论不存在", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `comments` WHERE id = ? AND `comments`.`deleted_at` IS NULL ORDER BY `comments`.`id` LIMIT ?").
			WithArgs(targetID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		resp, err := service.ToggleLike(ctx, userID, targetID, dto.LikeTargetTypeComment)

		assert.Error(t, err)
		assert.Equal(t, errno.ErrCommentNotFound, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("无效的 targetType", func(t *testing.T) {
		resp, err := service.ToggleLike(ctx, userID, targetID, 99)

		assert.Error(t, err)
		assert.Equal(t, errno.ErrBadRequest, err)
		assert.Nil(t, resp)
		// 没有数据库操作，不需要 mock 期望
	})

	// 错误场景：CheckLiked 失败
	t.Run("CheckLiked 失败", func(t *testing.T) {
		expectMemoryExists()
		mock.ExpectQuery("SELECT count\\(.*\\) FROM `likes` WHERE user_id = \\? AND target_id = \\? AND target_type = \\? AND `likes`.`deleted_at` IS NULL").
			WithArgs(userID, targetID, dto.LikeTargetTypeMemory).
			WillReturnError(errors.New("db error"))

		resp, err := service.ToggleLike(ctx, userID, targetID, dto.LikeTargetTypeMemory)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// 错误场景：CreateLike 失败
	t.Run("CreateLike 失败", func(t *testing.T) {
		expectMemoryExists()
		expectCheckLiked(false)
		// 模拟 CreateLike 失败
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `likes`").
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				userID,
				targetID,
				dto.LikeTargetTypeMemory,
			).WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		resp, err := service.ToggleLike(ctx, userID, targetID, dto.LikeTargetTypeMemory)

		assert.Error(t, err)
		assert.Equal(t, errno.ErrLikeCreateFail, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// 错误场景：DeleteLike 失败
	t.Run("DeleteLike 失败", func(t *testing.T) {
		expectMemoryExists()
		expectCheckLiked(true)
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `likes` SET .* WHERE user_id = \\? AND target_id = \\? AND target_type = \\?").
			WithArgs(sqlmock.AnyArg(), userID, targetID, dto.LikeTargetTypeMemory).
			WillReturnError(errors.New("delete error"))
		mock.ExpectRollback()

		resp, err := service.ToggleLike(ctx, userID, targetID, dto.LikeTargetTypeMemory)

		assert.Error(t, err)
		assert.Equal(t, errno.ErrLikeDeleteFail, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// 错误场景：更新记忆计数失败（虽然被忽略，但测试仍需模拟，因为代码会执行）
	t.Run("更新记忆计数失败（忽略）", func(t *testing.T) {
		expectMemoryExists()
		expectCheckLiked(false)
		expectCreateLike()
		// 模拟 UpdateCounts 失败
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `memories` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), 1, targetID).
			WillReturnError(errors.New("update error"))
		mock.ExpectRollback()
	})

	// 错误场景：GetLikeCount 失败
	t.Run("GetLikeCount 失败", func(t *testing.T) {
		expectMemoryExists()
		expectCheckLiked(false)
		expectCreateLike()
		expectUpdateMemoryCounts(1)
		// 模拟 GetLikeCount 失败
		mock.ExpectQuery("SELECT count\\(.*\\) FROM `likes` WHERE target_id = \\? AND target_type = \\? AND `likes`.`deleted_at` IS NULL").
			WithArgs(targetID, dto.LikeTargetTypeMemory).
			WillReturnError(errors.New("count error"))

		resp, err := service.ToggleLike(ctx, userID, targetID, dto.LikeTargetTypeMemory)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "count error")
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
