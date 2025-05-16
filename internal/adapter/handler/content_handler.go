package handler

import (
	"news-with-golang/internal/adapter/handler/response"
	"news-with-golang/internal/core/domain/entity"
	"news-with-golang/internal/core/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ContentHandler interface {
	GetContents(ctx *fiber.Ctx) error
	GetContentById(ctx *fiber.Ctx) error
	CreateContent(ctx *fiber.Ctx) error
	UpdateContent(ctx *fiber.Ctx) error
	DeleteContent(ctx *fiber.Ctx) error
	UploadImageR2(ctx *fiber.Ctx) error
}

type contentHandler struct {
	contentService service.ContentService
}

// CreateContent implements ContentHandler.
func (c *contentHandler) CreateContent(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// DeleteContent implements ContentHandler.
func (c *contentHandler) DeleteContent(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// GetContentById implements ContentHandler.
func (c *contentHandler) GetContentById(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// GetContents implements ContentHandler.
func (c *contentHandler) GetContents(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code := "[Handler] GetContents - 1"
		log.Errorw(code)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized",
			},
		}
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	results, err := c.contentService.GetContents(ctx.Context())
	if err != nil {
		code := "[Handler] GetContents - 2"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	respContents := []response.ContentResponse{}
	for _, result := range results {
		respContents = append(respContents, response.ContentResponse{
			ID:           result.ID,
			Title:        result.Title,
			Excerpt:      result.Excerpt,
			Description:  result.Description,
			Image:        result.Image,
			Tags:         result.Tags,
			Status:       result.Status,
			CategoryID:   result.CategoryID,
			CreatedAt:    result.CreatedAt.Format(time.RFC3339),
			CategoryName: result.Category.Title,
			Author:       result.User.Name,
		})
	}

	defaultSuccessResponse = response.DefaultSuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Contents fetched successfully",
		},
		Data: respContents,
	}

	return ctx.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}

// UpdateContent implements ContentHandler.
func (c *contentHandler) UpdateContent(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// UploadImageR2 implements ContentHandler.
func (c *contentHandler) UploadImageR2(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

func NewContentHandler(contentService service.ContentService) ContentHandler {
	return &contentHandler{
		contentService: contentService,
	}
}
