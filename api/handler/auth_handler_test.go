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
)

// mockAuthService 是 AuthService 的 mock 实现
type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) WechatLogin(ctx context.Context, req *dto.WechatLoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func (m *mockAuthService) GetUserProfile(ctx context.Context, userID int64) (*dto.UserProfileResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserProfileResponse), args.Error(1)
}

func (m *mockAuthService) UpdateUserProfile(ctx context.Context, userID int64, req *dto.UpdateProfileRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

// setupAuthTest 初始化测试环境，返回 router 和 mock service
func setupAuthTest(t *testing.T) (*gin.Engine, *mockAuthService) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockSvc := new(mockAuthService)
	handler := NewAuthHandler(mockSvc)

	// 注册路由
	r.POST("/api/auth/login", handler.WechatLogin)
	r.GET("/api/auth/profile", handler.GetProfile)
	r.PUT("/api/auth/profile", handler.UpdateProfile)

	return r, mockSvc
}

// TestWechatLogin 测试微信登录
func TestWechatLogin(t *testing.T) {
	r, mockSvc := setupAuthTest(t)

	t.Run("成功登录", func(t *testing.T) {
		reqBody := dto.WechatLoginRequest{
			Code:     "test_code",
			Nickname: "测试用户",
			Avatar:   "avatar.jpg",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		// mock 服务层返回
		mockResp := &dto.LoginResponse{
			Token: "test-token",
			User: dto.UserInfoDTO{
				ID:       1,
				Nickname: "测试用户",
				Avatar:   "avatar.jpg",
			},
		}
		mockSvc.On("WechatLogin", mock.Anything, &reqBody).Return(mockResp, nil)

		// 发起请求
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// 断言
		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "test-token", resp.Token)
		assert.Equal(t, int64(1), resp.User.ID)
		mockSvc.AssertExpectations(t)
	})

	t.Run("请求参数错误", func(t *testing.T) {
		// 发送无效 JSON
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewReader([]byte("{invalid}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, float64(400), resp["code"])
	})

	t.Run("服务层错误", func(t *testing.T) {
		reqBody := dto.WechatLoginRequest{Code: "test_code"}
		bodyBytes, _ := json.Marshal(reqBody)

		mockSvc.On("WechatLogin", mock.Anything, &reqBody).Return(nil, errors.New("微信接口错误"))

		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestGetProfile 测试获取个人信息
func TestGetProfile(t *testing.T) {
	_, mockSvc := setupAuthTest(t)

	t.Run("成功获取", func(t *testing.T) {
		// 模拟中间件设置 user_id
		req, _ := http.NewRequest("GET", "/api/auth/profile", nil)
		w := httptest.NewRecorder()

		// 构造带有上下文的请求（手动设置 user_id）
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("user_id", int64(1))

		mockResp := &dto.UserProfileResponse{
			ID:       1,
			Nickname: "测试用户",
			Avatar:   "avatar.jpg",
			Role:     0,
			Status:   1,
		}
		mockSvc.On("GetUserProfile", mock.Anything, int64(1)).Return(mockResp, nil)

		// 直接调用 handler（不走路由）
		handler := NewAuthHandler(mockSvc)
		handler.GetProfile(ctx)

		// 断言
		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.UserProfileResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "测试用户", resp.Nickname)
		mockSvc.AssertExpectations(t)
	})

	t.Run("未登录", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/auth/profile", nil)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		// 不设置 user_id

		handler := NewAuthHandler(mockSvc)
		handler.GetProfile(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(401), resp["code"])
	})

	t.Run("服务层错误", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/auth/profile", nil)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("user_id", int64(1))

		mockSvc.On("GetUserProfile", mock.Anything, int64(1)).Return(nil, errors.New("db error"))

		handler := NewAuthHandler(mockSvc)
		handler.GetProfile(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestUpdateProfile 测试更新个人信息
func TestUpdateProfile(t *testing.T) {
	_, mockSvc := setupAuthTest(t)

	t.Run("成功更新", func(t *testing.T) {
		reqBody := dto.UpdateProfileRequest{
			Nickname: strPtr("新昵称"),
			Avatar:   strPtr("new.jpg"),
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("PUT", "/api/auth/profile", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("user_id", int64(1))

		mockSvc.On("UpdateUserProfile", mock.Anything, int64(1), &reqBody).Return(nil)

		handler := NewAuthHandler(mockSvc)
		handler.UpdateProfile(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "更新成功", resp["data"].(map[string]interface{})["message"])
		mockSvc.AssertExpectations(t)
	})

	t.Run("请求参数错误", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/api/auth/profile", bytes.NewReader([]byte("{invalid}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("user_id", int64(1))

		handler := NewAuthHandler(mockSvc)
		handler.UpdateProfile(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("未登录", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/api/auth/profile", nil)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		// 不设置 user_id

		handler := NewAuthHandler(mockSvc)
		handler.UpdateProfile(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("服务层错误", func(t *testing.T) {
		reqBody := dto.UpdateProfileRequest{Nickname: strPtr("新昵称")}
		bodyBytes, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("PUT", "/api/auth/profile", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("user_id", int64(1))

		mockSvc.On("UpdateUserProfile", mock.Anything, int64(1), &reqBody).Return(errors.New("update error"))

		handler := NewAuthHandler(mockSvc)
		handler.UpdateProfile(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// 辅助函数：返回字符串指针
func strPtr(s string) *string {
	return &s
}
