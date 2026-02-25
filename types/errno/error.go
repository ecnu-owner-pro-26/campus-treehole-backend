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
	// 通用错误 (2xx, 4xx, 5xx)
	ErrSuccess      = New(200, "成功")
	ErrBadRequest   = New(400, "请求参数错误")
	ErrUnauthorized = New(401, "未登录或token无效")
	ErrForbidden    = New(403, "无权限访问")
	ErrNotFound     = New(404, "资源不存在")
	ErrServerError  = New(500, "服务器内部错误")

	// 记忆相关 (11xxx)
	ErrMemoryNotFound   = New(11001, "记忆不存在")
	ErrMemoryCreateFail = New(11002, "创建记忆失败")
	ErrMemoryUpdateFail = New(11003, "更新记忆失败")
	ErrMemoryDeleteFail = New(11004, "删除记忆失败")

	// 评论相关 (12xxx)
	ErrCommentNotFound   = New(12001, "评论不存在")
	ErrCommentCreateFail = New(12002, "创建评论失败")
	ErrCommentDeleteFail = New(12003, "删除评论失败")

	// 文件相关 (13xxx)
	ErrFileUploadFail = New(13001, "文件上传失败")

	// 点赞相关 (14xxx)
	ErrAlreadyLiked   = New(14001, "已经点赞过了")
	ErrNotLiked       = New(14002, "还没有点赞")
	ErrLikeCreateFail = New(14003, "点赞失败")
	ErrLikeDeleteFail = New(14004, "取消点赞失败")

	// 用户相关 (15xxx)
	ErrUserNotFound    = New(15001, "用户不存在")
	ErrWechatLoginFail = New(15002, "微信登录失败")

	// 校区地点相关 (16xxx)
	ErrCampusNotFound   = New(16001, "校区不存在")
	ErrLocationNotFound = New(16002, "地点不存在")
)
