package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"campus-memory/application/dto"
	"campus-memory/infra/repo"
	"campus-memory/utils"
)

// 保存原始函数以便恢复
var (
	originalGetWechatOpenID = utils.GetWechatOpenID
	originalGenerateToken   = utils.GenerateToken
)

// setupMockDB 创建 mock 数据库连接（使用 SQLite dialector + sqlmock）
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	// 使用 sqlite.Dialector，将 mock 的 *sql.DB 注入 Conn 字段
	dialector := sqlite.Dialector{
		DriverName: "sqlite3",
		DSN:        "file::memory:?cache=shared",
		Conn:       mockDB,
	}
	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)
	return db, mock
}

// TestAuthService_WechatLogin 测试微信登录
func TestAuthService_WechatLogin(t *testing.T) {
	db, mock := setupMockDB(t)
	userRepo := repo.NewUserRepo(db) // 假设存在此构造函数
	authService := NewAuthService(userRepo)

	// mock 微信接口
	mockWechat := func(openID, unionID string, err error) {
		utils.GetWechatOpenID = func(code string) (*utils.WechatSession, error) {
			return &utils.WechatSession{OpenID: openID, UnionID: unionID}, err
		}
	}
	defer func() { utils.GetWechatOpenID = originalGetWechatOpenID }()

	// mock token 生成
	mockToken := func(token string, err error) {
		utils.GenerateToken = func(userID int64, openID string) (string, error) {
			return token, err
		}
	}
	defer func() { utils.GenerateToken = originalGenerateToken }()

	t.Run("新用户登录成功", func(t *testing.T) {
		mockWechat("open123", "union123", nil)
		mockToken("test-token", nil)

		req := &dto.WechatLoginRequest{
			Code:     "code123",
			Nickname: "张三",
			Avatar:   "avatar.jpg",
		}

		// 查询用户不存在
		mock.ExpectQuery("SELECT .* FROM `users` WHERE open_id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs("open123", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// 创建用户
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `users`").
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // deleted_at
				"open123",
				"union123",
				"张三",
				"avatar.jpg",
				int8(1), // Status
				int8(0), // Role
				nil,     // DefaultCampusID
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		resp, err := authService.WechatLogin(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "test-token", resp.Token)
		assert.Equal(t, int64(1), resp.User.ID)
		assert.Equal(t, "张三", resp.User.Nickname)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("已有用户登录（无需更新）", func(t *testing.T) {
		mockWechat("open123", "", nil)
		mockToken("test-token", nil)

		req := &dto.WechatLoginRequest{Code: "code123"}

		rows := sqlmock.NewRows([]string{
			"id", "open_id", "union_id", "nickname", "avatar", "status", "role", "default_campus_id", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, "open123", "", "李四", "old.jpg", int8(1), int8(0), nil, time.Now(), time.Now(), nil,
		)
		mock.ExpectQuery("SELECT .* FROM `users` WHERE open_id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs("open123", 1).
			WillReturnRows(rows)

		resp, err := authService.WechatLogin(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, "test-token", resp.Token)
		assert.Equal(t, "李四", resp.User.Nickname)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("已有用户登录并更新信息", func(t *testing.T) {
		mockWechat("open123", "newUnion", nil)
		mockToken("test-token", nil)

		req := &dto.WechatLoginRequest{
			Code:     "code123",
			Nickname: "王五",
			Avatar:   "new.jpg",
		}

		rows := sqlmock.NewRows([]string{
			"id", "open_id", "union_id", "nickname", "avatar", "status", "role", "default_campus_id", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			1, "open123", "oldUnion", "李四", "old.jpg", int8(1), int8(0), nil, time.Now(), time.Now(), nil,
		)
		mock.ExpectQuery("SELECT .* FROM `users` WHERE open_id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs("open123", 1).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `users` SET .* WHERE id = ?").
			WithArgs(
				sqlmock.AnyArg(), // updated_at
				"王五",
				"new.jpg",
				"newUnion",
				1, // id
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		resp, err := authService.WechatLogin(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, "王五", resp.User.Nickname)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("微信API调用失败", func(t *testing.T) {
		mockWechat("", "", errors.New("微信接口错误"))
		req := &dto.WechatLoginRequest{Code: "code123"}
		resp, err := authService.WechatLogin(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

// TestAuthService_GetUserProfile 测试获取用户信息
func TestAuthService_GetUserProfile(t *testing.T) {
	db, mock := setupMockDB(t)
	userRepo := repo.NewUserRepo(db)
	authService := NewAuthService(userRepo)

	t.Run("成功获取用户信息", func(t *testing.T) {
		userID := int64(1)
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"id", "nickname", "avatar", "default_campus_id", "status", "role", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			userID, "张三", "avatar.jpg", nil, int8(1), int8(0), now, now, nil,
		)
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs(userID, 1).
			WillReturnRows(rows)

		resp, err := authService.GetUserProfile(context.Background(), userID)

		assert.NoError(t, err)
		assert.Equal(t, userID, resp.ID)
		assert.Equal(t, "张三", resp.Nickname)
		assert.Equal(t, "avatar.jpg", resp.Avatar)
		assert.Equal(t, int8(1), resp.Status)
		assert.Equal(t, int8(0), resp.Role)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("用户不存在", func(t *testing.T) {
		userID := int64(999)
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs(userID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		resp, err := authService.GetUserProfile(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "用户不存在")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestAuthService_UpdateUserProfile 测试更新用户信息
func TestAuthService_UpdateUserProfile(t *testing.T) {
	db, mock := setupMockDB(t)
	userRepo := repo.NewUserRepo(db)
	authService := NewAuthService(userRepo)

	t.Run("成功更新所有字段", func(t *testing.T) {
		userID := int64(1)
		nickname := "新昵称"
		avatar := "new.jpg"
		campusID := int64(2)
		req := &dto.UpdateProfileRequest{
			Nickname:        &nickname,
			Avatar:          &avatar,
			DefaultCampusID: &campusID,
		}

		rows := sqlmock.NewRows([]string{
			"id", "nickname", "avatar", "default_campus_id", "status", "role", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			userID, "旧昵称", "old.jpg", nil, int8(1), int8(0), time.Now(), time.Now(), nil,
		)
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs(userID, 1).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `users` SET .* WHERE id = ?").
			WithArgs(
				sqlmock.AnyArg(), // updated_at
				"新昵称",
				"new.jpg",
				campusID,
				userID,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := authService.UpdateUserProfile(context.Background(), userID, req)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("只更新昵称", func(t *testing.T) {
		userID := int64(1)
		nickname := "新昵称"
		req := &dto.UpdateProfileRequest{Nickname: &nickname}

		rows := sqlmock.NewRows([]string{
			"id", "nickname", "avatar", "default_campus_id", "status", "role", "created_at", "updated_at", "deleted_at",
		}).AddRow(
			userID, "旧昵称", "old.jpg", nil, int8(1), int8(0), time.Now(), time.Now(), nil,
		)
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs(userID, 1).
			WillReturnRows(rows)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `users` SET .* WHERE id = ?").
			WithArgs(
				sqlmock.AnyArg(), // updated_at
				"新昵称",
				userID,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := authService.UpdateUserProfile(context.Background(), userID, req)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("用户不存在", func(t *testing.T) {
		userID := int64(999)
		req := &dto.UpdateProfileRequest{}

		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs(userID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		err := authService.UpdateUserProfile(context.Background(), userID, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "用户不存在")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
