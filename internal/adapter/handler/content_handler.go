package handler

import (
	"fmt"
	"news-with-golang/internal/adapter/handler/request"
	"news-with-golang/internal/adapter/handler/response"
	"news-with-golang/internal/core/domain/entity"
	"news-with-golang/internal/core/service"
	"news-with-golang/lib/conv"
	validatorLib "news-with-golang/lib/validator"
	"os"
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

	// FE
	GetContentWithQuery(ctx *fiber.Ctx) error
	GetContentDetail(ctx *fiber.Ctx) error
}

type contentHandler struct {
	contentService service.ContentService
}

// GetContentDetail implements ContentHandler.
func (c *contentHandler) GetContentDetail(ctx *fiber.Ctx) error {
	contentIDParam := ctx.Params("contentID")
	contentID, err := conv.StringToInt64(contentIDParam)
	if err != nil {
		code = "[Handler] GetContentDetail - 2"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	result, err := c.contentService.GetContentById(ctx.Context(), contentID)
	if err != nil {
		code = "[Handler] GetContentDetail - 3"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	responseContent := response.SuccessContentResponse{
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
		Data: responseContent,
	}

	return ctx.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}

// GetContentWithQuery implements ContentHandler.
func (c *contentHandler) GetContentWithQuery(ctx *fiber.Ctx) error {
	page := 1
	if ctx.Query("page") != "" {
		page, err = conv.StringToInt(ctx.Query("page"))
		if err != nil {
			code = "[Handler] GetContentWithQuery - 1"
			log.Errorw(code, err)
			errorResp = response.ErrorResponseDefault{
				Meta: response.Meta{
					Status:  false,
					Message: "Invalid page number",
				},
			}
			return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	limit := 6
	if ctx.Query("limit") != "" {
		limit, err = conv.StringToInt(ctx.Query("limit"))
		if err != nil {
			code = "[Handler] GetContentWithQuery - 2"
			log.Errorw(code, err)
			errorResp = response.ErrorResponseDefault{
				Meta: response.Meta{
					Status:  false,
					Message: "Invalid limit number",
				},
			}
			return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	orderBy := "created_at"
	if ctx.Query("orderBy") != "" {
		orderBy = ctx.Query("orderBy")
	}

	orderType := "desc"
	if ctx.Query("orderType") != "" {
		orderType = ctx.Query("orderType")
	}

	search := ""
	if ctx.Query("search") != "" {
		search = ctx.Query("search")
	}

	categoryID := 0
	if ctx.Query("categoryID") != "" {
		categoryID, err = conv.StringToInt(ctx.Query("categoryID"))
		if err != nil {
			code = "[Handler] GetContentWithQuery - 3"
			log.Errorw(code, err)
			errorResp = response.ErrorResponseDefault{
				Meta: response.Meta{
					Status:  false,
					Message: "Invalid category ID",
				},
			}
			return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	query := entity.QueryString{
		Page:       page,
		Limit:      limit,
		Search:     search,
		OrderBy:    orderBy,
		OrderType:  orderType,
		Status:     "PUBLISHED",
		CategoryID: int64(categoryID),
	}

	results, totalData, totalPages, err := c.contentService.GetContents(ctx.Context(), query)
	if err != nil {
		code = "[Handler] GetContentWithQuery - 4"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	responseContent := []response.SuccessContentResponse{}
	for _, content := range results {
		respContent := response.SuccessContentResponse{
			ID:           content.ID,
			Title:        content.Title,
			Excerpt:      content.Excerpt,
			Description:  content.Description,
			Image:        content.Image,
			Tags:         content.Tags,
			Status:       content.Status,
			CategoryID:   content.CategoryID,
			CreatedAt:    content.CreatedAt.Format(time.RFC3339),
			CategoryName: content.Category.Title,
			Author:       content.User.Name,
		}
		responseContent = append(responseContent, respContent)
	}

	defaultSuccessResponse = response.DefaultSuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Contents fetched successfullyz",
		},
		Data: responseContent,
		Pagination: &response.PaginationResponse{
			TotalRecords: int(totalData),
			Page:         page,
			PerPage:      limit,
			TotalPages:   int(totalPages),
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
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
	contentID, err := conv.StringToInt64(contentIDParam)
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
	id, err := conv.StringToInt64(idParam)
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

	page := 1
	if ctx.Query("page") != "" {
		page, err = conv.StringToInt(ctx.Query("page"))
		if err != nil {
			code = "[Handler] GetContents - 2"
			log.Errorw(code, err)
			errorResp = response.ErrorResponseDefault{
				Meta: response.Meta{
					Status:  false,
					Message: "Invalid page number",
				},
			}
			return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	limit := 10
	if ctx.Query("limit") != "" {
		limit, err = conv.StringToInt(ctx.Query("limit"))
		if err != nil {
			code = "[Handler] GetContents - 3"
			log.Errorw(code, err)
			errorResp = response.ErrorResponseDefault{
				Meta: response.Meta{
					Status:  false,
					Message: "Invalid limit number",
				},
			}
			return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	search := ""
	if ctx.Query("search") != "" {
		search = ctx.Query("search")
	}

	orderBy := "created_at"
	if ctx.Query("orderBy") != "" {
		orderBy = ctx.Query("orderBy")
	}

	orderType := "desc"
	if ctx.Query("orderType") != "" {
		orderType = ctx.Query("orderType")
	}

	categoryID := 0
	if ctx.Query("categoryID") != "" {
		categoryID, err = conv.StringToInt(ctx.Query("categoryID"))
		if err != nil {
			code = "[Handler] GetContents - 4"
			log.Errorw(code, err)
			errorResp = response.ErrorResponseDefault{
				Meta: response.Meta{
					Status:  false,
					Message: "Invalid category ID",
				},
			}
			return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	query := entity.QueryString{
		Page:       page,
		Limit:      limit,
		Search:     search,
		OrderBy:    orderBy,
		OrderType:  orderType,
		CategoryID: int64(categoryID),
	}

	results, totalData, totalPages, err := c.contentService.GetContents(ctx.Context(), query)
	if err != nil {
		code = "[Handler] GetContents - 5"
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
		Pagination: &response.PaginationResponse{
			TotalRecords: int(totalData),
			Page:         page,
			PerPage:      limit,
			TotalPages:   int(totalPages),
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}

// UpdateContent implements ContentHandler.
func (c *contentHandler) UpdateContent(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] UpdateContent - 1"
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
		code = "[Handler] UpdateContent - 2"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(&req); err != nil {
		code = "[Handler] UpdateContent - 3"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	idParam := ctx.Params("contentID")
	contentId, err := conv.StringToInt(idParam)
	if err != nil {
		code = "[Handler] UpdateContent - 4"
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
		ID:          int64(contentId),
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: int64(userID),
	}

	err = c.contentService.UpdateContent(ctx.Context(), reqEntity)
	if err != nil {
		code = "[Handler] UpdateContent - 5"
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
			Message: "Content updated successfully",
		},
		Data: nil,
	}

	return ctx.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}

// UploadImageR2 implements ContentHandler.
func (c *contentHandler) UploadImageR2(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] UploadImageR2 - 1"
		log.Errorw(code)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized",
			},
		}
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req request.FileUploadRequest
	file, err := ctx.FormFile("image")
	if err != nil {
		code = "[Handler] UploadImageR2 - 2"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: "Invalid request body",
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	if err = ctx.SaveFile(file, fmt.Sprintf("./temp/content/%s", file.Filename)); err != nil {
		code = "[Handler] UploadImageR2 - 3"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	req.Image = fmt.Sprintf("./temp/content/%s", file.Filename)
	reqEntity := entity.FileUploadEntity{
		Name: fmt.Sprintf("%d-%d", int64(userID), time.Now().UnixNano()),
		Path: req.Image,
	}

	imageUrl, err := c.contentService.UploadImageR2(ctx.Context(), reqEntity)
	if err != nil {
		code = "[Handler] UploadImageR2 - 4"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	if req.Image != "" {
		err = os.Remove(req.Image)
		if err != nil {
			code = "[Handler] UploadImageR2 - 5"
			log.Errorw(code, err)
			errorResp = response.ErrorResponseDefault{
				Meta: response.Meta{
					Status:  false,
					Message: err.Error(),
				},
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
		}
	}

	urlImageResp := map[string]interface{}{
		"url": imageUrl,
	}

	defaultSuccessResponse = response.DefaultSuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Image uploaded successfully",
		},
		Data: urlImageResp,
	}

	return ctx.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func NewContentHandler(contentService service.ContentService) ContentHandler {
	return &contentHandler{
		contentService: contentService,
	}
}
