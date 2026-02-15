package handler

import (
	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	// TODO: 注入依赖
}

// WechatLogin 微信登录
func (h *AuthHandler) WechatLogin(c *gin.Context) {
	// TODO: 实现微信登录逻辑
}

// GetProfile 获取个人信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// TODO: 实现获取个人信息逻辑
}

// UpdateProfile 更新个人信息
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// TODO: 实现更新个人信息逻辑
}
