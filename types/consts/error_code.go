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

	// 业务错误码
	ErrCodeMemoryNotFound    = 1001
	ErrCodeMemoryCreateFail  = 1002
	ErrCodeCommentNotFound   = 2001
	ErrCodeCommentCreateFail = 2002
	ErrCodeFileUploadFail    = 3001
)
