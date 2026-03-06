package repo

import (
	"campus-memory/infra/model"
	"context"
	"time"

	"gorm.io/gorm"
)

// UserRepo 用户数据访问层
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo 创建用户数据访问层
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetUserByOpenID 根据微信OpenID获取用户
func (r *UserRepo) GetUserByOpenID(openid string) (*model.UserModel, error) {
	var user model.UserModel
	err := r.db.Where("openid = ?", openid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUnionID 根据微信UnionID获取用户（用于跨应用识别）
func (r *UserRepo) GetUserByUnionID(unionid string) (*model.UserModel, error) {
	var user model.UserModel
	err := r.db.Where("unionid = ? AND unionid != ''", unionid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建新用户（首次微信登录时）
func (r *UserRepo) CreateUser(user *model.UserModel) error {
	return r.db.Create(user).Error
}

// GetUserByID 根据用户ID获取用户
func (r *UserRepo) GetUserByID(ctx context.Context, id int64) (*model.UserModel, error) {
	var user model.UserModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息（昵称、头像、默认校区等）
func (r *UserRepo) UpdateUser(user *model.UserModel) error {
	user.UpdatedAt = time.Now()
	return r.db.Save(user).Error
}
