package handler

// ImageHandler 图片上传处理器
type ImageHandler struct {
	// TODO: 注入依赖
}

// NewImageHandler 创建图片处理器实例
func NewImageHandler() *ImageHandler {
	return &ImageHandler{}
}

// UploadImage 上传图片
func (h *ImageHandler) UploadImage() {
	// TODO: 实现图片上传逻辑
}

// DeleteImage 删除图片
func (h *ImageHandler) DeleteImage() {
	// TODO: 实现图片删除逻辑
}
