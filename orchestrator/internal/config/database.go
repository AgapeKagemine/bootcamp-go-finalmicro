package config

type DatabaseConfig struct {
	Driver   string
	Username string
	Password uint
	Host     string
	Port     uint
	Database string
}
