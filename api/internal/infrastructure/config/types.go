package config

type Config struct {
	AppEnv        string `env:"APP_ENV" envDefault:"development"`
	IsDevelopment bool
	Server        *ServerConfig
	Database      *DatabaseConfig
	Auth          *AuthConfig
}

type ServerConfig struct {
	Port int `env:"PORT" envDefault:"3000"`
}

type DatabaseConfig struct {
	ConnectionString string `env:"DATABASE_URL" envDefault:"postgres://postgres:password@localhost:5432/bouncy?sslmode=disable"`
}

type AuthConfig struct {
	JwtSecret     string `env:"JWT_SECRET" envDefault:"super-secret"`
	RefreshSecret string `env:"JWT_REFRESH_SECRET" envDefault:"refresh-secret"`
	Ttl           int    `env:"JWT_TTL" envDefault:"3600"`
}
