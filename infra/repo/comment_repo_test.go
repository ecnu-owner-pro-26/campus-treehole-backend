package repo

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"campus-memory/infra/model"
	"campus-memory/types/errno"
)

// TestCommentRepo_Create 测试创建评论
func TestCommentRepo_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepo(db)

	t.Run("成功创建", func(t *testing.T) {
		comment := &model.CommentModel{
			MemoryID: 1,
			UserID:   100,
			Content:  "这是一条评论",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `comment_models`")).
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // deleted_at
				comment.MemoryID,
				comment.UserID,
				comment.Content,
				comment.ParentID,
				comment.ReplyToUserID,
				comment.LikeCount,
				comment.Status,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.Background(), comment)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("comment 为 nil", func(t *testing.T) {
		err := repo.Create(context.Background(), nil)
		assert.Error(t, err)
		assert.Equal(t, "comment is nil", err.Error())
	})

	t.Run("数据库错误", func(t *testing.T) {
		comment := &model.CommentModel{
			MemoryID: 1,
			UserID:   100,
			Content:  "测试",
		}
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `comment_models`")).
			WithArgs(
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				comment.MemoryID, comment.UserID, comment.Content,
				comment.ParentID, comment.ReplyToUserID,
				comment.LikeCount, comment.Status,
			).WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.Create(context.Background(), comment)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCommentRepo_ListByMemoryID 测试根据记忆ID获取评论列表
func TestCommentRepo_ListByMemoryID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepo(db)

	t.Run("成功获取（最新）", func(t *testing.T) {
		memoryID := int64(1)
		page := 1
		pageSize := 10
		sortBy := "latest"

		// 计数查询
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `comment_models` WHERE memory_id = $1 AND status = $2 AND deleted_at IS NULL AND parent_id IS NULL")).
			WithArgs(memoryID, 1).
			WillReturnRows(countRows)

		// 数据查询
		rows := sqlmock.NewRows([]string{"id", "content", "like_count"}).
			AddRow(1, "评论1", 10).
			AddRow(2, "评论2", 5)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `comment_models` WHERE memory_id = $1 AND status = $2 AND deleted_at IS NULL AND parent_id IS NULL ORDER BY created_at DESC LIMIT $3 OFFSET $4")).
			WithArgs(memoryID, 1, pageSize, (page-1)*pageSize).
			WillReturnRows(rows)

		comments, total, err := repo.ListByMemoryID(context.Background(), memoryID, page, pageSize, sortBy)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, comments, 2)
		assert.Equal(t, "评论1", comments[0].Content)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功获取（热门）", func(t *testing.T) {
		memoryID := int64(1)
		page := 2
		pageSize := 5
		sortBy := "hot"

		// 计数
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(3)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `comment_models` WHERE memory_id = $1 AND status = $2 AND deleted_at IS NULL AND parent_id IS NULL")).
			WithArgs(memoryID, 1).
			WillReturnRows(countRows)

		// 数据查询
		rows := sqlmock.NewRows([]string{"id", "content"}).AddRow(3, "评论3")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `comment_models` WHERE memory_id = $1 AND status = $2 AND deleted_at IS NULL AND parent_id IS NULL ORDER BY like_count DESC, created_at DESC LIMIT $3 OFFSET $4")).
			WithArgs(memoryID, 1, pageSize, (page-1)*pageSize).
			WillReturnRows(rows)

		comments, total, err := repo.ListByMemoryID(context.Background(), memoryID, page, pageSize, sortBy)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, comments, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("计数失败", func(t *testing.T) {
		memoryID := int64(1)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `comment_models` WHERE memory_id = $1 AND status = $2 AND deleted_at IS NULL AND parent_id IS NULL")).
			WithArgs(memoryID, 1).
			WillReturnError(errors.New("count error"))

		comments, total, err := repo.ListByMemoryID(context.Background(), memoryID, 1, 10, "latest")
		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, comments)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCommentRepo_ListRepliesByParentID 测试根据父评论ID获取回复列表
func TestCommentRepo_ListRepliesByParentID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepo(db)

	t.Run("成功获取", func(t *testing.T) {
		parentID := int64(10)
		page := 1
		pageSize := 20

		// 计数
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(3)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `comment_models` WHERE parent_id = $1 AND status = $2 AND deleted_at IS NULL")).
			WithArgs(parentID, 1).
			WillReturnRows(countRows)

		// 数据
		rows := sqlmock.NewRows([]string{"id", "content"}).
			AddRow(101, "回复1").
			AddRow(102, "回复2").
			AddRow(103, "回复3")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `comment_models` WHERE parent_id = $1 AND status = $2 AND deleted_at IS NULL ORDER BY created_at ASC LIMIT $3 OFFSET $4")).
			WithArgs(parentID, 1, pageSize, (page-1)*pageSize).
			WillReturnRows(rows)

		replies, total, err := repo.ListRepliesByParentID(context.Background(), parentID, page, pageSize)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, replies, 3)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("总数为0", func(t *testing.T) {
		parentID := int64(10)
		// 计数返回 0
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `comment_models` WHERE parent_id = $1 AND status = $2 AND deleted_at IS NULL")).
			WithArgs(parentID, 1).
			WillReturnRows(countRows)

		replies, total, err := repo.ListRepliesByParentID(context.Background(), parentID, 1, 20)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, replies)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("计数失败", func(t *testing.T) {
		parentID := int64(10)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `comment_models` WHERE parent_id = $1 AND status = $2 AND deleted_at IS NULL")).
			WithArgs(parentID, 1).
			WillReturnError(errors.New("count error"))

		replies, total, err := repo.ListRepliesByParentID(context.Background(), parentID, 1, 20)
		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, replies)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCommentRepo_Delete 测试删除评论（软删除）
func TestCommentRepo_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepo(db)

	t.Run("成功删除", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `deleted_at`=CURRENT_TIMESTAMP WHERE id = $1 AND `comment_models`.`deleted_at` IS NULL")).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(context.Background(), 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `deleted_at`=CURRENT_TIMESTAMP WHERE id = $1 AND `comment_models`.`deleted_at` IS NULL")).
			WithArgs(1).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.Delete(context.Background(), 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCommentRepo_GetByID 测试根据ID获取评论
func TestCommentRepo_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepo(db)

	t.Run("成功获取", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "content"}).AddRow(1, "测试评论")
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `comment_models` WHERE id = $1 AND deleted_at IS NULL AND status IN ($2) ORDER BY `comment_models`.`id` LIMIT $3")).
			WithArgs(1, 1, 1).
			WillReturnRows(rows)

		comment, err := repo.GetByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, comment)
		assert.Equal(t, int64(1), comment.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("记录不存在", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `comment_models` WHERE id = $1 AND deleted_at IS NULL AND status IN ($2) ORDER BY `comment_models`.`id` LIMIT $3")).
			WithArgs(999, 1, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		comment, err := repo.GetByID(context.Background(), 999)
		assert.NoError(t, err) // 注意：方法中处理了 ErrRecordNotFound 返回 nil,nil
		assert.Nil(t, comment)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `comment_models` WHERE id = $1 AND deleted_at IS NULL AND status IN ($2) ORDER BY `comment_models`.`id` LIMIT $3")).
			WithArgs(1, 1, 1).
			WillReturnError(errors.New("db error"))

		comment, err := repo.GetByID(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, comment)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCommentRepo_GetReplyCount 测试获取回复数量
func TestCommentRepo_GetReplyCount(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepo(db)

	t.Run("成功获取", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `comment_models` WHERE parent_id = $1 AND status = $2")).
			WithArgs(1, 1).
			WillReturnRows(rows)

		count, err := repo.GetReplyCount(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `comment_models` WHERE parent_id = $1 AND status = $2")).
			WithArgs(1, 1).
			WillReturnError(errors.New("db error"))

		count, err := repo.GetReplyCount(context.Background(), 1)
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCommentRepo_BatchGetReplyCount 测试批量获取回复数量
func TestCommentRepo_BatchGetReplyCount(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepo(db)

	t.Run("成功获取", func(t *testing.T) {
		commentIDs := []int64{1, 2, 3}
		rows := sqlmock.NewRows([]string{"parent_id", "count"}).
			AddRow(1, 3).
			AddRow(2, 5)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT parent_id, COUNT(*) as count FROM `comment_models` WHERE parent_id IN ($1,$2,$3) AND status = $4 GROUP BY parent_id")).
			WithArgs(1, 2, 3, 1).
			WillReturnRows(rows)

		countMap, err := repo.BatchGetReplyCount(context.Background(), commentIDs)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), countMap[1])
		assert.Equal(t, int64(5), countMap[2])
		_, ok := countMap[3]
		assert.False(t, ok) // 没有结果应为 0 但不包含在 map 中
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("空数组", func(t *testing.T) {
		countMap, err := repo.BatchGetReplyCount(context.Background(), []int64{})
		assert.NoError(t, err)
		assert.Empty(t, countMap)
	})

	t.Run("数据库错误", func(t *testing.T) {
		commentIDs := []int64{1}
		mock.ExpectQuery(regexp.QuoteMeta("SELECT parent_id, COUNT(*) as count FROM `comment_models` WHERE parent_id IN ($1) AND status = $2 GROUP BY parent_id")).
			WithArgs(1, 1).
			WillReturnError(errors.New("db error"))

		countMap, err := repo.BatchGetReplyCount(context.Background(), commentIDs)
		assert.Error(t, err)
		assert.Nil(t, countMap)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCommentRepo_IncrementReplyCount 测试增加回复数
func TestCommentRepo_IncrementReplyCount(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepo(db)

	t.Run("成功增加", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `reply_count`=reply_count + $1 WHERE id = $2 AND status = $3 AND deleted_at IS NULL")).
			WithArgs(1, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.IncrementReplyCount(context.Background(), 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("记录不存在", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `reply_count`=reply_count + $1 WHERE id = $2 AND status = $3 AND deleted_at IS NULL")).
			WithArgs(1, 999, 1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.IncrementReplyCount(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrCommentNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `reply_count`=reply_count + $1 WHERE id = $2 AND status = $3 AND deleted_at IS NULL")).
			WithArgs(1, 1, 1).
			WillReturnError(errors.New("db error"))

		err := repo.IncrementReplyCount(context.Background(), 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCommentRepo_DecrementReplyCount 测试减少回复数
func TestCommentRepo_DecrementReplyCount(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepo(db)

	t.Run("成功减少", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `reply_count`=reply_count - $1 WHERE id = $2 AND reply_count > $3")).
			WithArgs(1, 1, 0).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DecrementReplyCount(context.Background(), 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("记录不存在或已为0", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `reply_count`=reply_count - $1 WHERE id = $2 AND reply_count > $3")).
			WithArgs(1, 999, 0).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DecrementReplyCount(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrCommentNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `reply_count`=reply_count - $1 WHERE id = $2 AND reply_count > $3")).
			WithArgs(1, 1, 0).
			WillReturnError(errors.New("db error"))

		err := repo.DecrementReplyCount(context.Background(), 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCommentRepo_IncrementLikeCount 测试增加点赞数
func TestCommentRepo_IncrementLikeCount(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepo(db)

	t.Run("成功增加", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `like_count`=like_count + $1 WHERE id = $2")).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.IncrementLikeCount(context.Background(), 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("记录不存在", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `like_count`=like_count + $1 WHERE id = $2")).
			WithArgs(1, 999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.IncrementLikeCount(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrCommentNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `like_count`=like_count + $1 WHERE id = $2")).
			WithArgs(1, 1).
			WillReturnError(errors.New("db error"))

		err := repo.IncrementLikeCount(context.Background(), 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCommentRepo_DecrementLikeCount 测试减少点赞数
func TestCommentRepo_DecrementLikeCount(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepo(db)

	t.Run("成功减少", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `like_count`=like_count - $1 WHERE id = $2 AND like_count > $3")).
			WithArgs(1, 1, 0).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DecrementLikeCount(context.Background(), 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("记录不存在或已为0", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `like_count`=like_count - $1 WHERE id = $2 AND like_count > $3")).
			WithArgs(1, 999, 0).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DecrementLikeCount(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrCommentNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `like_count`=like_count - $1 WHERE id = $2 AND like_count > $3")).
			WithArgs(1, 1, 0).
			WillReturnError(errors.New("db error"))

		err := repo.DecrementLikeCount(context.Background(), 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCommentRepo_UpdateLikeCount 测试更新点赞数（增量）
func TestCommentRepo_UpdateLikeCount(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCommentRepo(db)

	t.Run("成功更新（delta 不为空）", func(t *testing.T) {
		delta := int64(5)
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `like_count`=like_count + $1 WHERE id = $2")).
			WithArgs(delta, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateLikeCount(context.Background(), 1, &delta)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功更新（delta 为空）", func(t *testing.T) {
		// delta 为 nil 时，应该没有更新，直接返回 nil
		err := repo.UpdateLikeCount(context.Background(), 1, nil)
		assert.NoError(t, err)
		// 没有数据库操作，不需要 mock
	})

	t.Run("数据库错误", func(t *testing.T) {
		delta := int64(1)
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `comment_models` SET `like_count`=like_count + $1 WHERE id = $2")).
			WithArgs(delta, 1).
			WillReturnError(errors.New("db error"))

		err := repo.UpdateLikeCount(context.Background(), 1, &delta)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
