package interfaces

import "github.com/bouncy/bouncy-api/internal/infrastructure/auth"

type JwtService interface {
	Validate(token string) (*auth.Claims, error)
	GenerateToken(claims *auth.Claims) (string, error)
}

type AuthService interface {
	Login(email, password string) (string, error)
}
