package service

import (
	"campus-memory/application/assembler"
	"campus-memory/application/dto"
	"campus-memory/infra/model"
	"campus-memory/infra/repo"
	"campus-memory/types/errno"
	"context"
	"gorm.io/gorm"
)

// MemoryService 记忆服务
type MemoryService struct {
	memoryRepo   *repo.MemoryRepo
	userRepo     *repo.UserRepo
	locationRepo *repo.LocationRepo
	likeRepo     *repo.LikeRepo
	imageRepo    *repo.ImageRepo
	assembler    *assembler.MemoryAssembler
}

// NewMemoryService 创建记忆服务实例
func NewMemoryService(
	memoryRepo *repo.MemoryRepo,
	userRepo *repo.UserRepo,
	locationRepo *repo.LocationRepo,
	likeRepo *repo.LikeRepo,
	imageRepo *repo.ImageRepo,
) *MemoryService {
	return &MemoryService{
		memoryRepo:   memoryRepo,
		userRepo:     userRepo,
		locationRepo: locationRepo,
		likeRepo:     likeRepo,
		imageRepo:    imageRepo,
		assembler:    assembler.NewMemoryAssembler(),
	}
}

// CreateMemory 创建记忆
func (s *MemoryService) CreateMemory(req *dto.CreateMemoryRequest, creatorID int64) (*dto.MemoryResponse, error) {
	// 1. 验证地点是否存在
	location, err := s.locationRepo.GetByID(*req.LocationID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrLocationNotFound
		}
		return nil, err
	}

	// 2. 转换为 Model
	memory := s.assembler.ToMemoryModel(req, creatorID, location.Name)

	// 3. 保存到数据库
	if err := s.memoryRepo.Create(memory); err != nil {
		return nil, errno.ErrMemoryCreateFail
	}

	// 4. 保存图片关联(如果有)
	if len(req.ImageURLs) > 0 {
		for i, url := range req.ImageURLs {
			image := &model.ImageModel{
				MemoryID:  memory.ID,
				URL:       url,
				SortOrder: i,
			}
			// 忽略图片保存错误,不影响主流程
			_ = s.imageRepo.Create(image)
		}
	}

	// 5. 获取创建者信息
	creator, _ := s.userRepo.GetUserByID(context.Background(), creatorID)

	// 6. 获取图片列表
	images, _ := s.imageRepo.GetByMemoryID(memory.ID)

	// 7. 组装响应
	return s.assembler.ToMemoryResponse(memory, creator, images, false), nil
}

// GetMemory 获取记忆详情
func (s *MemoryService) GetMemory(id int64, currentUserID *int64) (*dto.MemoryResponse, error) {
	// 1. 获取记忆
	memory, err := s.memoryRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrMemoryNotFound
		}
		return nil, err
	}

	// 2. 增加浏览次数
	_ = s.memoryRepo.IncrementViewCount(id)

	// 3. 获取创建者信息
	creator, err := s.userRepo.GetUserByID(context.Background(), memory.CreatorID)
	if err != nil {
		return nil, err
	}

	// 4. 获取图片列表
	images, _ := s.imageRepo.GetByMemoryID(id)

	// 5. 检查当前用户是否已点赞
	isLiked := false
	if currentUserID != nil {
		isLiked, _ = s.likeRepo.CheckLiked(*currentUserID, id, dto.LikeTargetTypeMemory)
	}

	// 6. 组装响应
	return s.assembler.ToMemoryResponse(memory, creator, images, isLiked), nil
}

// ListMemories 获取记忆列表
func (s *MemoryService) ListMemories(req *dto.MemoryListRequest, currentUserID *int64) (*dto.MemoryListResponse, error) {
	// 1. 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.SortBy == "" {
		req.SortBy = "latest"
	}

	// 2. 查询记忆列表
	memories, total, err := s.memoryRepo.List(req.LocationID, req.Page, req.PageSize, req.SortBy)
	if err != nil {
		return nil, err
	}

	// 3. 组装响应列表
	memoryResponses := make([]dto.MemoryResponse, 0, len(memories))
	for _, memory := range memories {
		// 获取创建者
		creator, _ := s.userRepo.GetUserByID(context.Background(), memory.CreatorID)
		if creator == nil {
			continue
		}

		// 获取图片
		images, _ := s.imageRepo.GetByMemoryID(memory.ID)

		// 检查是否已点赞
		isLiked := false
		if currentUserID != nil {
			isLiked, _ = s.likeRepo.CheckLiked(*currentUserID, memory.ID, dto.LikeTargetTypeMemory)
		}

		memoryResponses = append(memoryResponses, *s.assembler.ToMemoryResponse(memory, creator, images, isLiked))
	}

	return &dto.MemoryListResponse{
		Memories: memoryResponses,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// UpdateMemory 更新记忆
func (s *MemoryService) UpdateMemory(id int64, req *dto.UpdateMemoryRequest, userID int64) error {
	// 1. 获取记忆
	memory, err := s.memoryRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errno.ErrMemoryNotFound
		}
		return err
	}

	// 2. 权限检查
	if memory.CreatorID != userID {
		return errno.ErrForbidden
	}

	// 3. 如果更新了 LocationID，需要同步更新 LocationName
	if req.LocationID != nil && (memory.LocationID == nil || *req.LocationID != *memory.LocationID) {
		location, err := s.locationRepo.GetByID(*req.LocationID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errno.ErrLocationNotFound
			}
			return err
		}
		memory.LocationName = location.Name
	}

	// 4. 更新其他字段
	s.assembler.UpdateMemoryModel(memory, req)

	// 5. 保存到数据库
	if err := s.memoryRepo.Update(memory); err != nil {
		return errno.ErrMemoryUpdateFail
	}

	// 6. 更新图片关联(如果有)
	if req.ImageURLs != nil {
		// 删除旧图片
		_ = s.imageRepo.DeleteByMemoryID(id)
		// 添加新图片
		for i, url := range req.ImageURLs {
			image := &model.ImageModel{
				MemoryID:  id,
				URL:       url,
				SortOrder: i,
			}
			_ = s.imageRepo.Create(image)
		}
	}

	return nil
}

// DeleteMemory 删除记忆
func (s *MemoryService) DeleteMemory(id int64, userID int64) error {
	// 1. 获取记忆
	memory, err := s.memoryRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errno.ErrMemoryNotFound
		}
		return err
	}

	// 2. 权限检查
	if memory.CreatorID != userID {
		return errno.ErrForbidden
	}

	// 3. 软删除
	if err := s.memoryRepo.Delete(id); err != nil {
		return errno.ErrMemoryDeleteFail
	}

	return nil
}
