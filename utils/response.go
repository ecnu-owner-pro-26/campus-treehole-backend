package utils

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(data interface{}) Response {
	return Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(message string, data interface{}) Response {
	return Response{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// Error 错误响应
func Error(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

// BadRequest 400错误
func BadRequest(message string) Response {
	return Error(400, message)
}

// Unauthorized 401错误
func Unauthorized(message string) Response {
	return Error(401, message)
}

// Forbidden 403错误
func Forbidden(message string) Response {
	return Error(403, message)
}

// NotFound 404错误
func NotFound(message string) Response {
	return Error(404, message)
}

// InternalServerError 500错误
func InternalServerError(message string) Response {
	return Error(500, message)
}