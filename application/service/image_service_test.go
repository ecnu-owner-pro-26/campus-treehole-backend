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

// TestImageService_UploadImage 测试上传图片
func TestImageService_UploadImage(t *testing.T) {
	db, mock := setupMockDB(t)
	imageRepo := repo.NewImageRepo(db)
	memoryRepo := repo.NewMemoryRepo(db)
	service := NewImageService(imageRepo, memoryRepo)

	ctx := context.Background()
	memoryID := int64(1)
	userID := int64(100)

	t.Run("成功上传图片", func(t *testing.T) {
		req := &dto.UploadImageRequest{
			MemoryID: memoryID,
		}
		url := "https://example.com/image.jpg"
		size := int64(1024)

		// 1. 查询记忆是否存在
		memoryRows := sqlmock.NewRows([]string{"id", "creator_id"}).AddRow(memoryID, userID)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(memoryID, 1).
			WillReturnRows(memoryRows)

		// 2. 查询当前记忆的图片列表（用于确定 sort_order）
		imageRows := sqlmock.NewRows([]string{"id"}).
			AddRow(10).
			AddRow(11)
		mock.ExpectQuery("SELECT .* FROM `images` WHERE memory_id = ? AND `images`.`deleted_at` IS NULL").
			WithArgs(memoryID).
			WillReturnRows(imageRows)

		// 3. 创建图片记录
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `images`").
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // deleted_at
				memoryID,
				url,
				size,
				2, // sort_order = len(images) = 2
			).WillReturnResult(sqlmock.NewResult(100, 1))
		mock.ExpectCommit()

		resp, err := service.UploadImage(ctx, req, url, size)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int64(100), resp.ID)
		assert.Equal(t, url, resp.URL)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("记忆不存在", func(t *testing.T) {
		req := &dto.UploadImageRequest{MemoryID: 999}
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		resp, err := service.UploadImage(ctx, req, "url", 1024)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrMemoryNotFound, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestImageService_DeleteImage 测试删除图片
func TestImageService_DeleteImage(t *testing.T) {
	db, mock := setupMockDB(t)
	imageRepo := repo.NewImageRepo(db)
	memoryRepo := repo.NewMemoryRepo(db)
	service := NewImageService(imageRepo, memoryRepo)

	ctx := context.Background()
	userID := int64(100)
	imageID := int64(1)

	t.Run("成功删除自己的图片", func(t *testing.T) {
		// 1. 查询图片
		imageRows := sqlmock.NewRows([]string{"id", "memory_id"}).AddRow(imageID, 10)
		mock.ExpectQuery("SELECT .* FROM `images` WHERE id = ? AND `images`.`deleted_at` IS NULL ORDER BY `images`.`id` LIMIT ?").
			WithArgs(imageID, 1).
			WillReturnRows(imageRows)

		// 2. 查询记忆，验证权限
		memoryRows := sqlmock.NewRows([]string{"id", "creator_id"}).AddRow(10, userID)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(10, 1).
			WillReturnRows(memoryRows)

		// 3. 删除图片
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `images` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), imageID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := service.DeleteImage(ctx, imageID, userID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("图片不存在", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `images` WHERE id = ? AND `images`.`deleted_at` IS NULL ORDER BY `images`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		err := service.DeleteImage(ctx, 999, userID)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrImageNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("无权删除他人记忆的图片", func(t *testing.T) {
		imageRows := sqlmock.NewRows([]string{"id", "memory_id"}).AddRow(imageID, 10)
		mock.ExpectQuery("SELECT .* FROM `images` WHERE id = ? AND `images`.`deleted_at` IS NULL ORDER BY `images`.`id` LIMIT ?").
			WithArgs(imageID, 1).
			WillReturnRows(imageRows)

		memoryRows := sqlmock.NewRows([]string{"id", "creator_id"}).AddRow(10, 200) // 其他用户
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(10, 1).
			WillReturnRows(memoryRows)

		err := service.DeleteImage(ctx, imageID, userID)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrForbidden, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestImageService_GetImagesByMemoryID 测试获取记忆的图片列表
func TestImageService_GetImagesByMemoryID(t *testing.T) {
	db, mock := setupMockDB(t)
	imageRepo := repo.NewImageRepo(db)
	memoryRepo := repo.NewMemoryRepo(db)
	service := NewImageService(imageRepo, memoryRepo)

	ctx := context.Background()
	memoryID := int64(1)

	t.Run("成功获取图片列表", func(t *testing.T) {
		imageRows := sqlmock.NewRows([]string{"id", "memory_id", "url"}).
			AddRow(1, memoryID, "url1").
			AddRow(2, memoryID, "url2")
		mock.ExpectQuery("SELECT .* FROM `images` WHERE memory_id = ? AND `images`.`deleted_at` IS NULL").
			WithArgs(memoryID).
			WillReturnRows(imageRows)

		images, err := service.GetImagesByMemoryID(ctx, memoryID)

		assert.NoError(t, err)
		assert.Len(t, images, 2)
		assert.Equal(t, int64(1), images[0].ID)
		assert.Equal(t, "url1", images[0].URL)
		assert.Equal(t, int64(2), images[1].ID)
		assert.Equal(t, "url2", images[1].URL)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("记忆无图片", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `images` WHERE memory_id = ? AND `images`.`deleted_at` IS NULL").
			WithArgs(memoryID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "memory_id", "url"}))

		images, err := service.GetImagesByMemoryID(ctx, memoryID)

		assert.NoError(t, err)
		assert.Len(t, images, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `images` WHERE memory_id = ? AND `images`.`deleted_at` IS NULL").
			WithArgs(memoryID).
			WillReturnError(errors.New("db error"))

		images, err := service.GetImagesByMemoryID(ctx, memoryID)

		assert.Error(t, err)
		assert.Nil(t, images)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
