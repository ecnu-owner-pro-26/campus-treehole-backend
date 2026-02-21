package assembler

// CommentAssembler 评论数据组装器
type CommentAssembler struct{}

// NewCommentAssembler 创建评论组装器
func NewCommentAssembler() *CommentAssembler {
	return &CommentAssembler{}
}

// ToDTO 将评论模型转换为DTO
func (a *CommentAssembler) ToDTO() {
	// TODO: 实现模型到DTO的转换
}

// ToDomain 将DTO转换为评论模型
func (a *CommentAssembler) ToDomain() {
	// TODO: 实现DTO到模型的转换
}
