package config

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	SSLMode  string
	Name     string
	User     string
	Password string
}
