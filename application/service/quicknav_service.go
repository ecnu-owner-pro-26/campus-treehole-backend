package service

import (
	"campus-memory/infra/model"

	"gorm.io/gorm"
)

// QuickNavService 快速导航服务
type QuickNavService struct {
	db *gorm.DB
}

// NewQuickNavService 创建服务实例
func NewQuickNavService(db *gorm.DB) *QuickNavService {
	return &QuickNavService{db: db}
}

// BuildNavTree 构建导航树
func (s *QuickNavService) BuildNavTree(campusID int64) ([]map[string]interface{}, error) {
	var locations []model.LocationModel
	query := s.db.Where("is_active = ?", 1).Order("campus_id, category, sort_order")

	if campusID > 0 {
		query = query.Where("campus_id = ?", campusID)
	}

	if err := query.Find(&locations).Error; err != nil {
		return nil, err
	}

	// 构建树形结构
	campusMap := make(map[int64]map[string][]model.LocationModel)
	for _, loc := range locations {
		if _, ok := campusMap[loc.CampusID]; !ok {
			campusMap[loc.CampusID] = make(map[string][]model.LocationModel)
		}
		campusMap[loc.CampusID][loc.Category] = append(campusMap[loc.CampusID][loc.Category], loc)
	}

	// 校区名称映射
	campusNames := map[int64]string{1: "主校区", 2: "东校区", 3: "西校区"}

	// 构建返回结果
	var result []map[string]interface{}
	for campusID, categories := range campusMap {
		var cats []map[string]interface{}
		for category, locs := range categories {
			cats = append(cats, map[string]interface{}{
				"category":  category,
				"locations": locs,
				"count":     len(locs),
			})
		}
		result = append(result, map[string]interface{}{
			"campus_id":   campusID,
			"campus_name": campusNames[campusID],
			"categories":  cats,
		})
	}
	return result, nil
}

// SearchLocations 搜索地点
func (s *QuickNavService) SearchLocations(keyword string, campusID int64, page, pageSize int) ([]model.LocationModel, int64, error) {
	var locations []model.LocationModel
	var total int64

	query := s.db.Model(&model.LocationModel{}).Where("is_active = ?", 1)
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if campusID > 0 {
		query = query.Where("campus_id = ?", campusID)
	}

	query.Count(&total)
	err := query.Order("memory_count desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&locations).Error
	return locations, total, err
}

// GetPopularLocations 获取热门地点
func (s *QuickNavService) GetPopularLocations(campusID int64, limit int) ([]model.LocationModel, error) {
	var locations []model.LocationModel
	query := s.db.Where("is_active = ?", 1).Order("memory_count desc")
	if campusID > 0 {
		query = query.Where("campus_id = ?", campusID)
	}
	err := query.Limit(limit).Find(&locations).Error
	return locations, err
}

// GetLocationsByCategory 根据类别获取地点
func (s *QuickNavService) GetLocationsByCategory(campusID int64, category string, page, pageSize int) ([]model.LocationModel, int64, error) {
	var locations []model.LocationModel
	var total int64

	// 查询总数
	s.db.Model(&model.LocationModel{}).
		Where("campus_id = ? AND category = ? AND is_active = ?", campusID, category, 1).
		Count(&total)

	// 查询数据
	err := s.db.Where("campus_id = ? AND category = ? AND is_active = ?", campusID, category, 1).
		Order("sort_order asc, memory_count desc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&locations).Error

	return locations, total, err
}
