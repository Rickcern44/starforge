package email

import (
	"fmt"
	"log/slog"
	"strings"
)

type MockEmailService struct{}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}

func (s *MockEmailService) SendInvitation(email, token, leagueName string) error {
	link := fmt.Sprintf("https://bouncy.app/register?token=%s", token)
	slog.Info("EMAIL SENT (INVITATION)", "to", email, "league", leagueName, "link", link)
	return nil
}

func (s *MockEmailService) SendGameNotification(emails []string, gameDate, location string, signupLink string) error {
	to := strings.Join(emails, ", ")
	slog.Info("EMAIL SENT (GAME NOTIFICATION)", "to", to, "date", gameDate, "location", location, "signup", signupLink)
	return nil
}
