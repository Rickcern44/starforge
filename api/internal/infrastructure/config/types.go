package config

type Config struct {
	Server   *ServerConfig
	Database *DatabaseConfig
	Auth     *AuthConfig
}

type ServerConfig struct {
	Port int
}
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type AuthConfig struct {
	JwtSecret     string
	RefreshSecret string
	Ttl           int `koanf:"jwt_ttl"`
}
