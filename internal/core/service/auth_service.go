package service

import (
	"context"
	"errors"
	"news-with-golang/config"
	"news-with-golang/internal/adapter/repository"
	"news-with-golang/internal/core/domain/entity"
	"news-with-golang/lib/auth"
	"news-with-golang/lib/conv"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

var code string

type AuthService interface {
	GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error)
}

type authService struct {
	authRepository repository.AuthRepository
	cfg            *config.Config
	jwtToken       auth.Jwt
}

// GetUserByEmail implements AuthService.
func (a *authService) GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error) {
	result, err := a.authRepository.GetUserByEmail(ctx, req)
	if err != nil {
		code = "[Service] GetUserByEmail - 1"
		log.Errorw(code, err)
		return nil, err
	}

	if !conv.CheckPasswordHas(req.Password, result.Password) {
		code = "[Service] GetUserByEmail - 2"
		err = errors.New("invalid password")
		log.Errorw(code, err)
		return nil, err
	}

	jwtData := entity.JwtData{
		UserID: float64(result.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			ID:        strconv.FormatInt(result.ID, 10),
		},
	}

	accessToken, expiresAt, err := a.jwtToken.GenerateToken(&jwtData)
	if err != nil {
		code = "[Service] GetUserByEmail - 3"
		log.Errorw(code, err)
		return nil, err
	}

	return &entity.AccessToken{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}

func NewAuthService(authRepository repository.AuthRepository, cfg *config.Config, jwtToken auth.Jwt) AuthService {
	return &authService{
		authRepository: authRepository,
		cfg:            cfg,
		jwtToken:       jwtToken,
	}
}
