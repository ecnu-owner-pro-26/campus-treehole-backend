package handler

import (
	"bytes"
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

type mockLocationService struct {
	mock.Mock
}

func (m *mockLocationService) GetLocation(ctx context.Context, id int64) (*dto.LocationResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LocationResponse), args.Error(1)
}

func (m *mockLocationService) ListLocations(ctx context.Context, req *dto.LocationListRequest) (*dto.LocationListResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LocationListResponse), args.Error(1)
}

func (m *mockLocationService) CreateLocation(ctx context.Context, req *dto.CreateLocationRequest) (*dto.LocationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LocationResponse), args.Error(1)
}

func (m *mockLocationService) UpdateLocation(ctx context.Context, id int64, req *dto.UpdateLocationRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *mockLocationService) DeleteLocation(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockLocationService) SearchLocations(ctx context.Context, keyword string, campusID *int64) ([]dto.LocationResponse, error) {
	args := m.Called(ctx, keyword, campusID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.LocationResponse), args.Error(1)
}

// setupLocationTest 初始化测试环境，返回 handler 和 mock service
func setupLocationTest(t *testing.T) (*LocationHandler, *mockLocationService) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockLocationService)
	handler := NewLocationHandler(mockSvc) // 假设 NewLocationHandler 接受接口类型
	return handler, mockSvc
}

