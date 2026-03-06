package assembler

import (
	"campus-memory/application/dto"
	"campus-memory/infra/model"
)

// CampusAssembler 校区数据组装器
type CampusAssembler struct{}

// NewCampusAssembler 创建校区组装器
func NewCampusAssembler() *CampusAssembler {
	return &CampusAssembler{}
}

// ToCampusResponse 将 Model 转换为 DTO Response
func (a *CampusAssembler) ToCampusResponse(campus *model.CampusModel) *dto.CampusResponse {
	return &dto.CampusResponse{
		ID:        campus.ID,
		Name:      campus.Name,
		CreatedAt: campus.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToCampusModel 将创建请求转换为 Model
func (a *CampusAssembler) ToCampusModel(name string, sortOrder int) *model.CampusModel {
	return &model.CampusModel{
		Name:      name,
		IsActive:  true,
		SortOrder: sortOrder,
	}
}

// UpdateCampusModel 更新 Model 字段
func (a *CampusAssembler) UpdateCampusModel(campus *model.CampusModel, name *string, isActive *bool, sortOrder *int) {
	if name != nil {
		campus.Name = *name
	}
	if isActive != nil {
		campus.IsActive = *isActive
	}
	if sortOrder != nil {
		campus.SortOrder = *sortOrder
	}
}
