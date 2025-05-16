package handler

import (
	"news-with-golang/internal/adapter/handler/request"
	"news-with-golang/internal/adapter/handler/response"
	"news-with-golang/internal/core/domain/entity"
	"news-with-golang/internal/core/service"
	"news-with-golang/lib/conv"
	validatorLib "news-with-golang/lib/validator"
	"strings"
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
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] CreateContent - 1"
		log.Errorw(code)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized",
			},
		}
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req request.ContentRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		code = "[Handler] CreateContent - 2"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[Handler] CreateContent - 3"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: int64(userID),
	}

	err = c.contentService.CreateContent(ctx.Context(), reqEntity)
	if err != nil {
		code = "[Handler] CreateContent - 4"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse = response.DefaultSuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Content created successfully",
		},
		Data: nil,
	}

	return ctx.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

// DeleteContent implements ContentHandler.
func (c *contentHandler) DeleteContent(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] DeleteContent - 1"
		log.Errorw(code)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized",
			},
		}
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	contentIDParam := ctx.Params("contentID")
	contentID, err := conv.StringToInt(contentIDParam)
	if err != nil {
		code = "[Handler] DeleteContent - 2"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	err = c.contentService.DeleteContent(ctx.Context(), contentID)
	if err != nil {
		code = "[Handler] DeleteContent - 3"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse = response.DefaultSuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Content deleted successfully",
		},
		Data: nil,
	}

	return ctx.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}

// GetContentById implements ContentHandler.
func (c *contentHandler) GetContentById(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] GetContentById - 1"
		log.Errorw(code)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized",
			},
		}
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := ctx.Params("contentID")
	id, err := conv.StringToInt(idParam)
	if err != nil {
		code = "[Handler] GetContentById - 2"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := c.contentService.GetContentById(ctx.Context(), id)
	if err != nil {
		code = "[Handler] GetContentById - 3"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	contentResponse := response.SuccessContentResponse{
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
	}

	defaultSuccessResponse = response.DefaultSuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Content fetched successfully",
		},
		Data:       contentResponse,
		Pagination: nil,
	}

	return ctx.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}

// GetContents implements ContentHandler.
func (c *contentHandler) GetContents(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] GetContents - 1"
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
		code = "[Handler] GetContents - 2"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	respContents := []response.SuccessContentResponse{}
	for _, result := range results {
		respContents = append(respContents, response.SuccessContentResponse{
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
