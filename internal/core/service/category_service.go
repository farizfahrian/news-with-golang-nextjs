package service

import (
	"context"
	"news-with-golang/internal/adapter/repository"
	"news-with-golang/internal/core/domain/entity"
	"news-with-golang/lib/conv"

	"github.com/gofiber/fiber/v2/log"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategory(ctx context.Context, id int64) error
	GetCategories(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryById(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	UpdateCategory(ctx context.Context, id int64, req entity.CategoryEntity) error
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

// CreateCategory implements CategoryService.
func (c *categoryService) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	slug := conv.GenerateSlug(req.Title)
	req.Slug = slug

	err := c.categoryRepository.CreateCategory(ctx, req)
	if err != nil {
		code := "[Service] CreateCategory - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// DeleteCategory implements CategoryService.
func (c *categoryService) DeleteCategory(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// GetCategories implements CategoryService.
func (c *categoryService) GetCategories(ctx context.Context) ([]entity.CategoryEntity, error) {
	results, err := c.categoryRepository.GetCategories(ctx)
	if err != nil {
		code := "[Service] GetCategories - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return results, nil
}

// GetCategoryById implements CategoryService.
func (c *categoryService) GetCategoryById(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	results, err := c.categoryRepository.GetCategoryById(ctx, id)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// UpdateCategory implements CategoryService.
func (c *categoryService) UpdateCategory(ctx context.Context, id int64, req entity.CategoryEntity) error {
	err := c.categoryRepository.UpdateCategory(ctx, id, req)
	if err != nil {
		return err
	}

	return nil
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository: categoryRepository}
}
