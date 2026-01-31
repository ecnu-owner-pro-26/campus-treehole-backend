package errno

// Error 自定义错误类型
type Error struct {
	Code    int
	Message string
}

// Error 实现error接口
func (e *Error) Error() string {
	return e.Message
}

// New 创建新错误
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// 预定义错误
var (
	ErrMemoryNotFound    = New(1001, "记忆不存在")
	ErrMemoryCreateFail  = New(1002, "创建记忆失败")
	ErrCommentNotFound   = New(2001, "留言不存在")
	ErrCommentCreateFail = New(2002, "创建留言失败")
	ErrFileUploadFail    = New(3001, "文件上传失败")
)
