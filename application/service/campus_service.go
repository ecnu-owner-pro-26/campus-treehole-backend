package service

import (
	"campus-memory/application/assembler"
	"campus-memory/application/dto"
	"campus-memory/infra/repo"
	"campus-memory/types/errno"

	"gorm.io/gorm"
)

// CampusService 校区服务
type CampusService struct {
	campusRepo   *repo.CampusRepo
	locationRepo *repo.LocationRepo
	assembler    *assembler.CampusAssembler
}

// NewCampusService 创建校区服务实例
func NewCampusService(
	campusRepo *repo.CampusRepo,
	locationRepo *repo.LocationRepo,
) *CampusService {
	return &CampusService{
		campusRepo:   campusRepo,
		locationRepo: locationRepo,
		assembler:    assembler.NewCampusAssembler(),
	}
}

// ListCampuses 获取所有校区列表
func (s *CampusService) ListCampuses() (*dto.CampusListResponse, error) {
	// 1. 查询所有校区
	campuses, err := s.campusRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// 2. 组装响应列表
	campusResponses := make([]dto.CampusResponse, 0, len(campuses))
	for _, campus := range campuses {
		campusResponses = append(campusResponses, *s.assembler.ToCampusResponse(&campus))
	}

	return &dto.CampusListResponse{
		Campuses: campusResponses,
		Total:    int64(len(campuses)),
	}, nil
}

// GetCampus 获取校区详情
func (s *CampusService) GetCampus(id int64) (*dto.CampusResponse, error) {
	// 1. 获取校区
	campus, err := s.campusRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrCampusNotFound
		}
		return nil, err
	}

	// 2. 组装响应
	return s.assembler.ToCampusResponse(campus), nil
}

// GetCampusWithLocations 获取校区及其地点列表
func (s *CampusService) GetCampusWithLocations(campusID int64) (*dto.CampusLocationsResponse, error) {
	// 1. 获取校区
	campus, err := s.campusRepo.GetByID(campusID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.ErrCampusNotFound
		}
		return nil, err
	}

	// 2. 获取该校区的所有地点
	locations, err := s.locationRepo.GetByCampusID(campusID)
	if err != nil {
		return nil, err
	}

	// 3. 组装响应
	campusResponse := s.assembler.ToCampusResponse(campus)
	locationResponses := make([]dto.LocationResponse, 0, len(locations))
	locationAssembler := assembler.NewLocationAssembler()
	for _, location := range locations {
		locationResponses = append(locationResponses, *locationAssembler.ToLocationResponse(&location))
	}

	return &dto.CampusLocationsResponse{
		Campus:    *campusResponse,
		Locations: locationResponses,
	}, nil
}
