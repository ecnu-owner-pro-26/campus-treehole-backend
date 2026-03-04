package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"campus-memory/infra/repo"
	"campus-memory/types/errno"
)

// TestCampusService_ListCampuses 测试获取校区列表
func TestCampusService_ListCampuses(t *testing.T) {
	db, mock := setupMockDB(t)
	campusRepo := repo.NewCampusRepo(db)
	locationRepo := repo.NewLocationRepo(db)
	service := NewCampusService(campusRepo, locationRepo)

	t.Run("成功获取列表", func(t *testing.T) {
		// 模拟查询所有校区
		rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "校区A", "描述A", time.Now(), time.Now(), nil).
			AddRow(2, "校区B", "描述B", time.Now(), time.Now(), nil)
		mock.ExpectQuery("SELECT .* FROM `campuses` WHERE `campuses`.`deleted_at` IS NULL").
			WillReturnRows(rows)

		resp, err := service.ListCampuses(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 2, len(resp.Campuses))
		assert.Equal(t, int64(2), resp.Total)
		assert.Equal(t, "校区A", resp.Campuses[0].Name)
		assert.Equal(t, "校区B", resp.Campuses[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `campuses` WHERE `campuses`.`deleted_at` IS NULL").
			WillReturnError(errors.New("db error"))

		resp, err := service.ListCampuses(context.Background())

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCampusService_GetCampus 测试获取单个校区
func TestCampusService_GetCampus(t *testing.T) {
	db, mock := setupMockDB(t)
	campusRepo := repo.NewCampusRepo(db)
	locationRepo := repo.NewLocationRepo(db)
	service := NewCampusService(campusRepo, locationRepo)

	t.Run("成功获取", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "校区A", "描述A", time.Now(), time.Now(), nil)
		mock.ExpectQuery("SELECT .* FROM `campuses` WHERE id = ? AND `campuses`.`deleted_at` IS NULL ORDER BY `campuses`.`id` LIMIT ?").
			WithArgs(1, 1).
			WillReturnRows(rows)

		resp, err := service.GetCampus(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int64(1), resp.ID)
		assert.Equal(t, "校区A", resp.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("校区不存在", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `campuses` WHERE id = ? AND `campuses`.`deleted_at` IS NULL ORDER BY `campuses`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		resp, err := service.GetCampus(context.Background(), 999)

		assert.Error(t, err)
		assert.Equal(t, errno.ErrCampusNotFound, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `campuses` WHERE id = ? AND `campuses`.`deleted_at` IS NULL ORDER BY `campuses`.`id` LIMIT ?").
			WithArgs(1, 1).
			WillReturnError(errors.New("db error"))

		resp, err := service.GetCampus(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCampusService_GetCampusWithLocations 测试获取校区及地点列表
func TestCampusService_GetCampusWithLocations(t *testing.T) {
	db, mock := setupMockDB(t)
	campusRepo := repo.NewCampusRepo(db)
	locationRepo := repo.NewLocationRepo(db)
	service := NewCampusService(campusRepo, locationRepo)

	t.Run("成功获取校区及地点", func(t *testing.T) {
		// 模拟查询校区
		campusRows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "校区A", "描述A", time.Now(), time.Now(), nil)
		mock.ExpectQuery("SELECT .* FROM `campuses` WHERE id = ? AND `campuses`.`deleted_at` IS NULL ORDER BY `campuses`.`id` LIMIT ?").
			WithArgs(1, 1).
			WillReturnRows(campusRows)

		// 模拟查询地点
		locationRows := sqlmock.NewRows([]string{"id", "campus_id", "name", "description", "created_at", "updated_at", "deleted_at"}).
			AddRow(101, 1, "地点1", "地点1描述", time.Now(), time.Now(), nil).
			AddRow(102, 1, "地点2", "地点2描述", time.Now(), time.Now(), nil)
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE campus_id = ? AND `locations`.`deleted_at` IS NULL").
			WithArgs(1).
			WillReturnRows(locationRows)

		resp, err := service.GetCampusWithLocations(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int64(1), resp.Campus.ID)
		assert.Equal(t, "校区A", resp.Campus.Name)
		assert.Equal(t, 2, len(resp.Locations))
		assert.Equal(t, int64(101), resp.Locations[0].ID)
		assert.Equal(t, "地点1", resp.Locations[0].Name)
		assert.Equal(t, int64(102), resp.Locations[1].ID)
		assert.Equal(t, "地点2", resp.Locations[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("校区不存在", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `campuses` WHERE id = ? AND `campuses`.`deleted_at` IS NULL ORDER BY `campuses`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		resp, err := service.GetCampusWithLocations(context.Background(), 999)

		assert.Error(t, err)
		assert.Equal(t, errno.ErrCampusNotFound, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("获取地点失败", func(t *testing.T) {
		// 校区查询成功
		campusRows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "校区A", "描述A", time.Now(), time.Now(), nil)
		mock.ExpectQuery("SELECT .* FROM `campuses` WHERE id = ? AND `campuses`.`deleted_at` IS NULL ORDER BY `campuses`.`id` LIMIT ?").
			WithArgs(1, 1).
			WillReturnRows(campusRows)

		// 地点查询返回错误
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE campus_id = ? AND `locations`.`deleted_at` IS NULL").
			WithArgs(1).
			WillReturnError(errors.New("db error"))

		resp, err := service.GetCampusWithLocations(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
