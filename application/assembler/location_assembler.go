package assembler

// LocationAssembler 地点数据组装器
type LocationAssembler struct{}

// NewLocationAssembler 创建地点组装器
func NewLocationAssembler() *LocationAssembler {
	return &LocationAssembler{}
}

// ToDTO 将地点模型转换为DTO
func (a *LocationAssembler) ToDTO() {
	// TODO: 实现模型到DTO的转换
}
