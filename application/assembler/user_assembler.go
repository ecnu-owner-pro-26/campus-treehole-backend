package assembler

// UserAssembler 用户数据组装器
type UserAssembler struct{}

// NewUserAssembler 创建用户组装器
func NewUserAssembler() *UserAssembler {
	return &UserAssembler{}
}

// ToDTO 将用户模型转换为DTO
func (a *UserAssembler) ToDTO() {
	// TODO: 实现模型到DTO的转换
}

// ToSimpleDTO 将用户模型转换为简单DTO
func (a *UserAssembler) ToSimpleDTO() {
	// TODO: 实现模型到简单DTO的转换
}
