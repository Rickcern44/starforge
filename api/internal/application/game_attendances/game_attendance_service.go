package game_attendances

import (
	"time"

	"github.com/bouncy/bouncy-api/internal/application/interfaces"
	"github.com/bouncy/bouncy-api/internal/domain/models"
)

type Service struct {
	repo interfaces.GameAttendanceRepository
}

func NewGameAttendanceService(repo interfaces.GameAttendanceRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Add(status models.AttendanceStatus, gameID, userId, comment string) error {
	gameAttendance := &models.GameAttendance{
		UserID:         userId,
		CheckedIn:      true,
		Status:         status,
		CheckInComment: comment,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	return s.repo.Add(gameAttendance, gameID)
}

func (s *Service) Remove(gameID, userID string) error {
	return s.repo.Remove(gameID, userID)
}
