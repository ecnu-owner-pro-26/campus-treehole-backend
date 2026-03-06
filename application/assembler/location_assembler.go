package assembler

import (
	"campus-memory/application/dto"
	"campus-memory/infra/model"
	"time"
)

// LocationAssembler 地点数据组装器
type LocationAssembler struct{}

// NewLocationAssembler 创建地点组装器
func NewLocationAssembler() *LocationAssembler {
	return &LocationAssembler{}
}

// ToLocationResponse 将 Model 转换为 DTO Response
func (a *LocationAssembler) ToLocationResponse(location *model.LocationModel) *dto.LocationResponse {
	return &dto.LocationResponse{
		ID:          location.ID,
		Name:        location.Name,
		CampusID:    location.CampusID,
		Category:    location.Category,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
		MemoryCount: location.MemoryCount,
		CreatedAt:   location.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToLocationModel 将创建请求转换为 Model
func (a *LocationAssembler) ToLocationModel(req *dto.CreateLocationRequest) *model.LocationModel {
	return &model.LocationModel{
		CampusID:    req.CampusID,
		Name:        req.Name,
		Category:    req.Category,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		IsActive:    true,
		SortOrder:   req.SortOrder,
		MemoryCount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// UpdateLocationModel 更新 Model 字段
func (a *LocationAssembler) UpdateLocationModel(location *model.LocationModel, req *dto.UpdateLocationRequest) {
	if req.Name != nil {
		location.Name = *req.Name
	}
	if req.Category != nil {
		location.Category = *req.Category
	}
	if req.Latitude != nil {
		location.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		location.Longitude = *req.Longitude
	}
	if req.IsActive != nil {
		location.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		location.SortOrder = *req.SortOrder
	}
	location.UpdatedAt = time.Now()
}
