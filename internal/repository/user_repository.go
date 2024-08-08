package repository

import (
	"context"

	"github.com/temuka-api-service/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetFollowers(ctx context.Context, userId int) ([]model.UserFollow, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, userId int, user *model.User) error
	DeleteUser(ctx context.Context, id int) error
	CreateUserFollow(ctx context.Context, user_follow *model.UserFollow) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepositoryImpl) CreateUserFollow(ctx context.Context, user_follow *model.UserFollow) error {
	return r.db.WithContext(ctx).Create(user_follow).Error
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetAllUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, userId int, user *model.User) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userId).Updates(user).Error
}

func (r *UserRepositoryImpl) GetFollowers(ctx context.Context, userID int) ([]model.UserFollow, error) {
	var followers []model.UserFollow
	if err := r.db.WithContext(ctx).Where("following_id = ?", userID).Find(&followers).Error; err != nil {
		return nil, err
	}
	return followers, nil
}
