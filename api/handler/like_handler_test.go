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

// mockLikeService 是 LikeServiceInterface 的 mock 实现
type mockLikeService struct {
	mock.Mock
}

func (m *mockLikeService) ToggleLike(ctx context.Context, userID, targetID int64, targetType int8) (*dto.ToggleLikeResponse, error) {
	args := m.Called(ctx, userID, targetID, targetType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ToggleLikeResponse), args.Error(1)
}

func setupLikeTest(t *testing.T) (*LikeHandler, *mockLikeService) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockLikeService)
	handler := NewLikeHandler(mockSvc) // 假设 NewLikeHandler 接受接口类型
	return handler, mockSvc
}

func TestLikeHandler_ToggleLike(t *testing.T) {
	handler, mockSvc := setupLikeTest(t)

	t.Run("成功点赞", func(t *testing.T) {
		mockResp := &dto.ToggleLikeResponse{
			IsLiked:   true,
			LikeCount: 10,
		}
		mockSvc.On("ToggleLike", mock.Anything, int64(100), int64(1), int8(1)).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("POST", "/likes/toggle?target_id=1&target_type=1", nil)

		handler.ToggleLike(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.ToggleLikeResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.IsLiked)
		assert.Equal(t, int64(10), resp.LikeCount)
		mockSvc.AssertExpectations(t)
	})

	t.Run("成功取消点赞", func(t *testing.T) {
		mockResp := &dto.ToggleLikeResponse{
			IsLiked:   false,
			LikeCount: 9,
		}
		mockSvc.On("ToggleLike", mock.Anything, int64(100), int64(1), int8(1)).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("POST", "/likes/toggle?target_id=1&target_type=1", nil)

		handler.ToggleLike(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.ToggleLikeResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.False(t, resp.IsLiked)
		assert.Equal(t, int64(9), resp.LikeCount)
		mockSvc.AssertExpectations(t)
	})

	t.Run("未登录", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/likes/toggle?target_id=1&target_type=1", nil)

		handler.ToggleLike(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrUnauthorized.Code), resp["code"])
	})

	t.Run("缺少必要参数", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("POST", "/likes/toggle?target_type=1", nil)

		handler.ToggleLike(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("参数格式错误", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("POST", "/likes/toggle?target_id=abc&target_type=1", nil)

		handler.ToggleLike(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("服务层业务错误", func(t *testing.T) {
		mockSvc.On("ToggleLike", mock.Anything, int64(100), int64(1), int8(1)).Return(nil, errno.ErrMemoryNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("POST", "/likes/toggle?target_id=1&target_type=1", nil)

		handler.ToggleLike(c)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrMemoryNotFound.Code), resp["code"])
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层系统错误", func(t *testing.T) {
		mockSvc.On("ToggleLike", mock.Anything, int64(100), int64(1), int8(1)).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("POST", "/likes/toggle?target_id=1&target_type=1", nil)

		handler.ToggleLike(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrServerError.Code), resp["code"])
		mockSvc.AssertExpectations(t)
	})
}
