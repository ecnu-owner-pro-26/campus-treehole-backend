package service

import (
	"campus-memory/application/dto"
	"campus-memory/infra/model"
	"campus-memory/infra/repo"
	"campus-memory/types/errno"

	"gorm.io/gorm"
)

// ImageService 图片服务
type ImageService struct {
	imageRepo  *repo.ImageRepo
	memoryRepo *repo.MemoryRepo
}

// NewImageService 创建图片服务实例
func NewImageService(
	imageRepo *repo.ImageRepo,
	memoryRepo *repo.MemoryRepo,
) *ImageService {
	return &ImageService{
		imageRepo:  imageRepo,
		memoryRepo: memoryRepo,
	}
}

// UploadImage 上传图片（创建图片记录）
func (s *ImageService) UploadImage(req *dto.UploadImageRequest, url string, size int64) (*dto.UploadImageResponse, error) {
	// 1. 验证记忆是否存在
	memory, err := s.memoryRepo.GetByID(req.MemoryID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrMemoryNotFound
		}
		return nil, err
	}

	// 2. 获取当前记忆的图片数量，用于设置排序
	images, _ := s.imageRepo.GetByMemoryID(memory.ID)
	sortOrder := len(images)

	// 3. 创建图片记录
	image := &model.ImageModel{
		MemoryID:  req.MemoryID,
		URL:       url,
		Size:      size,
		SortOrder: sortOrder,
	}

	if err := s.imageRepo.Create(image); err != nil {
		return nil, errno.ErrImageUploadFail
	}

	// 4. 组装响应
	return &dto.UploadImageResponse{
		ID:  image.ID,
		URL: image.URL,
	}, nil
}

// DeleteImage 删除图片
func (s *ImageService) DeleteImage(id int64, userID int64) error {
	// 1. 获取图片信息
	image, err := s.imageRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errno.ErrImageNotFound
		}
		return err
	}

	// 2. 获取关联的记忆，验证权限
	memory, err := s.memoryRepo.GetByID(image.MemoryID)
	if err != nil {
		return err
	}

	// 3. 权限检查（只能删除自己记忆的图片）
	if memory.CreatorID != userID {
		return errno.ErrForbidden
	}

	// 4. 删除图片记录
	if err := s.imageRepo.Delete(id); err != nil {
		return errno.ErrImageDeleteFail
	}

	return nil
}

// GetImagesByMemoryID 获取记忆的所有图片
func (s *ImageService) GetImagesByMemoryID(memoryID int64) ([]dto.ImageInfo, error) {
	images, err := s.imageRepo.GetByMemoryID(memoryID)
	if err != nil {
		return nil, err
	}

	imageInfos := make([]dto.ImageInfo, 0, len(images))
	for _, img := range images {
		imageInfos = append(imageInfos, dto.ImageInfo{
			ID:  img.ID,
			URL: img.URL,
		})
	}

	return imageInfos, nil
}
