package dto

// CampusResponse 校区响应
type CampusResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// CampusListResponse 校区列表响应
type CampusListResponse struct {
	Campuses []CampusResponse `json:"campuses"`
	Total    int64            `json:"total"`
}

// CampusLocationsResponse 校区地点列表响应
type CampusLocationsResponse struct {
	Campus    CampusResponse     `json:"campus"`
	Locations []LocationResponse `json:"locations"`
}

// CampusWithLocations 带地点的校区
type CampusWithLocations struct {
	Campus    CampusResponse     `json:"campus"`
	Locations []LocationResponse `json:"locations"`
}
