package handler

import (
	"campus-memory/application/dto"
	"campus-memory/infra/model"
	"campus-memory/types/errno"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockQuickNavService 是 QuickNavServiceInterface 的 mock 实现
type mockQuickNavService struct {
	mock.Mock
}

func (m *mockQuickNavService) BuildNavTree(campusID int64) ([]dto.CampusNavDTO, error) {
	args := m.Called(campusID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.CampusNavDTO), args.Error(1)
}

func (m *mockQuickNavService) GetLocationsByCategory(campusID int64, category string, page, pageSize int) ([]model.LocationModel, int64, error) {
	args := m.Called(campusID, category, page, pageSize)
	return args.Get(0).([]model.LocationModel), args.Get(1).(int64), args.Error(2)
}

func (m *mockQuickNavService) SearchLocations(keyword string, campusID int64, page, pageSize int) ([]model.LocationModel, int64, error) {
	args := m.Called(keyword, campusID, page, pageSize)
	return args.Get(0).([]model.LocationModel), args.Get(1).(int64), args.Error(2)
}

func (m *mockQuickNavService) GetPopularLocations(campusID int64, limit int) ([]model.LocationModel, error) {
	args := m.Called(campusID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.LocationModel), args.Error(1)
}

// setupQuickNavTest 初始化测试环境，返回 handler 和 mock service
func setupQuickNavTest(t *testing.T) (*QuickNavHandler, *mockQuickNavService) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockQuickNavService)
	handler := NewQuickNavHandler(mockSvc)
	return handler, mockSvc
}

