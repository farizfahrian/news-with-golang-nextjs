package handler

import (
	"news-with-golang/internal/adapter/handler/request"
	"news-with-golang/internal/adapter/handler/response"
	"news-with-golang/internal/core/domain/entity"
	"news-with-golang/internal/core/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var err error
var code string
var errorResp response.ErrorResponseDefault
var validate = validator.New()

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	authService service.AuthService
}

// Login implements AuthHandler.
func (a *authHandler) Login(c *fiber.Ctx) error {
	req := request.LoginRequest{}
	resp := response.SuccessAuthResponse{}
	err := c.BodyParser(&req)
	if err != nil {
		code = "[Handler] Login - 1"
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	validateErr := validate.Struct(req)
	if validateErr != nil {
		code = "[Handler] Login - 2"
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: validateErr.Error(),
			},
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqLogin := entity.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := a.authService.GetUserByEmail(c.Context(), reqLogin)
	if err != nil {
		code = "[Handler] Login - 3"
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	resp = response.SuccessAuthResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Success",
		},
		AccessToken: result.AccessToken,
		ExpiresAt:   result.ExpiresAt,
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}
