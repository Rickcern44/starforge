package game_attendances

import (
	"github.com/bouncy/bouncy-api/internal/application/interfaces"
)

type Service struct {
	repo interfaces.GameAttendanceRepository
}

func NewGameAttendanceService(repo interfaces.GameAttendanceRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Add(gameID, userID string) error {
	return s.repo.Add(gameID, userID)
}

func (s *Service) Remove(gameID, userID string) error {
	return s.repo.Remove(gameID, userID)
}
