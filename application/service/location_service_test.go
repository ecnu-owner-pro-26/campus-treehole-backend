package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"campus-memory/application/dto"
	"campus-memory/infra/repo"
	"campus-memory/types/errno"
)

// TestLocationService_GetLocation 测试获取地点详情
func TestLocationService_GetLocation(t *testing.T) {
	db, mock := setupMockDB(t)
	locationRepo := repo.NewLocationRepo(db)
	campusRepo := repo.NewCampusRepo(db)
	service := NewLocationService(locationRepo, campusRepo)

	ctx := context.Background()
	locationID := int64(1)

	t.Run("成功获取地点", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "campus_id", "name", "description", "category", "created_at", "updated_at", "deleted_at"}).
			AddRow(locationID, 10, "图书馆", "学校图书馆", "building", time.Now(), time.Now(), nil)
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE id = ? AND `locations`.`deleted_at` IS NULL ORDER BY `locations`.`id` LIMIT ?").
			WithArgs(locationID, 1).
			WillReturnRows(rows)

		resp, err := service.GetLocation(ctx, locationID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, locationID, resp.ID)
		assert.Equal(t, "图书馆", resp.Name)
		assert.Equal(t, "building", resp.Category)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("地点不存在", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE id = ? AND `locations`.`deleted_at` IS NULL ORDER BY `locations`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		resp, err := service.GetLocation(ctx, 999)

		assert.Error(t, err)
		assert.Equal(t, errno.ErrLocationNotFound, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE id = ? AND `locations`.`deleted_at` IS NULL ORDER BY `locations`.`id` LIMIT ?").
			WithArgs(locationID, 1).
			WillReturnError(errors.New("db error"))

		resp, err := service.GetLocation(ctx, locationID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestLocationService_ListLocations 测试获取地点列表
func TestLocationService_ListLocations(t *testing.T) {
	db, mock := setupMockDB(t)
	locationRepo := repo.NewLocationRepo(db)
	campusRepo := repo.NewCampusRepo(db)
	service := NewLocationService(locationRepo, campusRepo)

	ctx := context.Background()
	campusID := int64(10)
	category := "building"

	t.Run("成功获取列表（带分页和筛选）", func(t *testing.T) {
		req := &dto.LocationListRequest{
			CampusID: &campusID,
			Category: category,
			Page:     1,
			PageSize: 10,
		}

		// 查询列表
		locationRows := sqlmock.NewRows([]string{"id", "campus_id", "name", "category"}).
			AddRow(1, campusID, "图书馆", "building").
			AddRow(2, campusID, "教学楼A", "building")
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE campus_id = ? AND category = ? AND `locations`.`deleted_at` IS NULL ORDER BY .* LIMIT 10 OFFSET 0").
			WithArgs(campusID, category).
			WillReturnRows(locationRows)

		// 查询总数
		mock.ExpectQuery("SELECT count\\(.*\\) FROM `locations` WHERE campus_id = \\? AND category = \\? AND `locations`.`deleted_at` IS NULL").
			WithArgs(campusID, category).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		resp, err := service.ListLocations(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 2, len(resp.Locations))
		assert.Equal(t, int64(2), resp.Total)
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 10, resp.PageSize)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功获取列表（无筛选条件）", func(t *testing.T) {
		req := &dto.LocationListRequest{
			Page:     1,
			PageSize: 20,
		}

		locationRows := sqlmock.NewRows([]string{"id", "campus_id", "name"}).
			AddRow(1, 10, "图书馆").
			AddRow(2, 20, "食堂")
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE `locations`.`deleted_at` IS NULL ORDER BY .* LIMIT 20 OFFSET 0").
			WillReturnRows(locationRows)

		mock.ExpectQuery("SELECT count\\(.*\\) FROM `locations` WHERE `locations`.`deleted_at` IS NULL").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		resp, err := service.ListLocations(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(resp.Locations))
		assert.Equal(t, int64(2), resp.Total)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("默认分页值", func(t *testing.T) {
		req := &dto.LocationListRequest{} // page 和 pageSize 为 0，应被设为默认值

		locationRows := sqlmock.NewRows([]string{"id"})
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE `locations`.`deleted_at` IS NULL ORDER BY .* LIMIT 20 OFFSET 0").
			WillReturnRows(locationRows)
		mock.ExpectQuery("SELECT count\\(.*\\) FROM `locations` WHERE `locations`.`deleted_at` IS NULL").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		resp, err := service.ListLocations(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, 1, resp.Page)      // 默认 page = 1
		assert.Equal(t, 20, resp.PageSize) // 默认 pageSize = 20
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		req := &dto.LocationListRequest{Page: 1, PageSize: 10}
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE `locations`.`deleted_at` IS NULL ORDER BY .* LIMIT 10 OFFSET 0").
			WillReturnError(errors.New("db error"))

		resp, err := service.ListLocations(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestLocationService_CreateLocation 测试创建地点
func TestLocationService_CreateLocation(t *testing.T) {
	db, mock := setupMockDB(t)
	locationRepo := repo.NewLocationRepo(db)
	campusRepo := repo.NewCampusRepo(db)
	service := NewLocationService(locationRepo, campusRepo)

	ctx := context.Background()
	campusID := int64(10)

	t.Run("成功创建地点", func(t *testing.T) {
		req := &dto.CreateLocationRequest{
			CampusID: campusID,
			Name:     "新地点",
			Category: "building",
		}

		// 验证校区存在
		campusRows := sqlmock.NewRows([]string{"id"}).AddRow(campusID)
		mock.ExpectQuery("SELECT .* FROM `campuses` WHERE id = ? AND `campuses`.`deleted_at` IS NULL ORDER BY `campuses`.`id` LIMIT ?").
			WithArgs(campusID, 1).
			WillReturnRows(campusRows)

		// 创建地点
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `locations`").
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // deleted_at
				campusID,
				req.Name,
				req.Category,
			).WillReturnResult(sqlmock.NewResult(100, 1))
		mock.ExpectCommit()

		resp, err := service.CreateLocation(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int64(100), resp.ID)
		assert.Equal(t, req.Name, resp.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("校区不存在", func(t *testing.T) {
		req := &dto.CreateLocationRequest{CampusID: 999}
		mock.ExpectQuery("SELECT .* FROM `campuses` WHERE id = ? AND `campuses`.`deleted_at` IS NULL ORDER BY `campuses`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		resp, err := service.CreateLocation(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, errno.ErrCampusNotFound, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("创建失败", func(t *testing.T) {
		req := &dto.CreateLocationRequest{CampusID: campusID}
		campusRows := sqlmock.NewRows([]string{"id"}).AddRow(campusID)
		mock.ExpectQuery("SELECT .* FROM `campuses` WHERE id = ? AND `campuses`.`deleted_at` IS NULL ORDER BY `campuses`.`id` LIMIT ?").
			WithArgs(campusID, 1).
			WillReturnRows(campusRows)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `locations`").
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				campusID,
				req.Name,
				req.Category,
			).WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		resp, err := service.CreateLocation(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, errno.ErrLocationCreateFail, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestLocationService_UpdateLocation 测试更新地点
func TestLocationService_UpdateLocation(t *testing.T) {
	db, mock := setupMockDB(t)
	locationRepo := repo.NewLocationRepo(db)
	campusRepo := repo.NewCampusRepo(db)
	service := NewLocationService(locationRepo, campusRepo)

	ctx := context.Background()
	locationID := int64(1)

	t.Run("成功更新地点", func(t *testing.T) {
		name := "新名称"
		description := "新描述"
		req := &dto.UpdateLocationRequest{
			Name: &name,
		}

		// 查询地点
		locationRows := sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(locationID, "旧名称", "旧描述")
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE id = ? AND `locations`.`deleted_at` IS NULL ORDER BY `locations`.`id` LIMIT ?").
			WithArgs(locationID, 1).
			WillReturnRows(locationRows)

		// 更新
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `locations` SET .* WHERE id = ?").
			WithArgs(
				sqlmock.AnyArg(), // updated_at
				name,
				description,
				locationID,
			).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := service.UpdateLocation(ctx, locationID, req)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("只更新部分字段", func(t *testing.T) {
		name := "新名称"
		req := &dto.UpdateLocationRequest{Name: &name}

		locationRows := sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(locationID, "旧名称", "旧描述")
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE id = ? AND `locations`.`deleted_at` IS NULL ORDER BY `locations`.`id` LIMIT ?").
			WithArgs(locationID, 1).
			WillReturnRows(locationRows)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `locations` SET .* WHERE id = ?").
			WithArgs(
				sqlmock.AnyArg(),
				name,
				locationID,
			).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := service.UpdateLocation(ctx, locationID, req)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("地点不存在", func(t *testing.T) {
		req := &dto.UpdateLocationRequest{}
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE id = ? AND `locations`.`deleted_at` IS NULL ORDER BY `locations`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		err := service.UpdateLocation(ctx, 999, req)

		assert.Error(t, err)
		assert.Equal(t, errno.ErrLocationNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("更新失败", func(t *testing.T) {
		name := "新名称"
		req := &dto.UpdateLocationRequest{Name: &name}
		locationRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(locationID, "旧名称")
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE id = ? AND `locations`.`deleted_at` IS NULL ORDER BY `locations`.`id` LIMIT ?").
			WithArgs(locationID, 1).
			WillReturnRows(locationRows)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `locations` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), name, locationID).
			WillReturnError(errors.New("update error"))
		mock.ExpectRollback()

		err := service.UpdateLocation(ctx, locationID, req)

		assert.Error(t, err)
		assert.Equal(t, errno.ErrLocationUpdateFail, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestLocationService_DeleteLocation 测试删除地点
func TestLocationService_DeleteLocation(t *testing.T) {
	db, mock := setupMockDB(t)
	locationRepo := repo.NewLocationRepo(db)
	campusRepo := repo.NewCampusRepo(db)
	service := NewLocationService(locationRepo, campusRepo)

	ctx := context.Background()
	locationID := int64(1)

	t.Run("成功删除地点", func(t *testing.T) {
		// 查询地点是否存在
		rows := sqlmock.NewRows([]string{"id"}).AddRow(locationID)
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE id = ? AND `locations`.`deleted_at` IS NULL ORDER BY `locations`.`id` LIMIT ?").
			WithArgs(locationID, 1).
			WillReturnRows(rows)

		// 软删除
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `locations` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), locationID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := service.DeleteLocation(ctx, locationID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("地点不存在", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE id = ? AND `locations`.`deleted_at` IS NULL ORDER BY `locations`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		err := service.DeleteLocation(ctx, 999)

		assert.Error(t, err)
		assert.Equal(t, errno.ErrLocationNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("删除失败", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(locationID)
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE id = ? AND `locations`.`deleted_at` IS NULL ORDER BY `locations`.`id` LIMIT ?").
			WithArgs(locationID, 1).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `locations` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), locationID).
			WillReturnError(errors.New("delete error"))
		mock.ExpectRollback()

		err := service.DeleteLocation(ctx, locationID)

		assert.Error(t, err)
		assert.Equal(t, errno.ErrLocationDeleteFail, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestLocationService_SearchLocations 测试搜索地点
func TestLocationService_SearchLocations(t *testing.T) {
	db, mock := setupMockDB(t)
	locationRepo := repo.NewLocationRepo(db)
	campusRepo := repo.NewCampusRepo(db)
	service := NewLocationService(locationRepo, campusRepo)

	ctx := context.Background()
	keyword := "图书馆"
	campusID := int64(10)

	t.Run("成功搜索（指定校区）", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "campus_id", "name", "category"}).
			AddRow(1, campusID, "图书馆", "building").
			AddRow(2, campusID, "图书馆分馆", "building")
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE \\(name LIKE \\? OR description LIKE \\?\\) AND campus_id = \\? AND `locations`.`deleted_at` IS NULL").
			WithArgs("%"+keyword+"%", "%"+keyword+"%", campusID).
			WillReturnRows(rows)

		resp, err := service.SearchLocations(ctx, keyword, &campusID)

		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		assert.Equal(t, "图书馆", resp[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功搜索（不指定校区）", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "campus_id", "name"}).
			AddRow(1, 10, "图书馆").
			AddRow(3, 20, "图书馆分馆")
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE \\(name LIKE \\? OR description LIKE \\?\\) AND `locations`.`deleted_at` IS NULL").
			WithArgs("%"+keyword+"%", "%"+keyword+"%").
			WillReturnRows(rows)

		resp, err := service.SearchLocations(ctx, keyword, nil)

		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("搜索结果为空", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE \\(name LIKE \\? OR description LIKE \\?\\) AND `locations`.`deleted_at` IS NULL").
			WithArgs("%"+keyword+"%", "%"+keyword+"%").
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		resp, err := service.SearchLocations(ctx, keyword, nil)

		assert.NoError(t, err)
		assert.Len(t, resp, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE \\(name LIKE \\? OR description LIKE \\?\\) AND `locations`.`deleted_at` IS NULL").
			WithArgs("%"+keyword+"%", "%"+keyword+"%").
			WillReturnError(errors.New("search error"))

		resp, err := service.SearchLocations(ctx, keyword, nil)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