// TestLocationHandler_GetLocation 测试获取地点详情
func TestLocationHandler_GetLocation(t *testing.T) {
	handler, mockSvc := setupLocationTest(t)

	t.Run("成功获取", func(t *testing.T) {
		mockResp := &dto.LocationResponse{
			ID:       1,
			Name:     "教学楼A",
			CampusID: 1,
		}
		mockSvc.On("GetLocation", mock.Anything, int64(1)).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.GetLocation(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.LocationResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "教学楼A", resp.Name)
		mockSvc.AssertExpectations(t)
	})

	t.Run("地点ID无效", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/abc", nil)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		handler.GetLocation(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("地点不存在", func(t *testing.T) {
		mockSvc.On("GetLocation", mock.Anything, int64(999)).Return(nil, errno.ErrLocationNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/999", nil)
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		handler.GetLocation(c)

		assert.Equal(t, errno.ErrLocationNotFound.Code, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层系统错误", func(t *testing.T) {
		mockSvc.On("GetLocation", mock.Anything, int64(1)).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.GetLocation(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestLocationHandler_ListLocations 测试获取地点列表
func TestLocationHandler_ListLocations(t *testing.T) {
	handler, mockSvc := setupLocationTest(t)

	t.Run("成功获取列表", func(t *testing.T) {
		req := dto.LocationListRequest{
			CampusID: int64Ptr(1),
			Category: "教学楼",
			Page:     1,
			PageSize: 10,
		}
		mockResp := &dto.LocationListResponse{
			Locations: []dto.LocationResponse{
				{ID: 1, Name: "教学楼A"},
				{ID: 2, Name: "教学楼B"},
			},
			Total:    2,
			Page:     1,
			PageSize: 10,
		}
		mockSvc.On("ListLocations", mock.Anything, &req).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations?campus_id=1&category=教学楼&page=1&page_size=10", nil)

		handler.ListLocations(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.LocationListResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp.Locations, 2)
		mockSvc.AssertExpectations(t)
	})

	t.Run("参数错误", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations?page=abc", nil)

		handler.ListLocations(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("服务层错误", func(t *testing.T) {
		req := dto.LocationListRequest{Page: 1, PageSize: 20}
		mockSvc.On("ListLocations", mock.Anything, &req).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations", nil)

		handler.ListLocations(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestLocationHandler_CreateLocation 测试创建地点
func TestLocationHandler_CreateLocation(t *testing.T) {
	handler, mockSvc := setupLocationTest(t)

	t.Run("成功创建", func(t *testing.T) {
		reqBody := dto.CreateLocationRequest{
			Name:      "新食堂",
			CampusID:  1,
			Category:  "食堂",
			SortOrder: 1,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		mockResp := &dto.LocationResponse{
			ID:       101,
			Name:     "新食堂",
			CampusID: 1,
			Category: "食堂",
		}
		mockSvc.On("CreateLocation", mock.Anything, &reqBody).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/locations", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateLocation(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.LocationResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, int64(101), resp.ID)
		mockSvc.AssertExpectations(t)
	})

	t.Run("请求参数错误", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/locations", bytes.NewReader([]byte("{invalid}")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateLocation(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("服务层业务错误", func(t *testing.T) {
		reqBody := dto.CreateLocationRequest{Name: "地点", CampusID: 1}
		bodyBytes, _ := json.Marshal(reqBody)

		mockSvc.On("CreateLocation", mock.Anything, &reqBody).Return(nil, errno.ErrCampusNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/locations", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateLocation(c)

		assert.Equal(t, errno.ErrCampusNotFound.Code, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层系统错误", func(t *testing.T) {
		reqBody := dto.CreateLocationRequest{Name: "地点", CampusID: 1}
		bodyBytes, _ := json.Marshal(reqBody)

		mockSvc.On("CreateLocation", mock.Anything, &reqBody).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/locations", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateLocation(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestLocationHandler_UpdateLocation 测试更新地点
func TestLocationHandler_UpdateLocation(t *testing.T) {
	handler, mockSvc := setupLocationTest(t)

	t.Run("成功更新", func(t *testing.T) {
		reqBody := dto.UpdateLocationRequest{
			Name:     strPtr("更新后名称"),
			Category: strPtr("教学楼"),
		}
		bodyBytes, _ := json.Marshal(reqBody)

		mockSvc.On("UpdateLocation", mock.Anything, int64(1), &reqBody).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("PUT", "/locations/1", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateLocation(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("地点ID无效", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Request, _ = http.NewRequest("PUT", "/locations/abc", nil)

		handler.UpdateLocation(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("请求参数错误", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("PUT", "/locations/1", bytes.NewReader([]byte("{invalid}")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateLocation(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("地点不存在", func(t *testing.T) {
		reqBody := dto.UpdateLocationRequest{Name: strPtr("新名称")}
		bodyBytes, _ := json.Marshal(reqBody)

		mockSvc.On("UpdateLocation", mock.Anything, int64(999), &reqBody).Return(errno.ErrLocationNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request, _ = http.NewRequest("PUT", "/locations/999", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateLocation(c)

		assert.Equal(t, errno.ErrLocationNotFound.Code, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层系统错误", func(t *testing.T) {
		reqBody := dto.UpdateLocationRequest{Name: strPtr("新名称")}
		bodyBytes, _ := json.Marshal(reqBody)

		mockSvc.On("UpdateLocation", mock.Anything, int64(1), &reqBody).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("PUT", "/locations/1", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateLocation(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestLocationHandler_DeleteLocation 测试删除地点
func TestLocationHandler_DeleteLocation(t *testing.T) {
	handler, mockSvc := setupLocationTest(t)

	t.Run("成功删除", func(t *testing.T) {
		mockSvc.On("DeleteLocation", mock.Anything, int64(1)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("DELETE", "/locations/1", nil)

		handler.DeleteLocation(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("地点ID无效", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Request, _ = http.NewRequest("DELETE", "/locations/abc", nil)

		handler.DeleteLocation(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("地点不存在", func(t *testing.T) {
		mockSvc.On("DeleteLocation", mock.Anything, int64(999)).Return(errno.ErrLocationNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request, _ = http.NewRequest("DELETE", "/locations/999", nil)

		handler.DeleteLocation(c)

		assert.Equal(t, errno.ErrLocationNotFound.Code, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层系统错误", func(t *testing.T) {
		mockSvc.On("DeleteLocation", mock.Anything, int64(1)).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("DELETE", "/locations/1", nil)

		handler.DeleteLocation(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestLocationHandler_SearchLocations 测试搜索地点
func TestLocationHandler_SearchLocations(t *testing.T) {
	handler, mockSvc := setupLocationTest(t)

	// 修正语法错误：添加括号和逗号
	t.Run("成功搜索（不带校区）", func(t *testing.T) {
		keyword := "食堂"
		mockResp := []dto.LocationResponse{
			{ID: 1, Name: "第一食堂"},
			{ID: 2, Name: "第二食堂"},
		}
		mockSvc.On("SearchLocations", mock.Anything, keyword, (*int64)(nil)).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/search?keyword=食堂", nil)

		handler.SearchLocations(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []dto.LocationResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		mockSvc.AssertExpectations(t)
	})

	t.Run("成功搜索（带校区）", func(t *testing.T) {
		keyword := "食堂"
		campusID := int64(1)
		mockResp := []dto.LocationResponse{
			{ID: 1, Name: "第一食堂"},
		}
		mockSvc.On("SearchLocations", mock.Anything, keyword, &campusID).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/search?keyword=食堂&campus_id=1", nil)

		handler.SearchLocations(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []dto.LocationResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		mockSvc.AssertExpectations(t)
	})

	t.Run("关键词为空", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/search?keyword=", nil)

		handler.SearchLocations(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("校区ID格式错误（忽略）", func(t *testing.T) {
		keyword := "食堂"
		mockResp := []dto.LocationResponse{{ID: 1, Name: "第一食堂"}}
		mockSvc.On("SearchLocations", mock.Anything, keyword, (*int64)(nil)).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/search?keyword=食堂&campus_id=abc", nil)

		handler.SearchLocations(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层错误", func(t *testing.T) {
		mockSvc.On("SearchLocations", mock.Anything, "食堂", (*int64)(nil)).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/locations/search?keyword=食堂", nil)

		handler.SearchLocations(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
