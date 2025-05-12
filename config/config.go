package config

import "github.com/spf13/viper"

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`

	JwtSecretKey string `json:"jwt_secret_key"`
	JwtIssuer    string `json:"jwt_issuer"`
	JwtExpire    int64  `json:"jwt_expire"`
}

type PsqlDB struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"db_name"`

	MaxOpenConns    int `json:"max_open_conns"`
	MaxIdleConns    int `json:"max_idle_conns"`
	ConnMaxLifetime int `json:"conn_max_lifetime"`
	ConnMaxIdleTime int `json:"conn_max_idle_time"`
}

type CloudflareR2 struct {
	Name      string `json:"name"`
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	Token     string `json:"token"`
	AccountId string `json:"account_id"`
	PublicUrl string `json:"public_url"`
}

type Config struct {
	App  App
	Psql PsqlDB
	R2   CloudflareR2
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort: viper.GetString("app_port"),
			AppEnv:  viper.GetString("app_env"),

			JwtSecretKey: viper.GetString("jwt_secret_key"),
			JwtIssuer:    viper.GetString("jwt_issuer"),
			JwtExpire:    viper.GetInt64("jwt_expire"),
		},
		Psql: PsqlDB{
			Host:     viper.GetString("DATABASE_HOST"),
			Port:     viper.GetString("DATABASE_PORT"),
			User:     viper.GetString("DATABASE_USER"),
			Password: viper.GetString("DATABASE_PASSWORD"),
			Name:     viper.GetString("DATABASE_NAME"),

			MaxOpenConns:    viper.GetInt("DATABASE_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("DATABASE_MAX_IDLE_CONNS"),
			ConnMaxLifetime: viper.GetInt("DATABASE_CONN_MAX_LIFETIME"),
			ConnMaxIdleTime: viper.GetInt("DATABASE_CONN_MAX_IDLE_TIME"),
		},
		R2: CloudflareR2{
			Name:      viper.GetString("CLOUDFLARE_R2_BUCKET_NAME"),
			ApiKey:    viper.GetString("CLOUDFLARE_R2_API_KEY"),
			ApiSecret: viper.GetString("CLOUDFLARE_R2_API_SECRET"),
			Token:     viper.GetString("CLOUDFLARE_R2_TOKEN"),
			AccountId: viper.GetString("CLOUDFLARE_R2_ACCOUNT_ID"),
			PublicUrl: viper.GetString("CLOUDFLARE_R2_PUBLIC_URL"),
		},
	}
}
