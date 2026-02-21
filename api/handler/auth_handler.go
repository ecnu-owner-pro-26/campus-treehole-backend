package handler

import (
	"campus-memory/application/dto"
	"campus-memory/application/service"
	"campus-memory/infra/util"

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
		util.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 2. 调用服务层
	resp, err := h.authService.WechatLogin(&req)
	if err != nil {
		util.ErrorResponse(c, 500, "登录失败: "+err.Error())
		return
	}

	// 3. 返回成功响应
	util.SuccessResponse(c, resp)
}

// GetProfile 获取个人信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// 1. 从context获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, 401, "未登录")
		return
	}

	// 2. 调用服务层
	resp, err := h.authService.GetUserProfile(userID.(int64))
	if err != nil {
		util.ErrorResponse(c, 500, "获取用户信息失败: "+err.Error())
		return
	}

	// 3. 返回成功响应
	util.SuccessResponse(c, resp)
}

// UpdateProfile 更新个人信息
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// 1. 从context获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		util.ErrorResponse(c, 401, "未登录")
		return
	}

	// 2. 绑定请求参数
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 3. 调用服务层
	err := h.authService.UpdateUserProfile(userID.(int64), &req)
	if err != nil {
		util.ErrorResponse(c, 500, "更新用户信息失败: "+err.Error())
		return
	}

	// 4. 返回成功响应
	util.SuccessResponse(c, gin.H{"message": "更新成功"})
}
