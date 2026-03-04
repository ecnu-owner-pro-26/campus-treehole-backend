package repo

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"campus-memory/infra/model"
)

// setupMockDB 创建 mock 数据库连接，返回 db, mock 和 repo 实例
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := sqlite.Dialector{
		DriverName: "sqlite3",
		DSN:        "file::memory:?cache=shared",
		Conn:       mockDB,
	}
	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return db, mock // 返回 db，而不是 repo
}

// TestCampusRepo_GetAll 测试获取所有校区
func TestCampusRepo_GetAll(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCampusRepo(db)

	t.Run("成功获取", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "is_active", "sort_order"}).
			AddRow(1, "普陀校区", true, 1).
			AddRow(2, "闵行校区", true, 2)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `campus_models` WHERE is_active = ? ORDER BY sort_order ASC, id ASC")).
			WithArgs(true).
			WillReturnRows(rows)

		campuses, err := repo.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Len(t, campuses, 2)
		assert.Equal(t, "普陀校区", campuses[0].Name)
		assert.Equal(t, "闵行校区", campuses[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("无数据", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "is_active", "sort_order"})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `campus_models` WHERE is_active = ? ORDER BY sort_order ASC, id ASC")).
			WithArgs(true).
			WillReturnRows(rows)

		campuses, err := repo.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, campuses)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `campus_models` WHERE is_active = ? ORDER BY sort_order ASC, id ASC")).
			WithArgs(true).
			WillReturnError(errors.New("db error"))

		campuses, err := repo.GetAll(context.Background())
		assert.Error(t, err)
		assert.Nil(t, campuses)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCampusRepo_GetByID 测试根据 ID 获取校区
func TestCampusRepo_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCampusRepo(db)

	t.Run("成功获取", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "is_active", "sort_order"}).
			AddRow(1, "普陀校区", true, 1)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `campus_models` WHERE id = ? ORDER BY `campus_models`.`id` LIMIT ?")).
			WithArgs(1, 1).
			WillReturnRows(rows)

		campus, err := repo.GetByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, campus)
		assert.Equal(t, int64(1), campus.ID)
		assert.Equal(t, "普陀校区", campus.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("记录不存在", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `campus_models` WHERE id = ? ORDER BY `campus_models`.`id` LIMIT ?")).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		campus, err := repo.GetByID(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Nil(t, campus)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `campus_models` WHERE id = ? ORDER BY `campus_models`.`id` LIMIT ?")).
			WithArgs(1, 1).
			WillReturnError(errors.New("db error"))

		campus, err := repo.GetByID(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, campus)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCampusRepo_Create 测试创建校区
func TestCampusRepo_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCampusRepo(db)

	t.Run("成功创建", func(t *testing.T) {
		campus := &model.CampusModel{
			Name:      "新校区",
			IsActive:  true,
			SortOrder: 3,
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `campus_models`")).
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // deleted_at
				campus.Name,
				campus.IsActive,
				campus.SortOrder,
			).WillReturnResult(sqlmock.NewResult(10, 1))
		mock.ExpectCommit()

		err := repo.Create(context.Background(), campus)
		assert.NoError(t, err)
		assert.Equal(t, int64(10), campus.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("创建失败", func(t *testing.T) {
		campus := &model.CampusModel{
			Name: "失败校区",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `campus_models`")).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				campus.Name,
				campus.IsActive,
				campus.SortOrder,
			).WillReturnError(errors.New("duplicate key"))
		mock.ExpectRollback()

		err := repo.Create(context.Background(), campus)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCampusRepo_Update 测试更新校区
func TestCampusRepo_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCampusRepo(db)

	t.Run("成功更新", func(t *testing.T) {
		campus := &model.CampusModel{
			ID:        1,
			Name:      "更新后的名称",
			IsActive:  true,
			SortOrder: 5,
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `campus_models` SET")).
			WithArgs(
				sqlmock.AnyArg(), // updated_at
				campus.Name,
				campus.IsActive,
				campus.SortOrder,
				campus.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.Background(), campus)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("更新失败（记录不存在）", func(t *testing.T) {
		campus := &model.CampusModel{
			ID:        999,
			Name:      "不存在的校区",
			IsActive:  true,
			SortOrder: 1,
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `campus_models` SET")).
			WithArgs(
				sqlmock.AnyArg(),
				campus.Name,
				campus.IsActive,
				campus.SortOrder,
				campus.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.Update(context.Background(), campus)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		campus := &model.CampusModel{
			ID:   1,
			Name: "更新错误",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `campus_models` SET")).
			WithArgs(
				sqlmock.AnyArg(),
				campus.Name,
				campus.IsActive,
				campus.SortOrder,
				campus.ID,
			).WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.Update(context.Background(), campus)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestCampusRepo_Delete 测试删除校区
func TestCampusRepo_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewCampusRepo(db)

	t.Run("成功删除", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `campus_models` WHERE `campus_models`.`id` = ?")).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(context.Background(), 1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("删除不存在的记录", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `campus_models` WHERE `campus_models`.`id` = ?")).
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.Delete(context.Background(), 999)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `campus_models` WHERE `campus_models`.`id` = ?")).
			WithArgs(1).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.Delete(context.Background(), 1)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
