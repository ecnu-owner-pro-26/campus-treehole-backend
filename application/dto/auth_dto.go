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
	// 基本信息
	ID              int64     `json:"id"`                        // 用户ID
	Nickname        string    `json:"nickname"`                  // 用户昵称
	Avatar          string    `json:"avatar"`                    // 用户头像URL
	BackgroundImage string    `json:"backgroundImage"`           // 背景图片URL
	Bio             string    `json:"bio"`                       // 个人简介
	DefaultCampusID *int64    `json:"defaultCampusId,omitempty"` // 默认校区ID（可为空）
	Status          int8      `json:"status"`                    // 状态（0-禁用 1-正常）
	Role            int8      `json:"role"`                      // 角色（0-普通用户 1-管理员）
	CreatedAt       time.Time `json:"createdAt"`                 // 注册时间
	// 统计数据
	MemoryCount  int64 `json:"memoryCount"`  // 发布的记忆数
	TotalLikes   int64 `json:"totalLikes"`   // 获得的点赞总数
	CommentCount int64 `json:"commentCount"` // 发布的评论数
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Nickname        *string `json:"nickname" binding:"omitempty,min=1,max=50"` // 用户昵称（1-50字符）
	Avatar          *string `json:"avatar" binding:"omitempty,url"`            // 用户头像URL（必须是有效URL）
	DefaultCampusID *int64  `json:"defaultCampusId" binding:"omitempty,min=1"` // 默认校区ID（必须大于0）
	BackgroundImage *string `json:"backgroundImage" binding:"omitempty,url"` // 背景图片
	Bio             *string `json:"bio" binding:"omitempty,max=200"`         // 个人简介（最多200字）
}

// UserSimpleDTO 用户简化信息（用于记忆、评论的创建者信息）
type UserSimpleDTO struct {
	ID       int64  `json:"id"`       // 用户ID
	Nickname string `json:"nickname"` // 用户昵称
	Avatar   string `json:"avatar"`   // 用户头像URL
}
// PublicUserProfileResponse 他人主页响应（只返回公开信息）
type PublicUserProfileResponse struct {
	ID              int64     `json:"id"`
	Nickname        string    `json:"nickname"`
	Avatar          string    `json:"avatar"`
	BackgroundImage string    `json:"backgroundImage"`
	Bio             string    `json:"bio"`
	DefaultCampusID *int64    `json:"defaultCampusId,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	MemoryCount     int64     `json:"memoryCount"`
	TotalLikes      int64     `json:"totalLikes"`
	CommentCount    int64     `json:"commentCount"`
}

