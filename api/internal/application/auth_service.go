package application

import (
	"errors"
	"strings"
	"time"

	"github.com/bouncy/bouncy-api/internal/infrastructure/auth"
)

type AuthService struct {
	jwt *JwtService
	ttl time.Duration
}

func NewAuthService(jwt *JwtService, ttl time.Duration) *AuthService {
	return &AuthService{
		jwt: jwt,
		ttl: ttl,
	}
}

func (s AuthService) Login(email, password string) (string, error) {
	// This will be replaced with a DB search for a user
	if !strings.Contains(email, "test") || password != "password" {
		return "", errors.New("invalid credentials")
	}

	if strings.Contains(email, "admin") {
		return s.jwt.GenerateToken(&auth.Claims{
			UserId: "test-admin",
			Email:  email,
			Roles:  []string{"admin"},
		})
	}

	return s.jwt.GenerateToken(&auth.Claims{
		UserId: "test-user",
		Email:  email,
		Roles:  []string{"user"},
	})

}
