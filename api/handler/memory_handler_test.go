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

// mockMemoryService 是 MemoryServiceInterface 的 mock 实现
type mockMemoryService struct {
	mock.Mock
}

func (m *mockMemoryService) CreateMemory(ctx context.Context, req *dto.CreateMemoryRequest, creatorID int64) (*dto.MemoryResponse, error) {
	args := m.Called(ctx, req, creatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MemoryResponse), args.Error(1)
}

func (m *mockMemoryService) GetMemory(ctx context.Context, id int64, currentUserID *int64) (*dto.MemoryResponse, error) {
	args := m.Called(ctx, id, currentUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MemoryResponse), args.Error(1)
}

func (m *mockMemoryService) ListMemories(ctx context.Context, req *dto.MemoryListRequest, currentUserID *int64) (*dto.MemoryListResponse, error) {
	args := m.Called(ctx, req, currentUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MemoryListResponse), args.Error(1)
}

func (m *mockMemoryService) UpdateMemory(ctx context.Context, id int64, req *dto.UpdateMemoryRequest, userID int64) error {
	args := m.Called(ctx, id, req, userID)
	return args.Error(0)
}

func (m *mockMemoryService) DeleteMemory(ctx context.Context, id int64, userID int64) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

// setupMemoryTest 初始化测试环境，返回 handler 和 mock service
func setupMemoryTest(t *testing.T) (*MemoryHandler, *mockMemoryService) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockMemoryService)
	handler := NewMemoryHandler(mockSvc) // 假设 NewMemoryHandler 接受接口类型
	return handler, mockSvc
}

