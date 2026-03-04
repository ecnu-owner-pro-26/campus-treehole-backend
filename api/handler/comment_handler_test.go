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

// mockCommentService 是 CommentServiceInterface 的 mock 实现
type mockCommentService struct {
	mock.Mock
}

func (m *mockCommentService) CreateComment(ctx context.Context, req *dto.CreateCommentRequest, userID int64) (*dto.CommentResponse, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CommentResponse), args.Error(1)
}

func (m *mockCommentService) ListComments(ctx context.Context, req *dto.CommentListRequest, currentUserID *int64) (*dto.CommentListResponse, error) {
	args := m.Called(ctx, req, currentUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CommentListResponse), args.Error(1)
}

func (m *mockCommentService) ListReplies(ctx context.Context, parentID int64, page, pageSize int, currentUserID *int64) (*dto.CommentListResponse, error) {
	args := m.Called(ctx, parentID, page, pageSize, currentUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CommentListResponse), args.Error(1)
}

func (m *mockCommentService) DeleteComment(ctx context.Context, id int64, userID int64) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

// setupCommentTest 初始化测试环境，返回 handler 和 mock service
func setupCommentTest(t *testing.T) (*CommentHandler, *mockCommentService) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockCommentService)
	handler := NewCommentHandler(mockSvc)
	return handler, mockSvc
}

// TestCreateComment 测试创建评论
func TestCreateComment(t *testing.T) {
	handler, mockSvc := setupCommentTest(t)

	t.Run("成功创建评论", func(t *testing.T) {
		reqBody := dto.CreateCommentRequest{
			MemoryID: 1,
			Content:  "这是一条评论",
			ParentID: nil,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		mockResp := &dto.CommentResponse{
			ID:      1,
			Content: "这是一条评论",
			User: dto.UserSimpleInfo{
				ID:       100,
				Nickname: "测试用户",
			},
		}
		mockSvc.On("CreateComment", mock.Anything, &reqBody, int64(100)).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("POST", "/api/comments", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateComment(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CommentResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), resp.ID)
		assert.Equal(t, "这是一条评论", resp.Content)
		mockSvc.AssertExpectations(t)
	})

	t.Run("未登录", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// 不设置 user_id
		c.Request, _ = http.NewRequest("POST", "/api/comments", nil)

		handler.CreateComment(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrUnauthorized.Code), resp["code"])
	})

	t.Run("请求参数错误", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("POST", "/api/comments", bytes.NewReader([]byte("{invalid}")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateComment(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrBadRequest.Code), resp["code"])
	})

	t.Run("服务层返回业务错误", func(t *testing.T) {
		reqBody := dto.CreateCommentRequest{MemoryID: 1, Content: "test"}
		bodyBytes, _ := json.Marshal(reqBody)

		mockSvc.On("CreateComment", mock.Anything, &reqBody, int64(100)).Return(nil, errno.ErrMemoryNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("POST", "/api/comments", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateComment(c)

		assert.Equal(t, errno.ErrMemoryNotFound.Code, w.Code) // 注意：ErrorResponse 会设置 HTTP 状态码为 200，业务码在 body 中
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrMemoryNotFound.Code), resp["code"])
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层返回系统错误", func(t *testing.T) {
		reqBody := dto.CreateCommentRequest{MemoryID: 1, Content: "test"}
		bodyBytes, _ := json.Marshal(reqBody)

		mockSvc.On("CreateComment", mock.Anything, &reqBody, int64(100)).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("POST", "/api/comments", bytes.NewReader(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateComment(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrServerError.Code), resp["code"])
		mockSvc.AssertExpectations(t)
	})
}

