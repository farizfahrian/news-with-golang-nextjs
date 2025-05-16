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

type UserHandler interface {
	UpdatePassword(ctx *fiber.Ctx) error
	GetUserByID(ctx *fiber.Ctx) error
}

type userHandler struct {
	userService service.UserService
}

// GetUserByID implements UserHandler.
func (u *userHandler) GetUserByID(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] GetUserByID - 1"
		log.Errorw(code)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized",
			},
		}
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	user, err := u.userService.GetUserByID(ctx.Context(), int64(claims.UserID))
	if err != nil {
		code = "[Handler] GetUserByID - 2"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	userResponse := response.SuccessUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	defaultSuccessResponse := response.DefaultSuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "User fetched successfully",
		},
		Data: userResponse,
	}

	return ctx.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}

// UpdatePassword implements UserHandler.
func (u *userHandler) UpdatePassword(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] UpdatePassword - 1"
		log.Errorw(code)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: "Unauthorized",
			},
		}
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req request.UpdatePasswordRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		code = "[Handler] UpdatePassword - 2"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: "Bad Request",
			},
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(&req); err != nil {
		code = "[Handler] UpdatePassword - 3"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = u.userService.UpdatePassword(ctx.Context(), req.NewPassword, int64(claims.UserID))
	if err != nil {
		code = "[Handler] UpdatePassword - 4"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse := response.DefaultSuccessResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Password updated successfully",
		},
		Data: nil,
	}

	return ctx.Status(fiber.StatusOK).JSON(defaultSuccessResponse)
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService: userService}
}
