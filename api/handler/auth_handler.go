package handler

import (
	"campus-memory/application/dto"
	"campus-memory/infra/util"
	"context"

	"github.com/gin-gonic/gin"
)

// AuthServiceInterface 定义了认证服务需要实现的方法，用于依赖注入和测试
type AuthServiceInterface interface {
	WechatLogin(ctx context.Context, req *dto.WechatLoginRequest) (*dto.LoginResponse, error)
	GetUserProfile(ctx context.Context, userID int64) (*dto.UserProfileResponse, error)
	UpdateUserProfile(ctx context.Context, userID int64, req *dto.UpdateProfileRequest) error
}

// AuthHandler 认证处理器
type AuthHandler struct {
	authService AuthServiceInterface
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService AuthServiceInterface) *AuthHandler {
	return &AuthHandler{authService: authService}
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
	resp, err := h.authService.WechatLogin(c.Request.Context(), &req)
	if err != nil {
		util.ErrorResponse(c, 500, "登录失败: "+err.Error())
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
		util.ErrorResponse(c, 401, "未登录")
		return
	}

	// 调用服务层
	resp, err := h.authService.GetUserProfile(c.Request.Context(), userID.(int64))
	if err != nil {
		util.ErrorResponse(c, 500, "获取用户信息失败: "+err.Error())
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
		util.ErrorResponse(c, 401, "未登录")
		return
	}

	// 绑定请求参数
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 调用服务层
	err := h.authService.UpdateUserProfile(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		util.ErrorResponse(c, 500, "更新用户信息失败: "+err.Error())
		return
	}

	// 返回成功响应
	util.SuccessResponse(c, gin.H{"message": "更新成功"})
}