// TestMemoryHandler_CreateMemory 测试创建记忆
func TestMemoryHandler_CreateMemory(t *testing.T) {
	handler, mockSvc := setupMemoryTest(t)

	t.Run("成功创建", func(t *testing.T) {
		reqBody := dto.CreateMemoryRequest{
			Title:      "测试记忆",
			Content:    "这是一条测试记忆",
			LocationID: int64Ptr(1),
			Tags:       []string{"tag1", "tag2"},
			ImageURLs:  []string{"http://example.com/1.jpg"},
		}
		bodyBytes, _ := json.Marshal(reqBody)

		mockResp := &dto.MemoryResponse{
			ID:      1,
			Title:   "测试记忆",
			Content: "这是一条测试记忆",
			Creator: dto.UserSimpleInfo{
				ID:       100,
				Nickname: "测试用户",
			},
		}
		mockSvc.On("CreateMemory", mock.Anything, &reqBody, int64(100)).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("POST", "/memories", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateMemory(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.MemoryResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), resp.ID)
		mockSvc.AssertExpectations(t)
	})

	t.Run("未登录", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/memories", nil)

		handler.CreateMemory(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("请求参数错误", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("POST", "/memories", bytes.NewReader([]byte("{invalid}")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateMemory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("服务层业务错误", func(t *testing.T) {
		reqBody := dto.CreateMemoryRequest{Title: "测试", LocationID: int64Ptr(1)}
		bodyBytes, _ := json.Marshal(reqBody)

		mockSvc.On("CreateMemory", mock.Anything, &reqBody, int64(100)).Return(nil, errno.ErrLocationNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("POST", "/memories", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateMemory(c)

		assert.Equal(t, errno.ErrLocationNotFound.Code, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层系统错误", func(t *testing.T) {
		reqBody := dto.CreateMemoryRequest{Title: "测试", LocationID: int64Ptr(1)}
		bodyBytes, _ := json.Marshal(reqBody)

		mockSvc.On("CreateMemory", mock.Anything, &reqBody, int64(100)).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("POST", "/memories", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateMemory(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestMemoryHandler_GetMemory 测试获取记忆详情
func TestMemoryHandler_GetMemory(t *testing.T) {
	handler, mockSvc := setupMemoryTest(t)

	t.Run("成功获取（未登录）", func(t *testing.T) {
		mockResp := &dto.MemoryResponse{
			ID:      1,
			Title:   "测试记忆",
			Content: "内容",
		}
		mockSvc.On("GetMemory", mock.Anything, int64(1), (*int64)(nil)).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/memories/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.GetMemory(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.MemoryResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), resp.ID)
		mockSvc.AssertExpectations(t)
	})

	t.Run("成功获取（已登录）", func(t *testing.T) {
		currentUserID := int64(100)
		mockResp := &dto.MemoryResponse{
			ID:      1,
			Title:   "测试记忆",
			IsLiked: true,
		}
		mockSvc.On("GetMemory", mock.Anything, int64(1), &currentUserID).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", currentUserID)
		c.Request, _ = http.NewRequest("GET", "/memories/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.GetMemory(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.MemoryResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.IsLiked)
		mockSvc.AssertExpectations(t)
	})

	t.Run("记忆ID无效", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/memories/abc", nil)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		handler.GetMemory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("记忆不存在", func(t *testing.T) {
		mockSvc.On("GetMemory", mock.Anything, int64(999), (*int64)(nil)).Return(nil, errno.ErrMemoryNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/memories/999", nil)
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		handler.GetMemory(c)

		assert.Equal(t, errno.ErrMemoryNotFound.Code, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层系统错误", func(t *testing.T) {
		mockSvc.On("GetMemory", mock.Anything, int64(1), (*int64)(nil)).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/memories/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.GetMemory(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestMemoryHandler_ListMemories 测试获取记忆列表
func TestMemoryHandler_ListMemories(t *testing.T) {
	handler, mockSvc := setupMemoryTest(t)

	t.Run("成功获取列表（未登录）", func(t *testing.T) {
		req := dto.MemoryListRequest{
			LocationID: int64Ptr(1),
			Page:       1,
			PageSize:   10,
			SortBy:     "latest",
		}
		mockResp := &dto.MemoryListResponse{
			Memories: []dto.MemoryResponse{
				{ID: 1, Title: "记忆1"},
				{ID: 2, Title: "记忆2"},
			},
			Total:    2,
			Page:     1,
			PageSize: 10,
		}
		mockSvc.On("ListMemories", mock.Anything, &req, (*int64)(nil)).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/memories?location_id=1&page=1&page_size=10&sort_by=latest", nil)

		handler.ListMemories(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.MemoryListResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp.Memories, 2)
		mockSvc.AssertExpectations(t)
	})

	t.Run("成功获取列表（已登录）", func(t *testing.T) {
		req := dto.MemoryListRequest{Page: 1, PageSize: 20}
		currentUserID := int64(100)
		mockResp := &dto.MemoryListResponse{
			Memories: []dto.MemoryResponse{
				{ID: 1, Title: "记忆1", IsLiked: true},
			},
			Total: 1,
		}
		mockSvc.On("ListMemories", mock.Anything, &req, &currentUserID).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", currentUserID)
		c.Request, _ = http.NewRequest("GET", "/memories?page=1&page_size=20", nil)

		handler.ListMemories(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.MemoryListResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Memories[0].IsLiked)
		mockSvc.AssertExpectations(t)
	})

	t.Run("参数错误", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/memories?page=abc", nil)

		handler.ListMemories(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("服务层错误", func(t *testing.T) {
		req := dto.MemoryListRequest{Page: 1, PageSize: 20}
		mockSvc.On("ListMemories", mock.Anything, &req, (*int64)(nil)).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/memories?page=1&page_size=20", nil)

		handler.ListMemories(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestMemoryHandler_UpdateMemory 测试更新记忆
func TestMemoryHandler_UpdateMemory(t *testing.T) {
	handler, mockSvc := setupMemoryTest(t)

	t.Run("成功更新", func(t *testing.T) {
		reqBody := dto.UpdateMemoryRequest{
			Title:     strPtr("新标题"),
			Content:   strPtr("新内容"),
			ImageURLs: []string{"new1.jpg", "new2.jpg"},
		}
		bodyBytes, _ := json.Marshal(reqBody)

		mockSvc.On("UpdateMemory", mock.Anything, int64(1), &reqBody, int64(100)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("PUT", "/memories/1", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateMemory(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("记忆ID无效", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Request, _ = http.NewRequest("PUT", "/memories/abc", nil)

		handler.UpdateMemory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("未登录", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("PUT", "/memories/1", nil)

		handler.UpdateMemory(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("请求参数错误", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("PUT", "/memories/1", bytes.NewReader([]byte("{invalid}")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateMemory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("无权更新", func(t *testing.T) {
		reqBody := dto.UpdateMemoryRequest{Title: strPtr("新标题")}
		bodyBytes, _ := json.Marshal(reqBody)

		mockSvc.On("UpdateMemory", mock.Anything, int64(1), &reqBody, int64(100)).Return(errno.ErrForbidden)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("PUT", "/memories/1", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateMemory(c)

		assert.Equal(t, errno.ErrForbidden.Code, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层系统错误", func(t *testing.T) {
		reqBody := dto.UpdateMemoryRequest{Title: strPtr("新标题")}
		bodyBytes, _ := json.Marshal(reqBody)

		mockSvc.On("UpdateMemory", mock.Anything, int64(1), &reqBody, int64(100)).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("PUT", "/memories/1", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateMemory(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestMemoryHandler_DeleteMemory 测试删除记忆
func TestMemoryHandler_DeleteMemory(t *testing.T) {
	handler, mockSvc := setupMemoryTest(t)

	t.Run("成功删除", func(t *testing.T) {
		mockSvc.On("DeleteMemory", mock.Anything, int64(1), int64(100)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("DELETE", "/memories/1", nil)

		handler.DeleteMemory(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("记忆ID无效", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Request, _ = http.NewRequest("DELETE", "/memories/abc", nil)

		handler.DeleteMemory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("未登录", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("DELETE", "/memories/1", nil)

		handler.DeleteMemory(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("无权删除", func(t *testing.T) {
		mockSvc.On("DeleteMemory", mock.Anything, int64(1), int64(100)).Return(errno.ErrForbidden)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("DELETE", "/memories/1", nil)

		handler.DeleteMemory(c)

		assert.Equal(t, errno.ErrForbidden.Code, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层系统错误", func(t *testing.T) {
		mockSvc.On("DeleteMemory", mock.Anything, int64(1), int64(100)).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("DELETE", "/memories/1", nil)

		handler.DeleteMemory(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// 辅助函数：返回 *int64
func int64Ptr(i int64) *int64 {
	return &i
}
