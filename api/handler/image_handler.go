package handler

import (
	"campus-memory/application/dto"
	"campus-memory/application/service"
	"campus-memory/infra/util"
	"campus-memory/types/errno"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ImageHandler 图片上传处理器
type ImageHandler struct {
	imageService *service.ImageService
}

// NewImageHandler 创建图片处理器实例
func NewImageHandler(imageService *service.ImageService) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
	}
}

// UploadImage 上传图片
func (h *ImageHandler) UploadImage(c *gin.Context) {
	// 1. 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, errno.ErrUnauthorized.Code, errno.ErrUnauthorized.Message)
		return
	}
	_ = userID // 如果需要权限验证可以使用

	// 2. 获取记忆ID
	memoryIDStr := c.PostForm("memory_id")
	memoryID, err := strconv.ParseInt(memoryIDStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的记忆ID")
		return
	}

	// 3. 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "请选择要上传的图片")
		return
	}

	// 4. 验证文件类型
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "只支持 jpg、png、gif 格式的图片")
		return
	}

	// 5. 验证文件大小（限制为5MB）
	if file.Size > 5*1024*1024 {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "图片大小不能超过5MB")
		return
	}

	// 6. 生成文件名（使用时间戳+原文件名）
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	uploadDir := filepath.Join("uploads", "images")
	savePath := filepath.Join(uploadDir, filename)

	// 7. 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		util.ErrorResponse(c, errno.ErrServerError.Code, "创建上传目录失败")
		return
	}

	// 8. 保存文件到本地
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		util.ErrorResponse(c, errno.ErrServerError.Code, "文件保存失败")
		return
	}

	// 9. 生成访问URL（这里假设静态文件服务路径为 /uploads）
	imageURL := "/uploads/images/" + filename

	// 10. 调用服务层保存图片记录
	req := &dto.UploadImageRequest{
		MemoryID: memoryID,
	}
	resp, err := h.imageService.UploadImage(req, imageURL, file.Size)
	if err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 11. 返回成功响应
	util.SuccessResponse(c, resp)
}

// DeleteImage 删除图片
func (h *ImageHandler) DeleteImage(c *gin.Context) {
	// 1. 获取图片ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的图片ID")
		return
	}

	// 2. 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, errno.ErrUnauthorized.Code, errno.ErrUnauthorized.Message)
		return
	}

	// 3. 调用服务层
	if err := h.imageService.DeleteImage(id, userID.(int64)); err != nil {
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 4. 返回成功响应
	util.SuccessResponse(c, nil)
}

// GetImagesByMemoryID 获取记忆的所有图片
func (h *ImageHandler) GetImagesByMemoryID(c *gin.Context) {
	// 1. 获取记忆ID
	idStr := c.Param("memory_id")
	memoryID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, "无效的记忆ID")
		return
	}

	// 2. 调用服务层
	images, err := h.imageService.GetImagesByMemoryID(memoryID)
	if err != nil {
		util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		return
	}

	// 3. 返回成功响应
	util.SuccessResponse(c, images)
}
