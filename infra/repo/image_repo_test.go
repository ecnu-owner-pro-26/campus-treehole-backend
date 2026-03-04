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
)

// TestImageRepo_Create 测试创建图片
func TestImageRepo_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewImageRepo(db)

	t.Run("成功创建", func(t *testing.T) {
		image := &model.ImageModel{
			MemoryID:  1,
			URL:       "https://example.com/image.jpg",
			Size:      1024,
			SortOrder: 0,
		}
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `image_models`")).
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // deleted_at
				image.MemoryID,
				image.URL,
				image.Size,
				image.SortOrder,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.Background(), image)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestImageRepo_GetByID 测试根据ID获取图片
func TestImageRepo_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewImageRepo(db)

	t.Run("成功获取", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "memory_id", "url", "size", "sort_order"}).
			AddRow(1, 1, "https://example.com/image.jpg", 1024, 0)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `image_models` WHERE id = $1 ORDER BY `image_models`.`id` LIMIT $2")).
			WithArgs(1, 1).
			WillReturnRows(rows)

		image, err := repo.GetByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, image)
		assert.Equal(t, int64(1), image.ID)
		assert.Equal(t, "https://example.com/image.jpg", image.URL)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("记录不存在", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `image_models` WHERE id = $1 ORDER BY `image_models`.`id` LIMIT $2")).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		image, err := repo.GetByID(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Nil(t, image)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `image_models` WHERE id = $1 ORDER BY `image_models`.`id` LIMIT $2")).
			WithArgs(1, 1).
			WillReturnError(errors.New("db error"))

		image, err := repo.GetByID(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, image)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestImageRepo_GetByMemoryID 测试根据记忆ID获取图片列表
func TestImageRepo_GetByMemoryID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewImageRepo(db)

	t.Run("成功获取列表", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "memory_id", "url", "size", "sort_order"}).
			AddRow(1, 1, "img1.jpg", 100, 0).
			AddRow(2, 1, "img2.jpg", 200, 1)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `image_models` WHERE memory_id = $1 ORDER BY sort_order ASC")).
			WithArgs(1).
			WillReturnRows(rows)

		images, err := repo.GetByMemoryID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Len(t, images, 2)
		assert.Equal(t, "img1.jpg", images[0].URL)
		assert.Equal(t, "img2.jpg", images[1].URL)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("无图片", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "memory_id", "url", "size", "sort_order"})
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `image_models` WHERE memory_id = $1 ORDER BY sort_order ASC")).
			WithArgs(1).
			WillReturnRows(rows)

		images, err := repo.GetByMemoryID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Empty(t, images)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `image_models` WHERE memory_id = $1 ORDER BY sort_order ASC")).
			WithArgs(1).
			WillReturnError(errors.New("db error"))

		images, err := repo.GetByMemoryID(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, images)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestImageRepo_Delete 测试删除单张图片
func TestImageRepo_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewImageRepo(db)

	t.Run("成功删除", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `image_models` WHERE `image_models`.`id` = ?")).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(context.Background(), 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("删除不存在的记录", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `image_models` WHERE `image_models`.`id` = ?")).
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.Delete(context.Background(), 999)
		assert.NoError(t, err) // GORM 的 Delete 即使影响行数为 0 也不返回错误
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `image_models` WHERE `image_models`.`id` = ?")).
			WithArgs(1).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.Delete(context.Background(), 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestImageRepo_DeleteByMemoryID 测试删除记忆的所有图片
func TestImageRepo_DeleteByMemoryID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewImageRepo(db)

	t.Run("成功删除", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `image_models` WHERE memory_id = ?")).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 3)) // 假设删除了 3 张图片
		mock.ExpectCommit()

		err := repo.DeleteByMemoryID(context.Background(), 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("删除不存在的记忆（无图片）", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `image_models` WHERE memory_id = ?")).
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.DeleteByMemoryID(context.Background(), 999)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `image_models` WHERE memory_id = ?")).
			WithArgs(1).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.DeleteByMemoryID(context.Background(), 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
