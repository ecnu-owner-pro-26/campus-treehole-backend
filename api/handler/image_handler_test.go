package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"campus-memory/application/dto"
	"campus-memory/types/errno"
)

// mockImageService 是 ImageServiceInterface 的 mock 实现
type mockImageService struct {
	mock.Mock
}

func (m *mockImageService) UploadImage(ctx context.Context, req *dto.UploadImageRequest, url string, size int64) (*dto.UploadImageResponse, error) {
	args := m.Called(ctx, req, url, size)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UploadImageResponse), args.Error(1)
}

func (m *mockImageService) DeleteImage(ctx context.Context, id int64, userID int64) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *mockImageService) GetImagesByMemoryID(ctx context.Context, memoryID int64) ([]dto.ImageInfo, error) {
	args := m.Called(ctx, memoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.ImageInfo), args.Error(1)
}

// setupImageTest 初始化测试环境，返回 handler 和 mock service
func setupImageTest(t *testing.T) (*ImageHandler, *mockImageService) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockImageService)
	handler := NewImageHandler(mockSvc)
	return handler, mockSvc
}

// createMultipartRequest 创建一个包含文件的 multipart 请求
func createMultipartRequest(t *testing.T, fieldName, filename, content string) (*http.Request, *multipart.Writer) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filename)
	assert.NoError(t, err)
	_, err = io.Copy(part, bytes.NewReader([]byte(content)))
	assert.NoError(t, err)
	err = writer.Close()
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, writer
}

