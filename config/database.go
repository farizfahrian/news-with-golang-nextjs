package config

import (
	"fmt"
	"news-with-golang/database/seeds"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (cfg Config) ConnectionPostgres() (*Postgres, error) {
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Psql.User,
		cfg.Psql.Password,
		cfg.Psql.Host,
		cfg.Psql.Port,
		cfg.Psql.Name)

	db, err := gorm.Open(
		postgres.Open(dbConnString),
		&gorm.Config{},
	)

	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres-1] Failed to connect to postgres" + cfg.Psql.Host + ":" + cfg.Psql.Port)
		return nil, err
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres-2] Failed to get database postgres" + cfg.Psql.Host + ":" + cfg.Psql.Port)
		return nil, err
	}

	seeds.SeedRoles(db)

	sqlDB.SetMaxOpenConns(cfg.Psql.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Psql.MaxIdleConns)

	return &Postgres{DB: db}, nil
}
