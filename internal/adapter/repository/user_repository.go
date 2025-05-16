package repository

import (
	"context"

	"news-with-golang/internal/core/domain/entity"
	"news-with-golang/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type UserRepository interface {
	UpdatePassword(ctx context.Context, newPassword string, id int64) error
	GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error)
}

type userRepository struct {
	db *gorm.DB
}

// GetUserByID implements UserRepository.
func (u *userRepository) GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error) {
	var modelUser model.User
	err = u.db.Where("id = ?", id).First(&modelUser).Error
	if err != nil {
		code := "[Repository] GetUserByID -1"
		log.Errorw(code, err)
		return nil, err
	}

	return &entity.UserEntity{
		ID:    modelUser.ID,
		Name:  modelUser.Name,
		Email: modelUser.Email,
	}, nil
}

// UpdatePassword implements UserRepository.
func (u *userRepository) UpdatePassword(ctx context.Context, newPassword string, id int64) error {
	err = u.db.Model(&model.User{}).Where("id = ?", id).Update("password", newPassword).Error
	if err != nil {
		code := "[Repository] UpdatePassword -1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
