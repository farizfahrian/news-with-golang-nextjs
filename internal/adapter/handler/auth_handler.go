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

var err error
var code string
var errorResp response.ErrorResponseDefault

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	authService service.AuthService
}

// Login implements AuthHandler.
func (a *authHandler) Login(c *fiber.Ctx) error {
	req := request.LoginRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		code = "[Handler] Login - 1"
		log.Errorw(code, err)
		errorResp = response.ErrorResponseDefault{
			Meta: response.Meta{
				Status:  false,
				Message: err.Error(),
			},
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	validateErr := validatorLib.ValidateStruct(req)
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
		if err.Error() == "invalid password" {
			return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	resp := response.SuccessAuthResponse{
		Meta: response.Meta{
			Status:  true,
			Message: "Login success",
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
