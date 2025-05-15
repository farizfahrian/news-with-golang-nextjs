package handler

import (
	"news-with-golang/internal/adapter/handler/request"
	"news-with-golang/internal/adapter/handler/response"
	"news-with-golang/internal/core/domain/entity"
	"news-with-golang/internal/core/service"
	validatorLib "news-with-golang/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var defaultSuccessResponse response.DefaultSuccessResponse

type CategoryHandler interface {
	CreateCategory(ctx *fiber.Ctx) error
	DeleteCategory(ctx *fiber.Ctx) error
	GetCategories(ctx *fiber.Ctx) error
	GetCategoryById(ctx *fiber.Ctx) error
	UpdateCategory(ctx *fiber.Ctx) error
}

type categoryHandler struct {
	categoryService service.CategoryService
}

// CreateCategory implements CategoryHandler.
func (c *categoryHandler) CreateCategory(ctx *fiber.Ctx) error {
	var req request.CategoryRequest
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code := "[Handler] CreateCategory - 1"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized",
			},
		}
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	err := ctx.BodyParser(&req)
	if err != nil {
		code := "[Handler] CreateCategory - 2"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: "Bad Request",
			},
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		code := "[Handler] CreateCategory - 3"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	categoryEntity := entity.CategoryEntity{
		Title: req.Title,
		User: entity.UserEntity{
			ID: int64(userID),
		},
	}
	err = c.categoryService.CreateCategory(ctx.Context(), categoryEntity)
	if err != nil {
		code := "[Handler] CreateCategory - 4"
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
			Message: "Category created successfully",
		},
		Data:       nil,
		Pagination: nil,
	}

	return ctx.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}

// DeleteCategory implements CategoryHandler.
func (c *categoryHandler) DeleteCategory(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// GetCategories implements CategoryHandler.
func (c *categoryHandler) GetCategories(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code := "[Handler] GetCategories - 1"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized",
			},
		}
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	results, err := c.categoryService.GetCategories(ctx.Context())
	if err != nil {
		code := "[Handler] GetCategories - 2"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	categoryResponses := []response.SuccessCategoryResponse{}
	for _, result := range results {
		categoryResponses = append(categoryResponses, response.SuccessCategoryResponse{
			ID:            result.ID,
			Title:         result.Title,
			Slug:          result.Slug,
			CreatedByName: result.User.Name,
		})
	}
	defaultSuccessResponse = response.DefaultSuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Categories fetched successfully",
		},
		Data:       categoryResponses,
		Pagination: nil,
	}

	return ctx.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}

// GetCategoryById implements CategoryHandler.
func (c *categoryHandler) GetCategoryById(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// UpdateCategory implements CategoryHandler.
func (c *categoryHandler) UpdateCategory(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &categoryHandler{categoryService: categoryService}
}
