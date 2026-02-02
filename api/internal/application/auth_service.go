package application

import (
	"errors"
	"log/slog"
	"net/mail"
	"time"

	"github.com/bouncy/bouncy-api/internal/application/interfaces"
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user with that email already exists")
)

type AuthService struct {
	jwt      *JwtService
	userRepo interfaces.UserRepository
}

func NewAuthService(jwt *JwtService, userRepo interfaces.UserRepository) *AuthService {
	return &AuthService{
		jwt:      jwt,
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		// Consider logging the internal error, but return a generic one
		return "", ErrUserNotFound
	}

	if !checkPasswordHash(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	return s.jwt.GenerateToken(&auth.Claims{
		UserId: user.ID,
		Email:  user.Email,
		Roles:  user.Roles,
	})
}

func (s *AuthService) Register(name, email, password string) error {
	if err := s.ValidateLoginRequirements(email, password); err != nil {
		return err
	}

	// Check if user already exists
	_, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		slog.Info("user with that email already exists")
		return ErrUserAlreadyExists
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		// Log internal error
		return errors.New("internal server error")
	}

	newUser := &models.User{
		ID:           uuid.NewString(),
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		Roles:        []string{"user"}, // Assign []string directly
		CreatedAt:    time.Now(),
	}

	return s.userRepo.CreateUser(newUser)
}

func (s *AuthService) ValidateLoginRequirements(email, password string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("invalid email address")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
