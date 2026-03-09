package application

import (
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"time"

	"github.com/bouncy/bouncy-api/internal/application/payments"
	"github.com/bouncy/bouncy-api/internal/domain/interfaces"
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user with that email already exists")
	ErrInvalidInvitation  = errors.New("invalid or expired invitation")
)

type AuthService struct {
	jwt            *JwtService
	userRepo       interfaces.UserRepository
	emailService   interfaces.EmailService
	paymentService *payments.Service
	leagueRepo     interfaces.LeagueMemberRepository
}

func NewAuthService(jwt *JwtService, userRepo interfaces.UserRepository, email interfaces.EmailService, payments *payments.Service, leagueRepo interfaces.LeagueMemberRepository) *AuthService {
	return &AuthService{
		jwt:            jwt,
		userRepo:       userRepo,
		emailService:   email,
		paymentService: payments,
		leagueRepo:     leagueRepo,
	}
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
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

func (s *AuthService) InviteLeagueCreator(email, invitedBy string) error {
	token := uuid.NewString()
	inv := &models.Invitation{
		Token:       token,
		Email:       email,
		LeagueID:    "platform", // Global invitation
		InvitedBy:   invitedBy,
		SpecialRole: "league_creator",
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.userRepo.CreateInvitation(inv); err != nil {
		return err
	}

	return s.emailService.SendInvitation(email, token, "Bouncy Platform")
}

func (s *AuthService) InviteUser(email, leagueID, invitedBy string) error {
	isAdmin, err := s.leagueRepo.IsAdmin(leagueID, invitedBy)
	if err != nil {
		return err
	}
	if !isAdmin {
		return fmt.Errorf("insufficient permissions to invite to this league")
	}

	token := uuid.NewString()
	inv := &models.Invitation{
		Token:     token,
		Email:     email,
		LeagueID:  leagueID,
		InvitedBy: invitedBy,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := s.userRepo.CreateInvitation(inv); err != nil {
		return err
	}

	// For now, we'll assume the league name is just "Bouncy League" or fetch it later
	return s.emailService.SendInvitation(email, token, "Bouncy League")
}

func (s *AuthService) GetLeagueInvitations(leagueID, requesterID string) ([]models.Invitation, error) {
	isAdmin, err := s.leagueRepo.IsAdmin(leagueID, requesterID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, fmt.Errorf("insufficient permissions to view invitations for this league")
	}

	return s.userRepo.GetInvitationsByLeague(leagueID)
}

func (s *AuthService) RegisterWithInvitation(token, name, email, password string) error {
	inv, err := s.userRepo.GetInvitationByToken(token)
	if err != nil || !inv.IsValid() {
		return ErrInvalidInvitation
	}

	// Registration email must match invitation email for security
	if inv.Email != email {
		return fmt.Errorf("registration email does not match invitation")
	}

	if err := s.RegisterWithRoles(name, email, password, inv.SpecialRole); err != nil {
		return err
	}

	// Mark invitation as used
	_ = s.userRepo.MarkInvitationAsUsed(token, time.Now())

	// Auto-claim any spreadsheet records
	user, _ := s.userRepo.GetUserByEmail(email)
	if user != nil {
		_ = s.paymentService.ClaimUnclaimedRecords(user.ID, name)
	}

	return nil
}

func (s *AuthService) Register(name, email, password string) error {
	return s.RegisterWithRoles(name, email, password, "")
}

func (s *AuthService) RegisterWithRoles(name, email, password, specialRole string) error {
	if err := s.ValidateLoginRequirements(email, password); err != nil {
		return err
	}

	_, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		slog.Info("user with that email already exists")
		return ErrUserAlreadyExists
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return errors.New("internal server error")
	}

	roles := []string{"user"}
	if specialRole != "" {
		roles = append(roles, specialRole)
	}

	newUser := &models.User{
		ID:           uuid.NewString(),
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		Roles:        roles,
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
