package middleware

import (
	"news-with-golang/config"
	"news-with-golang/internal/adapter/handler/response"
	"news-with-golang/lib/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Middleware interface {
	CheckToken() fiber.Handler
}

type Options struct {
	authJwt auth.Jwt
}

// CheckToken implements Middleware.
func (o *Options) CheckToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var errorResponse response.ErrorResponseDefault
		authHandler := c.Get("Authorization")
		if authHandler == "" {
			errorResponse.Meta = response.Meta{
				Status:  false,
				Message: "Unauthorized",
			}
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		tokenString := strings.Split(authHandler, "Bearer ")[1]
		claims, err := o.authJwt.VerifyAccessToken(tokenString)
		if err != nil {
			errorResponse.Meta = response.Meta{
				Status:  false,
				Message: "Unauthorized",
			}
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		c.Locals("user_id", claims)

		return c.Next()
	}
}

func NewMiddleware(cfg *config.Config) Middleware {
	opt := new(Options)
	opt.authJwt = auth.NewJwt(cfg)

	return opt
}
