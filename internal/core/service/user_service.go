package service

import (
	"context"

	"news-with-golang/internal/adapter/repository"
	"news-with-golang/internal/core/domain/entity"
	"news-with-golang/lib/conv"
)

type UserService interface {
	UpdatePassword(ctx context.Context, newPassword string, id int64) error
	GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error)
}

type userService struct {
	userRepository repository.UserRepository
}

// GetUserByID implements UserService.
func (u *userService) GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error) {
	result, err := u.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdatePassword implements UserService.
func (u *userService) UpdatePassword(ctx context.Context, newPassword string, id int64) error {
	password, err := conv.HashPassword(newPassword)
	if err != nil {
		return err
	}

	err = u.userRepository.UpdatePassword(ctx, password, id)
	if err != nil {
		return err
	}

	return nil
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