// TestListComments 测试获取评论列表
func TestListComments(t *testing.T) {
	handler, mockSvc := setupCommentTest(t)

	t.Run("成功获取列表（未登录）", func(t *testing.T) {
		req := dto.CommentListRequest{
			MemoryID: 1,
			Page:     1,
			PageSize: 20,
			SortBy:   "latest",
		}
		mockResp := &dto.CommentListResponse{
			Comments: []*dto.CommentResponse{
				{ID: 1, Content: "评论1"},
				{ID: 2, Content: "评论2"},
			},
			Total:    2,
			Page:     1,
			PageSize: 20,
		}
		mockSvc.On("ListComments", mock.Anything, &req, (*int64)(nil)).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/comments?memory_id=1&page=1&page_size=20", nil)

		handler.ListComments(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CommentListResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(resp.Comments))
		mockSvc.AssertExpectations(t)
	})

	t.Run("成功获取列表（已登录）", func(t *testing.T) {
		req := dto.CommentListRequest{
			MemoryID: 1,
			Page:     1,
			PageSize: 20,
			SortBy:   "latest",
		}
		currentUserID := int64(100)
		mockResp := &dto.CommentListResponse{
			Comments: []*dto.CommentResponse{
				{ID: 1, Content: "评论1", IsLiked: true},
			},
			Total: 1,
		}
		mockSvc.On("ListComments", mock.Anything, &req, &currentUserID).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", currentUserID)
		c.Request, _ = http.NewRequest("GET", "/api/comments?memory_id=1&page=1&page_size=20", nil)

		handler.ListComments(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CommentListResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Comments[0].IsLiked)
		mockSvc.AssertExpectations(t)
	})

	t.Run("查询参数错误", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/comments?memory_id=abc", nil) // memory_id 应为数字

		handler.ListComments(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("服务层错误", func(t *testing.T) {
		req := dto.CommentListRequest{
			MemoryID: 1,
			Page:     1,
			PageSize: 20,
			SortBy:   "latest",
		}
		mockSvc.On("ListComments", mock.Anything, &req, (*int64)(nil)).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/comments?memory_id=1", nil)

		handler.ListComments(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestListReplies 测试获取回复列表
func TestListReplies(t *testing.T) {
	handler, mockSvc := setupCommentTest(t)

	t.Run("成功获取回复列表", func(t *testing.T) {
		parentID := int64(10)
		page := 1
		pageSize := 20
		currentUserID := int64(100)

		mockResp := &dto.CommentListResponse{
			Comments: []*dto.CommentResponse{
				{ID: 101, Content: "回复1"},
				{ID: 102, Content: "回复2"},
			},
			Total: 2,
		}
		mockSvc.On("ListReplies", mock.Anything, parentID, page, pageSize, &currentUserID).Return(mockResp, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", currentUserID)
		c.Request, _ = http.NewRequest("GET", "/api/comments/10/replies?page=1&page_size=20", nil)
		c.Params = gin.Params{{Key: "parent_id", Value: "10"}}

		handler.ListReplies(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.CommentListResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(resp.Comments))
		mockSvc.AssertExpectations(t)
	})

	t.Run("路径参数错误", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/comments/abc/replies", nil)
		c.Params = gin.Params{{Key: "parent_id", Value: "abc"}}

		handler.ListReplies(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("分页参数错误", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/comments/10/replies?page=abc", nil)
		c.Params = gin.Params{{Key: "parent_id", Value: "10"}}

		handler.ListReplies(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestDeleteComment 测试删除评论
func TestDeleteComment(t *testing.T) {
	handler, mockSvc := setupCommentTest(t)

	t.Run("成功删除", func(t *testing.T) {
		mockSvc.On("DeleteComment", mock.Anything, int64(1), int64(100)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("DELETE", "/api/comments/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteComment(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(200), resp["code"]) // 假设成功响应 code 为 200
		mockSvc.AssertExpectations(t)
	})

	t.Run("评论ID无效", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("DELETE", "/api/comments/abc", nil)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		handler.DeleteComment(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("未登录", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// 不设置 user_id
		c.Request, _ = http.NewRequest("DELETE", "/api/comments/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteComment(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("无权删除", func(t *testing.T) {
		mockSvc.On("DeleteComment", mock.Anything, int64(1), int64(100)).Return(errno.ErrForbidden)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("DELETE", "/api/comments/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteComment(c)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrForbidden.Code), resp["code"])
		mockSvc.AssertExpectations(t)
	})
}
