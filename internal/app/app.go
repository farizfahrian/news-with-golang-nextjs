package app

import (
	"news-with-golang/config"
	"news-with-golang/lib/auth"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

func RunServer() {
	cfg := config.NewConfig()
	_, err := cfg.ConnectionPostgres()

	if err != nil {
		log.Fatal().Msgf("Error connecting to database %v", err)
		return
	}

	cfdR2 := cfg.LoadAwsConfig()
	_ = s3.NewFromConfig(cfdR2)

	_ = auth.NewJwt(cfg)
}
