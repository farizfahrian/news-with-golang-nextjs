package repository

import (
	"context"
	"news-with-golang/internal/core/domain/entity"
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]entity.CategoryEntity, error)
}
