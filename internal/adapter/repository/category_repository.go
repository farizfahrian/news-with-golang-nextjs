package repository

import (
	"context"
	"errors"
	"fmt"
	"news-with-golang/internal/core/domain/entity"
	"news-with-golang/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryById(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	UpdateCategory(ctx context.Context, id int64, req entity.CategoryEntity) error
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryRepository struct {
	db *gorm.DB
}

// CreateCategory implements CategoryRepository.
func (c *categoryRepository) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	var countSlug int64
	err = c.db.Table("categories").Where("slug = ?", req.Slug).Count(&countSlug).Error
	if err != nil {
		code := "[Repository] CreateCategory -1"
		log.Errorw(code, err)
		return err
	}

	countSlug = countSlug + 1
	slug := fmt.Sprintf("%s-%d", req.Slug, countSlug)
	modelCategory := model.Category{
		Title:       req.Title,
		Slug:        slug,
		CreatedByID: req.User.ID,
	}

	err = c.db.Create(&modelCategory).Error
	if err != nil {
		code := "[Repository] CreateCategory -2"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// DeleteCategory implements CategoryRepository.
func (c *categoryRepository) DeleteCategory(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// GetCategories implements CategoryRepository.
func (c *categoryRepository) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	var modelCagories []model.Category

	err = c.db.Order("created_at DESC").Preload("User").Find(&modelCagories).Error
	if err != nil {
		code := "[Repository] GetCategories -1"
		log.Errorw(code, err)
		return nil, err
	}

	if len(modelCagories) == 0 {
		code := "[Repository] GetCategories -2"
		err = errors.New("category not found")
		log.Errorw(code, err)
		return nil, err
	}

	var categories []entity.CategoryEntity
	for _, modelCategory := range modelCagories {
		categories = append(categories, entity.CategoryEntity{
			ID:    modelCategory.ID,
			Title: modelCategory.Title,
			Slug:  modelCategory.Slug,
			User: entity.UserEntity{
				ID:       modelCategory.User.ID,
				Name:     modelCategory.User.Name,
				Email:    modelCategory.User.Email,
				Password: modelCategory.User.Password,
			},
		})
	}
	return categories, nil
}

// GetCategoryById implements CategoryRepository.
func (c *categoryRepository) GetCategoryById(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	panic("unimplemented")
}

// UpdateCategory implements CategoryRepository.
func (c *categoryRepository) UpdateCategory(ctx context.Context, id int64, req entity.CategoryEntity) error {
	panic("unimplemented")
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}
