package service

import (
	"context"
	"news-with-golang/config"
	"news-with-golang/internal/adapter/cloudflare"
	"news-with-golang/internal/adapter/repository"
	"news-with-golang/internal/core/domain/entity"

	"github.com/gofiber/fiber/v2/log"
)

type ContentService interface {
	GetContents(ctx context.Context) ([]entity.ContentEntity, error)
	GetContentById(ctx context.Context, id int64) (entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	UpdateContent(ctx context.Context, req entity.ContentEntity) error
	DeleteContent(ctx context.Context, id int64) error
	UploadImageR2(ctx context.Context, req entity.FileUploadEntity) (string, error)
}

type contentService struct {
	contentRepo repository.ContentRepository
	cfg         *config.Config
	r2          cloudflare.CloudflareR2Adapter
}

// CreateContent implements ContentService.
func (c *contentService) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	panic("unimplemented")
}

// DeleteContent implements ContentService.
func (c *contentService) DeleteContent(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// GetContentById implements ContentService.
func (c *contentService) GetContentById(ctx context.Context, id int64) (entity.ContentEntity, error) {
	panic("unimplemented")
}

// GetContents implements ContentService.
func (c *contentService) GetContents(ctx context.Context) ([]entity.ContentEntity, error) {
	result, err := c.contentRepo.GetContents(ctx)
	if err != nil {
		code := "[Service] GetContents - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

// UpdateContent implements ContentService.
func (c *contentService) UpdateContent(ctx context.Context, req entity.ContentEntity) error {
	panic("unimplemented")
}

// UploadImageR2 implements ContentService.
func (c *contentService) UploadImageR2(ctx context.Context, req entity.FileUploadEntity) (string, error) {
	panic("unimplemented")
}

func NewContentService(contentRepository repository.ContentRepository, cfg *config.Config, r2 cloudflare.CloudflareR2Adapter) ContentService {
	return &contentService{
		contentRepo: contentRepository,
		cfg:         cfg,
		r2:          r2,
	}
}
