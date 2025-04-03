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
	Name     string `json:"name"`

	MaxOpenConns    int `json:"max_open_conns"`
	MaxIdleConns    int `json:"max_idle_conns"`
	ConnMaxLifetime int `json:"conn_max_lifetime"`
	ConnMaxIdleTime int `json:"conn_max_idle_time"`
}

type Config struct {
	App  App
	Psql PsqlDB
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
			Host:     viper.GetString("psql.host"),
			Port:     viper.GetString("psql.port"),
			User:     viper.GetString("psql.user"),
			Password: viper.GetString("psql.password"),
			Name:     viper.GetString("psql.name"),

			MaxOpenConns:    viper.GetInt("psql.max_open_conns"),
			MaxIdleConns:    viper.GetInt("psql.max_idle_conns"),
			ConnMaxLifetime: viper.GetInt("psql.conn_max_lifetime"),
			ConnMaxIdleTime: viper.GetInt("psql.conn_max_idle_time"),
		},
	}
}
