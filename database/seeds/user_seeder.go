package seeds

import (
	"news-with-golang/internal/core/domain/model"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	bytes, err := bcrypt.GenerateFromPassword([]byte("hsc999"), 14)
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
