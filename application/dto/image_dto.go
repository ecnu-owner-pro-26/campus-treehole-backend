package dto

// ImageInfo 图片信息
type ImageInfo struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

// UploadImageRequest 上传图片请求
type UploadImageRequest struct {
	MemoryID int64 `json:"memory_id" binding:"required"`
}

// UploadImageResponse 上传图片响应
type UploadImageResponse struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}
