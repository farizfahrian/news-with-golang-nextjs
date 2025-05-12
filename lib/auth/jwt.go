package auth

import (
	"fmt"
	"news-with-golang/config"
	"news-with-golang/internal/core/domain/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt interface {
	GenerateToken(data entity.JwtData) (string, int64, error)
	VerifyAccessToken(token string) (*entity.JwtData, error)
}

type Options struct {
	signingKey string
	issuer     string
}

// GenerateToken implements Jwt.
func (o *Options) GenerateToken(data entity.JwtData) (string, int64, error) {
	now := time.Now().Local()
	expiresAt := now.Add(time.Hour * 24)
	data.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	data.RegisteredClaims.IssuedAt = jwt.NewNumericDate(now)
	data.RegisteredClaims.NotBefore = jwt.NewNumericDate(now)
	data.RegisteredClaims.Issuer = o.issuer

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	accessToken, err := token.SignedString([]byte(o.signingKey))
	if err != nil {
		return "", 0, err
	}

	return accessToken, expiresAt.Unix(), nil
}

// VerifyAccessToken implements Jwt.
func (o *Options) VerifyAccessToken(token string) (*entity.JwtData, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(o.signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if parsedToken.Valid {
		claim, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok || !parsedToken.Valid {
			return nil, err
		}

		return &entity.JwtData{
			UserId: claim["user_id"].(float64),
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func NewJwt(cfg *config.Config) Jwt {
	opt := new(Options)
	opt.signingKey = cfg.App.JwtSecretKey
	opt.issuer = cfg.App.JwtIssuer

	return opt
}
