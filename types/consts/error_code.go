package consts

// 错误码定义
const (
	// 通用错误码
	ErrCodeSuccess      = 200
	ErrCodeBadRequest   = 400
	ErrCodeUnauthorized = 401
	ErrCodeForbidden    = 403
	ErrCodeNotFound     = 404
	ErrCodeServerError  = 500

	// 记忆相关 (11xxx)
	ErrCodeMemoryNotFound   = 11001
	ErrCodeMemoryCreateFail = 11002

	// 评论相关 (12xxx)
	ErrCodeCommentNotFound   = 12001
	ErrCodeCommentCreateFail = 12002

	// 文件相关 (13xxx)
	ErrCodeFileUploadFail = 13001
)