// TestImageHandler_UploadImage 测试上传图片
func TestImageHandler_UploadImage(t *testing.T) {
	handler, mockSvc := setupImageTest(t)

	// 创建临时上传目录
	tempDir := t.TempDir()
	origDir, err := os.Getwd()
	assert.NoError(t, err)
	err = os.Chdir(tempDir)
	assert.NoError(t, err)
	defer os.Chdir(origDir)

	t.Run("成功上传图片", func(t *testing.T) {
		req, writer := createMultipartRequest(t, "file", "test.jpg", "fake image content")
		err := writer.WriteField("memory_id", "1")
		assert.NoError(t, err)
		writer.Close()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", int64(100))

		mockResp := &dto.UploadImageResponse{
			ID:  101,
			URL: "/uploads/images/test.jpg",
		}
		mockSvc.On("UploadImage", mock.Anything, &dto.UploadImageRequest{MemoryID: 1}, "/uploads/images/test.jpg", int64(len("fake image content"))).
			Return(mockResp, nil)

		handler.UploadImage(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp dto.UploadImageResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, int64(101), resp.ID)
		assert.Equal(t, "/uploads/images/test.jpg", resp.URL)

		_, err = os.Stat("uploads/images/test.jpg")
		assert.NoError(t, err, "文件应该被保存")

		mockSvc.AssertExpectations(t)
	})

	t.Run("未登录", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/upload", nil)

		handler.UploadImage(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrUnauthorized.Code), resp["code"])
	})

	t.Run("缺少 memory_id", func(t *testing.T) {
		req, writer := createMultipartRequest(t, "file", "test.jpg", "content")
		writer.Close()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", int64(100))

		handler.UploadImage(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// 修正：添加括号和逗号
	t.Run("memory_id 格式错误", func(t *testing.T) {
		req, writer := createMultipartRequest(t, "file", "test.jpg", "content")
		writer.WriteField("memory_id", "abc")
		writer.Close()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", int64(100))

		handler.UploadImage(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("没有上传文件", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("memory_id", "1")
		writer.Close()
		req := httptest.NewRequest("POST", "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", int64(100))

		handler.UploadImage(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("文件类型不支持", func(t *testing.T) {
		req, writer := createMultipartRequest(t, "file", "test.bmp", "content")
		writer.WriteField("memory_id", "1")
		writer.Close()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", int64(100))

		handler.UploadImage(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("文件过大", func(t *testing.T) {
		largeContent := make([]byte, 6*1024*1024)
		req, writer := createMultipartRequest(t, "file", "test.jpg", string(largeContent))
		writer.WriteField("memory_id", "1")
		writer.Close()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", int64(100))

		handler.UploadImage(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("服务层错误", func(t *testing.T) {
		req, writer := createMultipartRequest(t, "file", "test.jpg", "content")
		writer.WriteField("memory_id", "1")
		writer.Close()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", int64(100))

		mockSvc.On("UploadImage", mock.Anything, &dto.UploadImageRequest{MemoryID: 1}, "/uploads/images/test.jpg", int64(7)).
			Return(nil, errno.ErrMemoryNotFound)

		handler.UploadImage(c)

		assert.Equal(t, errno.ErrMemoryNotFound.Code, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrMemoryNotFound.Code), resp["code"])
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层系统错误", func(t *testing.T) {
		req, writer := createMultipartRequest(t, "file", "test.jpg", "content")
		writer.WriteField("memory_id", "1")
		writer.Close()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("user_id", int64(100))

		mockSvc.On("UploadImage", mock.Anything, &dto.UploadImageRequest{MemoryID: 1}, "/uploads/images/test.jpg", int64(7)).
			Return(nil, errors.New("db error"))

		handler.UploadImage(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrServerError.Code), resp["code"])
		mockSvc.AssertExpectations(t)
	})
}

// TestImageHandler_DeleteImage 测试删除图片
func TestImageHandler_DeleteImage(t *testing.T) {
	handler, mockSvc := setupImageTest(t)

	t.Run("成功删除", func(t *testing.T) {
		mockSvc.On("DeleteImage", mock.Anything, int64(1), int64(100)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("DELETE", "/images/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteImage(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(200), resp["code"])
		mockSvc.AssertExpectations(t)
	})

	t.Run("图片ID无效", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("DELETE", "/images/abc", nil)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		handler.DeleteImage(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("未登录", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("DELETE", "/images/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteImage(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("服务层业务错误", func(t *testing.T) {
		mockSvc.On("DeleteImage", mock.Anything, int64(1), int64(100)).Return(errno.ErrForbidden)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("DELETE", "/images/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteImage(c)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(errno.ErrForbidden.Code), resp["code"])
		mockSvc.AssertExpectations(t)
	})

	t.Run("服务层系统错误", func(t *testing.T) {
		mockSvc.On("DeleteImage", mock.Anything, int64(1), int64(100)).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", int64(100))
		c.Request, _ = http.NewRequest("DELETE", "/images/1", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteImage(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// TestImageHandler_GetImagesByMemoryID 测试获取记忆图片列表
func TestImageHandler_GetImagesByMemoryID(t *testing.T) {
	handler, mockSvc := setupImageTest(t)

	t.Run("成功获取", func(t *testing.T) {
		mockImages := []dto.ImageInfo{
			{ID: 1, URL: "/uploads/1.jpg"},
			{ID: 2, URL: "/uploads/2.jpg"},
		}
		mockSvc.On("GetImagesByMemoryID", mock.Anything, int64(1)).Return(mockImages, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/memories/1/images", nil)
		c.Params = gin.Params{{Key: "memory_id", Value: "1"}}

		handler.GetImagesByMemoryID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []dto.ImageInfo
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		assert.Equal(t, int64(1), resp[0].ID)
		mockSvc.AssertExpectations(t)
	})

	t.Run("记忆ID无效", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/memories/abc/images", nil)
		c.Params = gin.Params{{Key: "memory_id", Value: "abc"}}

		handler.GetImagesByMemoryID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("服务层错误", func(t *testing.T) {
		mockSvc.On("GetImagesByMemoryID", mock.Anything, int64(1)).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/memories/1/images", nil)
		c.Params = gin.Params{{Key: "memory_id", Value: "1"}}

		handler.GetImagesByMemoryID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
