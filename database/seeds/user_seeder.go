package seeds

import (
	"news-with-golang/internal/core/domain/model"
	"news-with-golang/lib/conv"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	bytes, err := conv.HashPassword("hsc999")
	if err != nil {
		log.Fatal().Err(err).Msg("[SeedRoles] Failed to generate password")
	}

	user := model.User{
		Name:     "Admin",
		Email:    "admin@gmail.com",
		Password: string(bytes),
	}

	err = db.FirstOrCreate(&user, model.User{Email: "admin@gmail.com"}).Error
	if err != nil {
		log.Fatal().Err(err).Msg("[SeedRoles] Failed to create user")
	} else {
		log.Info().Msg("[SeedRoles] User created")
	}
}
