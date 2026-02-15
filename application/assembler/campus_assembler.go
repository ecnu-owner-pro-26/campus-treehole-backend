package assembler

// CampusAssembler 校区数据组装器
type CampusAssembler struct{}

// NewCampusAssembler 创建校区组装器
func NewCampusAssembler() *CampusAssembler {
	return &CampusAssembler{}
}

// ToDTO 将校区模型转换为DTO
func (a *CampusAssembler) ToDTO() {
	// TODO: 实现模型到DTO的转换
}
