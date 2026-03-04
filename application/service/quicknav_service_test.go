package service

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestQuickNavService_BuildNavTree 测试构建导航树
func TestQuickNavService_BuildNavTree(t *testing.T) {
	db, mock := setupMockDB(t)
	service := NewQuickNavService(db)

	t.Run("成功构建所有校区树", func(t *testing.T) {
		// 模拟查询所有地点（is_active=1）
		rows := sqlmock.NewRows([]string{"id", "campus_id", "category", "name", "sort_order", "memory_count"}).
			AddRow(1, 1, "教学楼", "第一教学楼", 1, 100).
			AddRow(2, 1, "教学楼", "第二教学楼", 2, 80).
			AddRow(3, 1, "食堂", "第一食堂", 1, 200).
			AddRow(4, 2, "教学楼", "第三教学楼", 1, 150).
			AddRow(5, 2, "图书馆", "主图书馆", 1, 300)

		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE is_active = \\$1 ORDER BY campus_id, category, sort_order").
			WithArgs(1).
			WillReturnRows(rows)

		result, err := service.BuildNavTree(0) // 0 表示所有校区

		assert.NoError(t, err)
		assert.Len(t, result, 2) // 两个校区

		// 验证校区1
		assert.Equal(t, int64(1), result[0].CampusID)
		assert.Equal(t, "普陀校区", result[0].CampusName)
		assert.Len(t, result[0].Categories, 2) // 教学楼、食堂

		// 分类按字母排序（"教学楼" < "食堂"）
		assert.Equal(t, "教学楼", result[0].Categories[0].Category)
		assert.Len(t, result[0].Categories[0].Locations, 2)
		assert.Equal(t, "食堂", result[0].Categories[1].Category)
		assert.Len(t, result[0].Categories[1].Locations, 1)

		// 验证校区2
		assert.Equal(t, int64(2), result[1].CampusID)
		assert.Equal(t, "闵行校区", result[1].CampusName)
		assert.Len(t, result[1].Categories, 2) // 教学楼、图书馆

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功构建指定校区树", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "campus_id", "category", "name", "sort_order", "memory_count"}).
			AddRow(1, 1, "教学楼", "第一教学楼", 1, 100).
			AddRow(2, 1, "食堂", "第一食堂", 1, 200)

		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE is_active = \\$1 AND campus_id = \\$2 ORDER BY campus_id, category, sort_order").
			WithArgs(1, 1).
			WillReturnRows(rows)

		result, err := service.BuildNavTree(1)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), result[0].CampusID)
		assert.Len(t, result[0].Categories, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("无数据时返回空数组", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "campus_id", "category", "name"})
		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE is_active = \\$1 ORDER BY campus_id, category, sort_order").
			WithArgs(1).
			WillReturnRows(rows)

		result, err := service.BuildNavTree(0)

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE is_active = \\$1 ORDER BY campus_id, category, sort_order").
			WithArgs(1).
			WillReturnError(errors.New("db error"))

		result, err := service.BuildNavTree(0)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestQuickNavService_SearchLocations 测试搜索地点
func TestQuickNavService_SearchLocations(t *testing.T) {
	db, mock := setupMockDB(t)
	service := NewQuickNavService(db)

	t.Run("成功搜索（带关键词）", func(t *testing.T) {
		keyword := "教学楼"
		campusID := int64(1)
		page := 1
		pageSize := 10

		// 计数查询
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `locations` WHERE is_active = \\$1 AND name LIKE \\$2 ESCAPE '\\\\' AND campus_id = \\$3").
			WithArgs(1, "%教学楼%", campusID).
			WillReturnRows(countRows)

		// 数据查询
		dataRows := sqlmock.NewRows([]string{"id", "name", "memory_count"}).
			AddRow(1, "第一教学楼", 100).
			AddRow(2, "第二教学楼", 80)
		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE is_active = \\$1 AND name LIKE \\$2 ESCAPE '\\\\' AND campus_id = \\$3 ORDER BY memory_count desc LIMIT \\$4").
			WithArgs(1, "%教学楼%", campusID, pageSize).
			WillReturnRows(dataRows)

		locations, total, err := service.SearchLocations(keyword, campusID, page, pageSize)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, locations, 2)
		assert.Equal(t, "第一教学楼", locations[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功搜索（关键词含特殊字符）", func(t *testing.T) {
		keyword := "50%_off"
		campusID := int64(0) // 所有校区
		page := 1
		pageSize := 20

		// 转义后的关键词应该是：%50\%\_off%
		// 注意：在 SQL 中我们需要匹配最终生成的字符串，这里转义后的字符串可能根据驱动不同有差异，但大致是 "%50\\%\\_off%"
		// 我们假设最终生成的参数是 "%50\\%\\_off%"
		escapedParam := "%50\\%\\_off%"

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `locations` WHERE is_active = \\$1 AND name LIKE \\$2 ESCAPE '\\\\'").
			WithArgs(1, escapedParam).
			WillReturnRows(countRows)

		dataRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "50%_off 商店")
		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE is_active = \\$1 AND name LIKE \\$2 ESCAPE '\\\\' ORDER BY memory_count desc LIMIT \\$3").
			WithArgs(1, escapedParam, pageSize).
			WillReturnRows(dataRows)

		locations, total, err := service.SearchLocations(keyword, campusID, page, pageSize)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, locations, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("无关键词时返回所有", func(t *testing.T) {
		campusID := int64(0)
		page := 1
		pageSize := 10

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(5)
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `locations` WHERE is_active = \\$1").
			WithArgs(1).
			WillReturnRows(countRows)

		dataRows := sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2).AddRow(3).AddRow(4).AddRow(5)
		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE is_active = \\$1 ORDER BY memory_count desc LIMIT \\$2").
			WithArgs(1, pageSize).
			WillReturnRows(dataRows)

		locations, total, err := service.SearchLocations("", campusID, page, pageSize)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), total)
		assert.Len(t, locations, 5)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("分页参数边界", func(t *testing.T) {
		// page < 1 应自动修正为 1
		// pageSize 超出范围应修正为 20
		campusID := int64(0)
		page := 0
		pageSize := 200

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(10)
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `locations` WHERE is_active = \\$1").
			WithArgs(1).
			WillReturnRows(countRows)

		// 修正后 page=1, pageSize=20, offset=0
		dataRows := sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2)
		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE is_active = \\$1 ORDER BY memory_count desc LIMIT \\$2").
			WithArgs(1, 20).
			WillReturnRows(dataRows)

		locations, total, err := service.SearchLocations("", campusID, page, pageSize)

		assert.NoError(t, err)
		assert.Equal(t, int64(10), total)
		assert.Len(t, locations, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `locations` WHERE is_active = \\$1").
			WithArgs(1).
			WillReturnError(errors.New("count error"))

		locations, total, err := service.SearchLocations("", 0, 1, 10)

		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, locations)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestQuickNavService_GetPopularLocations 测试获取热门地点
func TestQuickNavService_GetPopularLocations(t *testing.T) {
	db, mock := setupMockDB(t)
	service := NewQuickNavService(db)

	t.Run("成功获取（指定校区）", func(t *testing.T) {
		campusID := int64(1)
		limit := 5

		rows := sqlmock.NewRows([]string{"id", "name", "memory_count"}).
			AddRow(1, "地点1", 100).
			AddRow(2, "地点2", 90)
		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE is_active = \\$1 AND campus_id = \\$2 ORDER BY memory_count desc LIMIT \\$3").
			WithArgs(1, campusID, limit).
			WillReturnRows(rows)

		locations, err := service.GetPopularLocations(campusID, limit)

		assert.NoError(t, err)
		assert.Len(t, locations, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("成功获取（所有校区）", func(t *testing.T) {
		campusID := int64(0)
		limit := 10

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2).AddRow(3)
		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE is_active = \\$1 ORDER BY memory_count desc LIMIT \\$2").
			WithArgs(1, limit).
			WillReturnRows(rows)

		locations, err := service.GetPopularLocations(campusID, limit)

		assert.NoError(t, err)
		assert.Len(t, locations, 3)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("limit 边界（超出范围自动修正为10）", func(t *testing.T) {
		campusID := int64(0)
		limit := 100 // 会修正为 10

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2)
		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE is_active = \\$1 ORDER BY memory_count desc LIMIT \\$2").
			WithArgs(1, 10). // 修正后的 limit
			WillReturnRows(rows)

		locations, err := service.GetPopularLocations(campusID, limit)

		assert.NoError(t, err)
		assert.Len(t, locations, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("limit 小于0时使用默认值", func(t *testing.T) {
		campusID := int64(0)
		limit := -5 // 会修正为 10

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE is_active = \\$1 ORDER BY memory_count desc LIMIT \\$2").
			WithArgs(1, 10).
			WillReturnRows(rows)

		locations, err := service.GetPopularLocations(campusID, limit)

		assert.NoError(t, err)
		assert.Len(t, locations, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE is_active = \\$1 ORDER BY memory_count desc LIMIT \\$2").
			WithArgs(1, 10).
			WillReturnError(errors.New("db error"))

		locations, err := service.GetPopularLocations(0, 10)

		assert.Error(t, err)
		assert.Nil(t, locations)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestQuickNavService_GetLocationsByCategory 测试按类别获取地点
func TestQuickNavService_GetLocationsByCategory(t *testing.T) {
	db, mock := setupMockDB(t)
	service := NewQuickNavService(db)

	t.Run("成功获取（分页）", func(t *testing.T) {
		campusID := int64(1)
		category := "教学楼"
		page := 2
		pageSize := 5

		// 计数
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(12)
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `locations` WHERE campus_id = \\$1 AND category = \\$2 AND is_active = \\$3").
			WithArgs(campusID, category, 1).
			WillReturnRows(countRows)

		// 数据查询
		rows := sqlmock.NewRows([]string{"id", "name", "sort_order", "memory_count"}).
			AddRow(6, "教学楼6", 1, 50).
			AddRow(7, "教学楼7", 2, 40)
		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE campus_id = \\$1 AND category = \\$2 AND is_active = \\$3 ORDER BY sort_order asc, memory_count desc LIMIT \\$4 OFFSET \\$5").
			WithArgs(campusID, category, 1, pageSize, (page-1)*pageSize).
			WillReturnRows(rows)

		locations, total, err := service.GetLocationsByCategory(campusID, category, page, pageSize)

		assert.NoError(t, err)
		assert.Equal(t, int64(12), total)
		assert.Len(t, locations, 2)
		assert.Equal(t, "教学楼6", locations[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("分页参数边界（自动修正）", func(t *testing.T) {
		campusID := int64(1)
		category := "食堂"
		page := 0       // 修正为1
		pageSize := 200 // 修正为20

		// 计数
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(3)
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `locations` WHERE campus_id = \\$1 AND category = \\$2 AND is_active = \\$3").
			WithArgs(campusID, category, 1).
			WillReturnRows(countRows)

		// 数据查询：limit 20 offset 0
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2).AddRow(3)
		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE campus_id = \\$1 AND category = \\$2 AND is_active = \\$3 ORDER BY sort_order asc, memory_count desc LIMIT \\$4 OFFSET \\$5").
			WithArgs(campusID, category, 1, 20, 0).
			WillReturnRows(rows)

		locations, total, err := service.GetLocationsByCategory(campusID, category, page, pageSize)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, locations, 3)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("无数据返回空切片", func(t *testing.T) {
		campusID := int64(1)
		category := "图书馆"

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `locations` WHERE campus_id = \\$1 AND category = \\$2 AND is_active = \\$3").
			WithArgs(campusID, category, 1).
			WillReturnRows(countRows)

		rows := sqlmock.NewRows([]string{"id"})
		mock.ExpectQuery("SELECT \\* FROM `locations` WHERE campus_id = \\$1 AND category = \\$2 AND is_active = \\$3 ORDER BY sort_order asc, memory_count desc LIMIT \\$4 OFFSET \\$5").
			WithArgs(campusID, category, 1, 20, 0).
			WillReturnRows(rows)

		locations, total, err := service.GetLocationsByCategory(campusID, category, 1, 20)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, locations)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库错误", func(t *testing.T) {
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `locations` WHERE campus_id = \\$1 AND category = \\$2 AND is_active = \\$3").
			WithArgs(1, "教学楼", 1).
			WillReturnError(errors.New("db error"))

		locations, total, err := service.GetLocationsByCategory(1, "教学楼", 1, 10)

		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, locations)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
