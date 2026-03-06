package dto

import "time"

// WechatLoginRequest 微信登录请求
type WechatLoginRequest struct {
	Code     string `json:"code" binding:"required"`      // 微信wx.login()返回的code（必需）
	Nickname string `json:"nickname" binding:"omitempty"` // 用户昵称
	Avatar   string `json:"avatar" binding:"omitempty"`   // 用户头像URL
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"` // JWT token
	User  UserInfoDTO `json:"user"`  // 用户基本信息
}

// UserInfoDTO 用户基本信息（用于登录响应、评论作者等场景）
type UserInfoDTO struct {
	ID              int64  `json:"id"`                        // 用户ID
	Nickname        string `json:"nickname"`                  // 用户昵称
	Avatar          string `json:"avatar"`                    // 用户头像URL
	DefaultCampusID *int64 `json:"defaultCampusId,omitempty"` // 默认校区ID（可为空）
}

// UserProfileResponse 用户详细信息响应（用于获取个人资料）
type UserProfileResponse struct {
	ID              int64     `json:"id"`                        // 用户ID
	Nickname        string    `json:"nickname"`                  // 用户昵称
	Avatar          string    `json:"avatar"`                    // 用户头像URL
	DefaultCampusID *int64    `json:"defaultCampusId,omitempty"` // 默认校区ID（可为空）
	Status          int8      `json:"status"`                    // 状态（0-禁用 1-正常）
	Role            int8      `json:"role"`                      // 角色（0-普通用户 1-管理员）
	CreatedAt       time.Time `json:"createdAt"`                 // 注册时间
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Nickname        *string `json:"nickname" binding:"omitempty,min=1,max=50"` // 用户昵称（1-50字符）
	Avatar          *string `json:"avatar" binding:"omitempty,url"`            // 用户头像URL（必须是有效URL）
	DefaultCampusID *int64  `json:"defaultCampusId" binding:"omitempty,min=1"` // 默认校区ID（必须大于0）
}

// UserSimpleDTO 用户简化信息（用于记忆、评论的创建者信息）
type UserSimpleDTO struct {
	ID       int64  `json:"id"`       // 用户ID
	Nickname string `json:"nickname"` // 用户昵称
	Avatar   string `json:"avatar"`   // 用户头像URL
}
