package handler

import (
	"campus-memory/application/dto"
	"campus-memory/application/service"
	"campus-memory/infra/util"
	"campus-memory/types/errno"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// WechatLogin 微信登录
func (h *AuthHandler) WechatLogin(c *gin.Context) {
	// 1. 绑定请求参数
	var req dto.WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, errno.ErrBadRequest.Message)
		return
	}

	// 2. 调用服务层
	resp, err := h.authService.WechatLogin(&req)
	if err != nil {
		// 判断是否是自定义错误
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrWechatLoginFail.Code, errno.ErrWechatLoginFail.Message)
		}
		return
	}

	// 3. 返回成功响应
	util.SuccessResponse(c, resp)
}

// GetProfile godoc
// @Summary 获取个人信息
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserProfileResponse
// @Router /api/auth/profile [get]
// GetProfile 获取个人信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// 从context获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, errno.ErrUnauthorized.Code, errno.ErrUnauthorized.Message)
		return
	}

	// 调用服务层
	resp, err := h.authService.GetUserProfile(c.Request.Context(), userID.(int64))
	if err != nil {
		// 判断是否是自定义错误
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 返回成功响应
	util.SuccessResponse(c, resp)
}

// UpdateProfile godoc
// @Summary 更新个人信息
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.UpdateProfileRequest true "更新信息"
// @Success 200 {object} util.Response
// @Router /api/auth/profile [put]
// UpdateProfile 更新个人信息
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// 1. 从context获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, errno.ErrUnauthorized.Code, errno.ErrUnauthorized.Message)
		return
	}

	// 绑定请求参数
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, errno.ErrBadRequest.Code, errno.ErrBadRequest.Message)
		return
	}

	// 调用服务层
	err := h.authService.UpdateUserProfile(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		// 判断是否是自定义错误
		if e, ok := err.(*errno.Error); ok {
			util.ErrorResponse(c, e.Code, e.Message)
		} else {
			util.ErrorResponse(c, errno.ErrServerError.Code, errno.ErrServerError.Message)
		}
		return
	}

	// 返回成功响应
	util.SuccessResponse(c, gin.H{"message": "更新成功"})
}
