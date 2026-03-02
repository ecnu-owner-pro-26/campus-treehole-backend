package service

import (
	"campus-memory/application/dto"
	"campus-memory/infra/model"
	"sort"
	"strings"

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
func (s *QuickNavService) BuildNavTree(campusID int64) ([]dto.CampusNavDTO, error) {
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
	campusNames := map[int64]string{1: "普陀校区", 2: "闵行校区", 3: "滴水湖校区"}

	// 构建返回结果
	var result []dto.CampusNavDTO
	for cid, categories := range campusMap {
		var cats []dto.CategoryDTO
		for cat, locs := range categories {
			cats = append(cats, dto.CategoryDTO{
				Category:  cat,
				Locations: locs,
				Count:     len(locs),
			})
		}
		// 对分类排序
		sort.Slice(cats, func(i, j int) bool {
			return cats[i].Category < cats[j].Category
		})
		result = append(result, dto.CampusNavDTO{
			CampusID:   cid,
			CampusName: campusNames[cid],
			Categories: cats,
		})
	}

	//对校区排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].CampusID < result[j].CampusID
	})
	return result, nil
}

// SearchLocations 搜索地点
func (s *QuickNavService) SearchLocations(keyword string, campusID int64, page, pageSize int) ([]model.LocationModel, int64, error) {
	var locations []model.LocationModel
	var total int64

	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	query := s.db.Model(&model.LocationModel{}).Where("is_active = ?", 1)
	if keyword != "" {
		// 转义关键词中的特殊字符
		escaped := strings.ReplaceAll(keyword, `%`, `\%`)
		escaped = strings.ReplaceAll(escaped, `_`, `\_`)
		// 使用 ESCAPE '\' 来转义
		query = query.Where("name LIKE ? ESCAPE '\\'", "%"+escaped+"%")
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
	if limit <= 0 || limit > 50 {
		limit = 10
	}

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

	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

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