// TestGetNavTree 测试获取导航树
func TestGetNavTree(t *testing.T) {
	handler, mockSvc := setupQuickNavTest(t)

	t.Run("成功获取（指定校区）", func(t *testing.T) {
		mockResp := []dto.CampusNavDTO{
			{
				CampusID:   1,
				CampusName: "普陀校区",
				Categories: []dto.CategoryDTO{
					{
						Category: "教学楼",
						Locations: []model.LocationModel{
							{ID: 1, Name: "第一教学楼"},
						},
						Count: 1,
					},
				},
			},
		}
		mockSvc.On("BuildNavTree", int64(1)).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/nav/tree?campus_id=1", nil)

		handler.GetNavTree(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []dto.CampusNavDTO
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "普陀校区", resp[0].CampusName)
		mockSvc.AssertExpectations(t)
	})

	t.Run("成功获取（所有校区）", func(t *testing.T) {
		mockResp := []dto.CampusNavDTO{
			{CampusID: 1, CampusName: "普陀校区"},
			{CampusID: 2, CampusName: "闵行校区"},
		}
		mockSvc.On("BuildNavTree", int64(0)).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/nav/tree?campus_id=0", nil)

		handler.GetNavTree(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []dto.CampusNavDTO
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		mockSvc.AssertExpectations(t)
	})

	t.Run("参数错误（campus_id 缺失）", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/nav/tree", nil)

		handler.GetNavTree(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("参数错误（campus_id 格式错误）", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/nav/tree?campus_id=abc", nil)

		handler.GetNavTree(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("服务层业务错误", func(t *testing.T) {
		mockSvc.On("BuildNavTree", int64(1)).Return(nil, errno.ErrServerError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/nav/tree?campus_id=1", nil)

		handler.GetNavTree(c)

		assert.Equal(t, errno.ErrServerError.Code, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层系统错误", func(t *testing.T) {
		mockSvc.On("BuildNavTree", int64(1)).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/nav/tree?campus_id=1", nil)

		handler.GetNavTree(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestGetLocationsByCategory 测试按类别获取地点
func TestGetLocationsByCategory(t *testing.T) {
	handler, mockSvc := setupQuickNavTest(t)

	t.Run("成功获取", func(t *testing.T) {
		mockLocations := []model.LocationModel{
			{ID: 1, Name: "第一教学楼"},
			{ID: 2, Name: "第二教学楼"},
		}
		mockSvc.On("GetLocationsByCategory", int64(1), "教学楼", 1, 10).Return(mockLocations, int64(2), nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/category?campus_id=1&category=教学楼&page=1&page_size=10", nil)

		handler.GetLocationsByCategory(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, float64(2), resp["total"])
		assert.Equal(t, float64(1), resp["page"])
		assert.Equal(t, float64(10), resp["size"])
		assert.Len(t, resp["list"], 2)
		mockSvc.AssertExpectations(t)
	})

	// 修正语法错误：添加括号和逗号
	t.Run("分页参数默认值", func(t *testing.T) {
		mockLocations := []model.LocationModel{{ID: 1}}
		mockSvc.On("GetLocationsByCategory", int64(1), "教学楼", 1, 20).Return(mockLocations, int64(1), nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/category?campus_id=1&category=教学楼", nil)

		handler.GetLocationsByCategory(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(1), resp["page"])
		assert.Equal(t, float64(20), resp["size"])
		mockSvc.AssertExpectations(t)
	})

	t.Run("分页参数边界", func(t *testing.T) {
		mockSvc.On("GetLocationsByCategory", int64(1), "教学楼", 1, 100).Return([]model.LocationModel{}, int64(0), nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/category?campus_id=1&category=教学楼&page=0&page_size=200", nil)

		handler.GetLocationsByCategory(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("参数错误（缺少 campus_id）", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/category?category=教学楼", nil)

		handler.GetLocationsByCategory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("参数错误（缺少 category）", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/category?campus_id=1", nil)

		handler.GetLocationsByCategory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("参数错误（campus_id 格式错误）", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/category?campus_id=abc&category=教学楼", nil)

		handler.GetLocationsByCategory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("服务层错误", func(t *testing.T) {
		mockSvc.On("GetLocationsByCategory", int64(1), "教学楼", 1, 10).Return([]model.LocationModel{}, int64(0), errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/category?campus_id=1&category=教学楼&page=1&page_size=10", nil)

		handler.GetLocationsByCategory(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestSearchLocations 测试搜索地点
func TestSearchLocations(t *testing.T) {
	handler, mockSvc := setupQuickNavTest(t)

	t.Run("成功搜索（带校区）", func(t *testing.T) {
		mockLocations := []model.LocationModel{{ID: 1, Name: "第一食堂"}}
		mockSvc.On("SearchLocations", "食堂", int64(1), 1, 10).Return(mockLocations, int64(1), nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/search?keyword=食堂&campus=1&page=1&page_size=10", nil)

		handler.SearchLocations(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, float64(1), resp["total"])
		assert.Equal(t, float64(1), resp["page"])
		assert.Equal(t, float64(10), resp["size"])
		assert.Len(t, resp["list"], 1)
		mockSvc.AssertExpectations(t)
	})

	t.Run("成功搜索（不带校区）", func(t *testing.T) {
		mockLocations := []model.LocationModel{{ID: 2, Name: "第二食堂"}}
		mockSvc.On("SearchLocations", "食堂", int64(0), 1, 20).Return(mockLocations, int64(1), nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/search?keyword=食堂", nil)

		handler.SearchLocations(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(1), resp["page"])
		assert.Equal(t, float64(20), resp["size"])
		mockSvc.AssertExpectations(t)
	})

	t.Run("关键词为空", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/search?keyword=", nil)

		handler.SearchLocations(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("campus 参数格式错误（忽略）", func(t *testing.T) {
		mockSvc.On("SearchLocations", "食堂", int64(0), 1, 20).Return([]model.LocationModel{}, int64(0), nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/search?keyword=食堂&campus=abc", nil)

		handler.SearchLocations(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层错误", func(t *testing.T) {
		mockSvc.On("SearchLocations", "食堂", int64(0), 1, 20).Return([]model.LocationModel{}, int64(0), errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/search?keyword=食堂", nil)

		handler.SearchLocations(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestGetPopularLocations 测试获取热门地点
func TestGetPopularLocations(t *testing.T) {
	handler, mockSvc := setupQuickNavTest(t)

	t.Run("成功获取", func(t *testing.T) {
		mockLocations := []model.LocationModel{
			{ID: 1, Name: "第一教学楼"},
			{ID: 2, Name: "第二教学楼"},
		}
		mockSvc.On("GetPopularLocations", int64(1), 5).Return(mockLocations, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/popular?campus_id=1&limit=5", nil)

		handler.GetPopularLocations(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []model.LocationModel
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		mockSvc.AssertExpectations(t)
	})

	// 修正语法错误：添加括号和逗号
	t.Run("limit 默认值", func(t *testing.T) {
		mockLocations := []model.LocationModel{{ID: 1}}
		mockSvc.On("GetPopularLocations", int64(1), 10).Return(mockLocations, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/popular?campus_id=1", nil)

		handler.GetPopularLocations(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("参数错误（缺少 campus_id）", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/popular", nil)

		handler.GetPopularLocations(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("campus_id 格式错误", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/popular?campus_id=abc", nil)

		handler.GetPopularLocations(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("服务层错误", func(t *testing.T) {
		mockSvc.On("GetPopularLocations", int64(1), 10).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/popular?campus_id=1", nil)

		handler.GetPopularLocations(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
