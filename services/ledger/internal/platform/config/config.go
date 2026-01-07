package config

type DatabaseConfig struct {
	Port     int64  `validate:"required"`
	Host     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
}
