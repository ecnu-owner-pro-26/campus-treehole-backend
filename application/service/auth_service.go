package service

import (
	"campus-memory/application/dto"
	"campus-memory/infra/model"
	"campus-memory/infra/repo"
	"campus-memory/types/errno"
	"campus-memory/utils"
	"context"
	"errors"

	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	userRepo *repo.UserRepo
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo *repo.UserRepo) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// WechatLogin 微信登录业务逻辑
func (s *AuthService) WechatLogin(req *dto.WechatLoginRequest) (*dto.LoginResponse, error) {
	// 1. 调用微信API获取OpenID
	session, err := utils.GetWechatOpenID(req.Code)
	if err != nil {
		return nil, err
	}

	// 2. 查询用户是否存在
	user, err := s.userRepo.GetUserByOpenID(session.OpenID)

	if err != nil {
		// 用户不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 3. 创建新用户
			user = &model.UserModel{
				OpenID:   session.OpenID,
				UnionID:  session.UnionID,
				Nickname: req.Nickname,
				Avatar:   req.Avatar,
				Status:   1, // 正常状态
				Role:     0, // 普通用户
			}

			// 如果没有提供昵称，使用默认值
			if user.Nickname == "" {
				user.Nickname = "微信用户"
			}

			if err := s.userRepo.CreateUser(user); err != nil {
				return nil, err
			}
		} else {
			// 其他数据库错误
			return nil, err
		}
	} else {
		// 4. 用户已存在，更新信息（如果提供了新信息）
		needUpdate := false

		if req.Nickname != "" && req.Nickname != user.Nickname {
			user.Nickname = req.Nickname
			needUpdate = true
		}

		if req.Avatar != "" && req.Avatar != user.Avatar {
			user.Avatar = req.Avatar
			needUpdate = true
		}

		// 更新 UnionID（如果微信返回了）
		if session.UnionID != "" && session.UnionID != user.UnionID {
			user.UnionID = session.UnionID
			needUpdate = true
		}

		if needUpdate {
			if err := s.userRepo.UpdateUser(user); err != nil {
				return nil, err
			}
		}
	}

	// 5. 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.OpenID)
	if err != nil {
		return nil, err
	}

	// 6. 组装响应
	response := &dto.LoginResponse{
		Token: token,
		User: dto.UserInfoDTO{
			ID:              user.ID,
			Nickname:        user.Nickname,
			Avatar:          user.Avatar,
			DefaultCampusID: user.DefaultCampusID,
		},
	}

	return response, nil
}

// GetUserProfile 获取用户详细信息
func (s *AuthService) GetUserProfile(ctx context.Context, userID int64) (*dto.UserProfileResponse, error) {
	// 查询用户
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}
		return nil, err
	}

	// 组装响应
	response := &dto.UserProfileResponse{
		ID:              user.ID,
		Nickname:        user.Nickname,
		Avatar:          user.Avatar,
		DefaultCampusID: user.DefaultCampusID,
		Status:          user.Status,
		Role:            user.Role,
		CreatedAt:       user.CreatedAt,
	}

	return response, nil
}

// UpdateUserProfile 更新用户信息
func (s *AuthService) UpdateUserProfile(ctx context.Context, userID int64, req *dto.UpdateProfileRequest) error {
	// 查询用户
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ErrUserNotFound
		}
		return err
	}

	// 更新字段（只更新提供的字段）
	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}

	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}

	if req.DefaultCampusID != nil {
		user.DefaultCampusID = req.DefaultCampusID
	}

	// 保存更新
	if err := s.userRepo.UpdateUser(user); err != nil {
		return err
	}

	return nil
}
