package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"campus-memory/application/dto"
	"campus-memory/types/errno"
)

// mockCampusService 是 CampusServiceInterface 的 mock 实现
type mockCampusService struct {
	mock.Mock
}

func (m *mockCampusService) ListCampuses(ctx context.Context) (*dto.CampusListResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CampusListResponse), args.Error(1)
}

func (m *mockCampusService) GetCampus(ctx context.Context, id int64) (*dto.CampusResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CampusResponse), args.Error(1)
}

func (m *mockCampusService) GetCampusWithLocations(ctx context.Context, campusID int64) (*dto.CampusLocationsResponse, error) {
	args := m.Called(ctx, campusID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CampusLocationsResponse), args.Error(1)
}

// setupCampusTest 初始化测试环境，返回 handler 和 mock service
func setupCampusTest(t *testing.T) (*CampusHandler, *mockCampusService) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockCampusService)
	handler := NewCampusHandler(mockSvc)
	return handler, mockSvc
}

// TestListCampuses 测试获取校区列表
func TestListCampuses(t *testing.T) {
	handler, mockSvc := setupCampusTest(t)

	t.Run("成功获取列表", func(t *testing.T) {
		// mock 服务层返回
		mockResp := &dto.CampusListResponse{
			Campuses: []dto.CampusResponse{
				{ID: 1, Name: "校区A"},
				{ID: 2, Name: "校区B"},
			},
			Total: 2,
		}
		mockSvc.On("ListCampuses", mock.Anything).Return(mockResp, nil)

		// 创建请求
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest("GET", "/api/campuses", nil)
		c.Request = req

		// 调用 handler
		handler.ListCampuses(c)

		// 断言
		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CampusListResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(resp.Campuses))
		assert.Equal(t, int64(2), resp.Total)
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层错误", func(t *testing.T) {
		mockSvc.On("ListCampuses", mock.Anything).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest("GET", "/api/campuses", nil)
		c.Request = req

		handler.ListCampuses(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrServerError.Code), resp["code"])
		mockSvc.AssertExpectations(t)
	})
}

// TestGetCampus 测试获取校区详情
func TestGetCampus(t *testing.T) {
	handler, mockSvc := setupCampusTest(t)

	t.Run("成功获取", func(t *testing.T) {
		mockResp := &dto.CampusResponse{
			ID:   1,
			Name: "普陀校区",
		}
		mockSvc.On("GetCampus", mock.Anything, int64(1)).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest("GET", "/api/campuses/1", nil)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.GetCampus(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CampusResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "普陀校区", resp.Name)
		mockSvc.AssertExpectations(t)
	})

	t.Run("无效的校区ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest("GET", "/api/campuses/abc", nil)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		handler.GetCampus(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrBadRequest.Code), resp["code"])
	})

	t.Run("校区不存在", func(t *testing.T) {
		mockSvc.On("GetCampus", mock.Anything, int64(999)).Return(nil, errno.ErrCampusNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest("GET", "/api/campuses/999", nil)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		handler.GetCampus(c)

		assert.Equal(t, errno.ErrCampusNotFound.Code, w.Code) // 注意：ErrorResponse 会设置 HTTP 状态码为 200，但返回的 code 字段是业务码，需要检查 body
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrCampusNotFound.Code), resp["code"])
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层其他错误", func(t *testing.T) {
		mockSvc.On("GetCampus", mock.Anything, int64(1)).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest("GET", "/api/campuses/1", nil)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.GetCampus(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrServerError.Code), resp["code"])
		mockSvc.AssertExpectations(t)
	})
}

// TestGetCampusWithLocations 测试获取校区及地点列表
func TestGetCampusWithLocations(t *testing.T) {
	handler, mockSvc := setupCampusTest(t)

	t.Run("成功获取", func(t *testing.T) {
		mockResp := &dto.CampusLocationsResponse{
			Campus: dto.CampusResponse{ID: 1, Name: "普陀校区"},
			Locations: []dto.LocationResponse{
				{ID: 101, Name: "教学楼A"},
				{ID: 102, Name: "食堂B"},
			},
		}
		mockSvc.On("GetCampusWithLocations", mock.Anything, int64(1)).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest("GET", "/api/campuses/1/locations", nil)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.GetCampusWithLocations(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CampusLocationsResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "普陀校区", resp.Campus.Name)
		assert.Equal(t, 2, len(resp.Locations))
		mockSvc.AssertExpectations(t)
	})

	t.Run("无效的校区ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest("GET", "/api/campuses/abc/locations", nil)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		handler.GetCampusWithLocations(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("校区不存在", func(t *testing.T) {
		mockSvc.On("GetCampusWithLocations", mock.Anything, int64(999)).Return(nil, errno.ErrCampusNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest("GET", "/api/campuses/999/locations", nil)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		handler.GetCampusWithLocations(c)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrCampusNotFound.Code), resp["code"])
		mockSvc.AssertExpectations(t)
	})
}
