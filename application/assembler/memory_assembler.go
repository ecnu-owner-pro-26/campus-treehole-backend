package assembler

import (
	"campus-memory/application/dto"
	"campus-memory/infra/model"
	"encoding/json"
	"time"
)

// MemoryAssembler 记忆数据组装器
type MemoryAssembler struct{}

// NewMemoryAssembler 创建记忆组装器
func NewMemoryAssembler() *MemoryAssembler {
	return &MemoryAssembler{}
}

// ToMemoryResponse 将 Model 转换为 DTO Response
func (a *MemoryAssembler) ToMemoryResponse(
	memory *model.MemoryModel,
	creator *model.UserModel,
	images []*model.ImageModel,
	isLiked bool,
) *dto.MemoryResponse {
	// 解析标签
	var tags []string
	if memory.Tags != "" {
		json.Unmarshal([]byte(memory.Tags), &tags)
	}

	// 转换图片列表
	imageInfos := make([]dto.ImageInfo, 0, len(images))
	for _, img := range images {
		imageInfos = append(imageInfos, dto.ImageInfo{
			ID:  img.ID,
			URL: img.URL,
		})
	}

	// 组装响应
	return &dto.MemoryResponse{
		ID:           memory.ID,
		Title:        memory.Title,
		Content:      memory.Content,
		LocationName: memory.LocationName,
		LocationID:   memory.LocationID,
		Latitude:     memory.Latitude,
		Longitude:    memory.Longitude,
		Creator: dto.UserSimpleInfo{
			ID:       creator.ID,
			Nickname: creator.Nickname,
			Avatar:   creator.Avatar,
		},
		LikeCount:    memory.LikeCount,
		CommentCount: memory.CommentCount,
		ViewCount:    memory.ViewCount,
		IsLiked:      isLiked,
		Tags:         tags,
		Images:       imageInfos,
		CreatedAt:    memory.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToMemoryModel 将创建请求转换为 Model
func (a *MemoryAssembler) ToMemoryModel(req *dto.CreateMemoryRequest, creatorID int64, locationName string) *model.MemoryModel {
	// 序列化标签
	tagsJSON, _ := json.Marshal(req.Tags)

	return &model.MemoryModel{
		Title:        req.Title,
		Content:      req.Content,
		LocationID:   req.LocationID,
		LocationName: locationName,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		IsPublic:     req.IsPublic,
		Tags:         string(tagsJSON),
		Status:       1, // 已发布
		CreatorID:    creatorID,
		ViewCount:    0,
		LikeCount:    0,
		CommentCount: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// UpdateMemoryModel 更新 Model 字段
func (a *MemoryAssembler) UpdateMemoryModel(memory *model.MemoryModel, req *dto.UpdateMemoryRequest) {
	if req.Title != nil {
		memory.Title = *req.Title
	}
	if req.Content != nil {
		memory.Content = *req.Content
	}
	if req.LocationID != nil {
		memory.LocationID = req.LocationID
	}
	if req.Latitude != nil {
		memory.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		memory.Longitude = *req.Longitude
	}
	if req.IsPublic != nil {
		memory.IsPublic = *req.IsPublic
	}
	if req.Tags != nil {
		tagsJSON, _ := json.Marshal(req.Tags)
		memory.Tags = string(tagsJSON)
	}
	memory.UpdatedAt = time.Now()
}
