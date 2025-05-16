package handler

import (
	"news-with-golang/internal/adapter/handler/request"
	"news-with-golang/internal/adapter/handler/response"
	"news-with-golang/internal/core/domain/entity"
	"news-with-golang/internal/core/service"
	"news-with-golang/lib/conv"
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
	EditCategoryById(ctx *fiber.Ctx) error
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

// EditCategoryById implements CategoryHandler.
func (c *categoryHandler) EditCategoryById(ctx *fiber.Ctx) error {
	var req request.CategoryRequest
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code := "[Handler] EditCategoryById - 1"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized",
			},
		}
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = ctx.BodyParser(&req); err != nil {
		code := "[Handler] EditCategoryById - 2"
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
		code := "[Handler] EditCategoryById - 3"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	idParam := ctx.Params("categoryID")
	id, err := conv.StringToInt(idParam)
	if err != nil {
		code := "[Handler] EditCategoryById - 4"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.CategoryEntity{
		ID:    id,
		Title: req.Title,
		User: entity.UserEntity{
			ID: int64(userID),
		},
	}

	err = c.categoryService.EditCategoryById(ctx.Context(), reqEntity)
	if err != nil {
		code := "[Handler] EditCategoryById - 5"
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
			Message: "Category updated successfully",
		},
		Data:       nil,
		Pagination: nil,
	}

	return ctx.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
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
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code := "[Handler] GetCategoryById - 1"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized",
			},
		}
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := ctx.Params("categoryID")
	id, err := conv.StringToInt(idParam)
	if err != nil {
		code := "[Handler] GetCategoryById - 2"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	category, err := c.categoryService.GetCategoryById(ctx.Context(), id)
	if err != nil {
		code := "[Handler] GetCategoryById - 3"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	categoryResponse := response.SuccessCategoryResponse{
		ID:            category.ID,
		Title:         category.Title,
		Slug:          category.Slug,
		CreatedByName: category.User.Name,
	}

	defaultSuccessResponse = response.DefaultSuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Category fetched successfully",
		},
		Data:       categoryResponse,
		Pagination: nil,
	}

	return ctx.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &categoryHandler{categoryService: categoryService}
}
