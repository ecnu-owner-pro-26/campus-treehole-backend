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

// TestMemoryService_CreateMemory 测试创建记忆
func TestMemoryService_CreateMemory(t *testing.T) {
	db, mock := setupMockDB(t)
	memoryRepo := repo.NewMemoryRepo(db)
	userRepo := repo.NewUserRepo(db)
	locationRepo := repo.NewLocationRepo(db)
	likeRepo := repo.NewLikeRepo(db)
	imageRepo := repo.NewImageRepo(db)
	service := NewMemoryService(memoryRepo, userRepo, locationRepo, likeRepo, imageRepo)

	ctx := context.Background()
	creatorID := int64(100)
	locationID := int64(1)

	t.Run("成功创建记忆（无图片）", func(t *testing.T) {
		req := &dto.CreateMemoryRequest{
			Title:      "测试记忆",
			Content:    "内容",
			LocationID: &locationID,
			Tags:       []string{"tag1", "tag2"},
		}

		// 1. 验证地点存在
		locationRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(locationID, "测试地点")
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE id = ? AND `locations`.`deleted_at` IS NULL ORDER BY `locations`.`id` LIMIT ?").
			WithArgs(locationID, 1).
			WillReturnRows(locationRows)

		// 2. 创建记忆
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `memories`").
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // deleted_at
				req.Title,
				req.Content,
				"测试地点", // LocationName
				locationID,
				true, // IsPublic (默认)
				0,    // ViewCount
				0,    // LikeCount
				0,    // CommentCount
				req.Tags,
				1, // Status (假设默认发布)
				creatorID,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// 3. 获取创建者信息
		userRows := sqlmock.NewRows([]string{"id", "nickname", "avatar"}).AddRow(creatorID, "创建者", "avatar.jpg")
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs(creatorID, 1).
			WillReturnRows(userRows)

		// 4. 获取图片列表（无图片）
		mock.ExpectQuery("SELECT .* FROM `images` WHERE memory_id = ? AND `images`.`deleted_at` IS NULL").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		resp, err := service.CreateMemory(ctx, req, creatorID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int64(1), resp.ID)
		assert.Equal(t, req.Title, resp.Title)
		assert.Equal(t, creatorID, resp.Creator.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功创建记忆（带图片）", func(t *testing.T) {
		req := &dto.CreateMemoryRequest{
			Title:      "带图片记忆",
			Content:    "内容",
			LocationID: &locationID,
			ImageURLs:  []string{"url1", "url2"},
		}

		locationRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(locationID, "测试地点")
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE id = ? AND `locations`.`deleted_at` IS NULL ORDER BY `locations`.`id` LIMIT ?").
			WithArgs(locationID, 1).
			WillReturnRows(locationRows)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `memories`").
			WithArgs(
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				req.Title, req.Content, "测试地点", locationID, true, 0, 0, 0, "", 1, creatorID,
			).WillReturnResult(sqlmock.NewResult(2, 1))
		mock.ExpectCommit()

		// 插入两张图片
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `images`").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 2, "url1", int64(0), 0).
			WillReturnResult(sqlmock.NewResult(101, 1))
		mock.ExpectCommit()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `images`").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 2, "url2", int64(0), 1).
			WillReturnResult(sqlmock.NewResult(102, 1))
		mock.ExpectCommit()

		userRows := sqlmock.NewRows([]string{"id", "nickname", "avatar"}).AddRow(creatorID, "创建者", "avatar.jpg")
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs(creatorID, 1).
			WillReturnRows(userRows)

		// 获取图片列表（返回刚刚插入的两张）
		imageRows := sqlmock.NewRows([]string{"id", "url", "sort_order"}).
			AddRow(101, "url1", 0).
			AddRow(102, "url2", 1)
		mock.ExpectQuery("SELECT .* FROM `images` WHERE memory_id = ? AND `images`.`deleted_at` IS NULL").
			WithArgs(2).
			WillReturnRows(imageRows)

		resp, err := service.CreateMemory(ctx, req, creatorID)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(resp.Images))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("地点不存在", func(t *testing.T) {
		req := &dto.CreateMemoryRequest{LocationID: &locationID}
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE id = ? AND `locations`.`deleted_at` IS NULL ORDER BY `locations`.`id` LIMIT ?").
			WithArgs(locationID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		resp, err := service.CreateMemory(ctx, req, creatorID)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrLocationNotFound, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestMemoryService_GetMemory 测试获取记忆详情
func TestMemoryService_GetMemory(t *testing.T) {
	db, mock := setupMockDB(t)
	memoryRepo := repo.NewMemoryRepo(db)
	userRepo := repo.NewUserRepo(db)
	locationRepo := repo.NewLocationRepo(db)
	likeRepo := repo.NewLikeRepo(db)
	imageRepo := repo.NewImageRepo(db)
	service := NewMemoryService(memoryRepo, userRepo, locationRepo, likeRepo, imageRepo)

	ctx := context.Background()
	memoryID := int64(1)
	creatorID := int64(100)
	currentUserID := int64(200)

	t.Run("成功获取记忆（未登录）", func(t *testing.T) {
		// 1. 查询记忆
		memoryRows := sqlmock.NewRows([]string{"id", "title", "content", "creator_id", "view_count"}).
			AddRow(memoryID, "标题", "内容", creatorID, 5)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(memoryID, 1).
			WillReturnRows(memoryRows)

		// 2. 增加浏览次数
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `memories` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), memoryID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// 3. 查询创建者
		userRows := sqlmock.NewRows([]string{"id", "nickname", "avatar"}).AddRow(creatorID, "创建者", "a.jpg")
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs(creatorID, 1).
			WillReturnRows(userRows)

		// 4. 查询图片
		imageRows := sqlmock.NewRows([]string{"id", "url"}).
			AddRow(10, "img1.jpg").
			AddRow(11, "img2.jpg")
		mock.ExpectQuery("SELECT .* FROM `images` WHERE memory_id = ? AND `images`.`deleted_at` IS NULL").
			WithArgs(memoryID).
			WillReturnRows(imageRows)

		// 5. 未登录，不检查点赞

		resp, err := service.GetMemory(ctx, memoryID, nil)

		assert.NoError(t, err)
		assert.Equal(t, "标题", resp.Title)
		assert.Equal(t, 2, len(resp.Images))
		assert.False(t, resp.IsLiked)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功获取记忆（已登录）", func(t *testing.T) {
		memoryRows := sqlmock.NewRows([]string{"id", "title", "creator_id"}).AddRow(memoryID, "标题", creatorID)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(memoryID, 1).
			WillReturnRows(memoryRows)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `memories` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), memoryID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		userRows := sqlmock.NewRows([]string{"id", "nickname"}).AddRow(creatorID, "创建者")
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?").
			WithArgs(creatorID, 1).
			WillReturnRows(userRows)

		imageRows := sqlmock.NewRows([]string{"id", "url"}).AddRow(10, "img1.jpg")
		mock.ExpectQuery("SELECT .* FROM `images` WHERE memory_id = ? AND `images`.`deleted_at` IS NULL").
			WithArgs(memoryID).
			WillReturnRows(imageRows)

		// 检查点赞
		likeRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery("SELECT count\\(.*\\) FROM `likes` WHERE user_id = \\? AND target_id = \\? AND target_type = \\? AND `likes`.`deleted_at` IS NULL").
			WithArgs(currentUserID, memoryID, dto.LikeTargetTypeMemory).
			WillReturnRows(likeRows)

		resp, err := service.GetMemory(ctx, memoryID, &currentUserID)

		assert.NoError(t, err)
		assert.True(t, resp.IsLiked)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("记忆不存在", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		resp, err := service.GetMemory(ctx, 999, nil)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrMemoryNotFound, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestMemoryService_ListMemories 测试获取记忆列表
func TestMemoryService_ListMemories(t *testing.T) {
	db, mock := setupMockDB(t)
	memoryRepo := repo.NewMemoryRepo(db)
	userRepo := repo.NewUserRepo(db)
	locationRepo := repo.NewLocationRepo(db)
	likeRepo := repo.NewLikeRepo(db)
	imageRepo := repo.NewImageRepo(db)
	service := NewMemoryService(memoryRepo, userRepo, locationRepo, likeRepo, imageRepo)

	ctx := context.Background()
	locationID := int64(1)
	currentUserID := int64(200)

	t.Run("成功获取列表（无分页参数）", func(t *testing.T) {
		req := &dto.MemoryListRequest{
			LocationID: &locationID,
		}

		// 1. 查询记忆列表
		memoryRows := sqlmock.NewRows([]string{"id", "title", "creator_id"}).
			AddRow(1, "记忆1", 101).
			AddRow(2, "记忆2", 102)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE location_id = \\? AND `memories`.`deleted_at` IS NULL ORDER BY .* LIMIT 20 OFFSET 0").
			WithArgs(locationID).
			WillReturnRows(memoryRows)

		// 2. 查询总数
		mock.ExpectQuery("SELECT count\\(.*\\) FROM `memories` WHERE location_id = \\? AND `memories`.`deleted_at` IS NULL").
			WithArgs(locationID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		// 3. 查询创建者（两个用户）
		userRows1 := sqlmock.NewRows([]string{"id", "nickname"}).AddRow(101, "用户1")
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?").
			WithArgs(101, 1).
			WillReturnRows(userRows1)

		userRows2 := sqlmock.NewRows([]string{"id", "nickname"}).AddRow(102, "用户2")
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?").
			WithArgs(102, 1).
			WillReturnRows(userRows2)

		// 4. 查询图片
		imageRows1 := sqlmock.NewRows([]string{"id", "url"}).AddRow(10, "img1.jpg")
		mock.ExpectQuery("SELECT .* FROM `images` WHERE memory_id = \\? AND `images`.`deleted_at` IS NULL").
			WithArgs(1).
			WillReturnRows(imageRows1)

		imageRows2 := sqlmock.NewRows([]string{"id", "url"}).AddRow(20, "img2.jpg")
		mock.ExpectQuery("SELECT .* FROM `images` WHERE memory_id = \\? AND `images`.`deleted_at` IS NULL").
			WithArgs(2).
			WillReturnRows(imageRows2)

		// 5. 检查点赞（未登录，不检查）

		resp, err := service.ListMemories(ctx, req, nil)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(resp.Memories))
		assert.Equal(t, int64(2), resp.Total)
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 20, resp.PageSize)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功获取列表（已登录）", func(t *testing.T) {
		req := &dto.MemoryListRequest{
			Page:     2,
			PageSize: 10,
			SortBy:   "hot",
		}

		memoryRows := sqlmock.NewRows([]string{"id", "creator_id"}).
			AddRow(1, 101)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE `memories`.`deleted_at` IS NULL ORDER BY .* LIMIT 10 OFFSET 10").
			WillReturnRows(memoryRows)

		mock.ExpectQuery("SELECT count\\(.*\\) FROM `memories` WHERE `memories`.`deleted_at` IS NULL").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		userRows := sqlmock.NewRows([]string{"id", "nickname"}).AddRow(101, "用户1")
		mock.ExpectQuery("SELECT .* FROM `users` WHERE id = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?").
			WithArgs(101, 1).
			WillReturnRows(userRows)

		imageRows := sqlmock.NewRows([]string{"id", "url"}).AddRow(10, "img1.jpg")
		mock.ExpectQuery("SELECT .* FROM `images` WHERE memory_id = \\? AND `images`.`deleted_at` IS NULL").
			WithArgs(1).
			WillReturnRows(imageRows)

		// 检查点赞
		likeRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery("SELECT count\\(.*\\) FROM `likes` WHERE user_id = \\? AND target_id = \\? AND target_type = \\? AND `likes`.`deleted_at` IS NULL").
			WithArgs(currentUserID, 1, dto.LikeTargetTypeMemory).
			WillReturnRows(likeRows)

		resp, err := service.ListMemories(ctx, req, &currentUserID)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Memories))
		assert.True(t, resp.Memories[0].IsLiked)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		req := &dto.MemoryListRequest{Page: 1, PageSize: 10}
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE `memories`.`deleted_at` IS NULL ORDER BY .* LIMIT 10 OFFSET 0").
			WillReturnError(errors.New("db error"))

		resp, err := service.ListMemories(ctx, req, nil)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestMemoryService_UpdateMemory 测试更新记忆
func TestMemoryService_UpdateMemory(t *testing.T) {
	db, mock := setupMockDB(t)
	memoryRepo := repo.NewMemoryRepo(db)
	userRepo := repo.NewUserRepo(db)
	locationRepo := repo.NewLocationRepo(db)
	likeRepo := repo.NewLikeRepo(db)
	imageRepo := repo.NewImageRepo(db)
	service := NewMemoryService(memoryRepo, userRepo, locationRepo, likeRepo, imageRepo)

	ctx := context.Background()
	memoryID := int64(1)
	userID := int64(100) // 创建者
	otherUserID := int64(200)

	t.Run("成功更新记忆（不更新位置和图片）", func(t *testing.T) {
		req := &dto.UpdateMemoryRequest{
			Title:   strPtr("新标题"),
			Content: strPtr("新内容"),
		}

		// 1. 查询记忆
		memoryRows := sqlmock.NewRows([]string{"id", "creator_id", "location_id", "location_name"}).
			AddRow(memoryID, userID, 1, "旧地点")
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(memoryID, 1).
			WillReturnRows(memoryRows)

		// 2. 更新记忆
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `memories` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), "新标题", "新内容", memoryID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// 图片更新部分：req.ImageURLs == nil，不操作图片

		err := service.UpdateMemory(ctx, memoryID, req, userID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功更新记忆（更新位置）", func(t *testing.T) {
		newLocationID := int64(2)
		req := &dto.UpdateMemoryRequest{
			LocationID: &newLocationID,
		}

		// 1. 查询记忆
		memoryRows := sqlmock.NewRows([]string{"id", "creator_id", "location_id", "location_name"}).
			AddRow(memoryID, userID, 1, "旧地点")
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(memoryID, 1).
			WillReturnRows(memoryRows)

		// 2. 查询新地点
		locationRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(newLocationID, "新地点")
		mock.ExpectQuery("SELECT .* FROM `locations` WHERE id = ? AND `locations`.`deleted_at` IS NULL ORDER BY `locations`.`id` LIMIT ?").
			WithArgs(newLocationID, 1).
			WillReturnRows(locationRows)

		// 3. 更新记忆
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `memories` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), "新地点", newLocationID, memoryID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := service.UpdateMemory(ctx, memoryID, req, userID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功更新记忆（更新图片）", func(t *testing.T) {
		req := &dto.UpdateMemoryRequest{
			ImageURLs: []string{"new1.jpg", "new2.jpg"},
		}

		memoryRows := sqlmock.NewRows([]string{"id", "creator_id"}).AddRow(memoryID, userID)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(memoryID, 1).
			WillReturnRows(memoryRows)

		// 更新记忆（无字段更新，但 Update 方法会更新 updated_at）
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `memories` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), memoryID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// 删除旧图片
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `images` SET .* WHERE memory_id = ?").
			WithArgs(sqlmock.AnyArg(), memoryID).
			WillReturnResult(sqlmock.NewResult(0, 2))
		mock.ExpectCommit()

		// 插入新图片
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `images`").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), memoryID, "new1.jpg", int64(0), 0).
			WillReturnResult(sqlmock.NewResult(101, 1))
		mock.ExpectCommit()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `images`").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), memoryID, "new2.jpg", int64(0), 1).
			WillReturnResult(sqlmock.NewResult(102, 1))
		mock.ExpectCommit()

		err := service.UpdateMemory(ctx, memoryID, req, userID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("无权更新（不是创建者）", func(t *testing.T) {
		req := &dto.UpdateMemoryRequest{Title: strPtr("新标题")}
		memoryRows := sqlmock.NewRows([]string{"id", "creator_id"}).AddRow(memoryID, otherUserID)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(memoryID, 1).
			WillReturnRows(memoryRows)

		err := service.UpdateMemory(ctx, memoryID, req, userID)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrForbidden, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("记忆不存在", func(t *testing.T) {
		req := &dto.UpdateMemoryRequest{}
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		err := service.UpdateMemory(ctx, 999, req, userID)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrMemoryNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestMemoryService_DeleteMemory 测试删除记忆
func TestMemoryService_DeleteMemory(t *testing.T) {
	db, mock := setupMockDB(t)
	memoryRepo := repo.NewMemoryRepo(db)
	userRepo := repo.NewUserRepo(db)
	locationRepo := repo.NewLocationRepo(db)
	likeRepo := repo.NewLikeRepo(db)
	imageRepo := repo.NewImageRepo(db)
	service := NewMemoryService(memoryRepo, userRepo, locationRepo, likeRepo, imageRepo)

	ctx := context.Background()
	memoryID := int64(1)
	userID := int64(100)
	otherUserID := int64(200)

	t.Run("成功删除自己的记忆", func(t *testing.T) {
		memoryRows := sqlmock.NewRows([]string{"id", "creator_id"}).AddRow(memoryID, userID)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(memoryID, 1).
			WillReturnRows(memoryRows)

		// 软删除
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `memories` SET .* WHERE id = ?").
			WithArgs(sqlmock.AnyArg(), memoryID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := service.DeleteMemory(ctx, memoryID, userID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("无权删除他人记忆", func(t *testing.T) {
		memoryRows := sqlmock.NewRows([]string{"id", "creator_id"}).AddRow(memoryID, otherUserID)
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(memoryID, 1).
			WillReturnRows(memoryRows)

		err := service.DeleteMemory(ctx, memoryID, userID)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrForbidden, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("记忆不存在", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM `memories` WHERE id = ? AND `memories`.`deleted_at` IS NULL ORDER BY `memories`.`id` LIMIT ?").
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		err := service.DeleteMemory(ctx, 999, userID)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrMemoryNotFound, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// 辅助函数，返回字符串指针
func strPtr(s string) *string {
	return &s
}
