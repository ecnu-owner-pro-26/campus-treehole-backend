package service

import (
	"campus-memory/application/assembler"
	"campus-memory/application/dto"
	"campus-memory/infra/repo"
	"campus-memory/types/errno"

	"gorm.io/gorm"
)

// LocationService 地点服务
type LocationService struct {
	locationRepo *repo.LocationRepo
	campusRepo   *repo.CampusRepo
	assembler    *assembler.LocationAssembler
}

// NewLocationService 创建地点服务实例
func NewLocationService(
	locationRepo *repo.LocationRepo,
	campusRepo *repo.CampusRepo,
) *LocationService {
	return &LocationService{
		locationRepo: locationRepo,
		campusRepo:   campusRepo,
		assembler:    assembler.NewLocationAssembler(),
	}
}

// GetLocation 获取地点详情
func (s *LocationService) GetLocation(id int64) (*dto.LocationResponse, error) {
	// 1. 获取地点
	location, err := s.locationRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrLocationNotFound
		}
		return nil, err
	}

	// 2. 组装响应
	return s.assembler.ToLocationResponse(location), nil
}

// ListLocations 获取地点列表
func (s *LocationService) ListLocations(req *dto.LocationListRequest) (*dto.LocationListResponse, error) {
	// 1. 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 2. 查询地点列表
	locations, total, err := s.locationRepo.List(req.CampusID, req.Category, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// 3. 组装响应列表
	locationResponses := make([]dto.LocationResponse, 0, len(locations))
	for _, location := range locations {
		locationResponses = append(locationResponses, *s.assembler.ToLocationResponse(&location))
	}

	return &dto.LocationListResponse{
		Locations: locationResponses,
		Total:     total,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}, nil
}

// CreateLocation 创建地点
func (s *LocationService) CreateLocation(req *dto.CreateLocationRequest) (*dto.LocationResponse, error) {
	// 1. 验证校区是否存在
	_, err := s.campusRepo.GetByID(req.CampusID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrCampusNotFound
		}
		return nil, err
	}

	// 2. 转换为 Model
	location := s.assembler.ToLocationModel(req)

	// 3. 保存到数据库
	if err := s.locationRepo.Create(location); err != nil {
		return nil, errno.ErrLocationCreateFail
	}

	// 4. 组装响应
	return s.assembler.ToLocationResponse(location), nil
}

// UpdateLocation 更新地点
func (s *LocationService) UpdateLocation(id int64, req *dto.UpdateLocationRequest) error {
	// 1. 获取地点
	location, err := s.locationRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errno.ErrLocationNotFound
		}
		return err
	}

	// 2. 更新字段
	s.assembler.UpdateLocationModel(location, req)

	// 3. 保存到数据库
	if err := s.locationRepo.Update(location); err != nil {
		return errno.ErrLocationUpdateFail
	}

	return nil
}

// DeleteLocation 删除地点
func (s *LocationService) DeleteLocation(id int64) error {
	// 1. 获取地点
	_, err := s.locationRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errno.ErrLocationNotFound
		}
		return err
	}

	// 2. 软删除
	if err := s.locationRepo.Delete(id); err != nil {
		return errno.ErrLocationDeleteFail
	}

	return nil
}

// SearchLocations 搜索地点
func (s *LocationService) SearchLocations(keyword string, campusID *int64) ([]dto.LocationResponse, error) {
	// 1. 搜索地点
	locations, err := s.locationRepo.Search(keyword, campusID)
	if err != nil {
		return nil, err
	}

	// 2. 组装响应列表
	locationResponses := make([]dto.LocationResponse, 0, len(locations))
	for _, location := range locations {
		locationResponses = append(locationResponses, *s.assembler.ToLocationResponse(&location))
	}

	return locationResponses, nil
}
