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

func (s *Service) Add(status models.AttendanceStatus, gameID, userId, comment string) (bool, error) {
	existingAttendance, err := s.repo.FindByGameAndUser(gameID, userId)
	if err != nil {
		return false, err
	}

	if existingAttendance != nil {
		existingAttendance.Status = status
		existingAttendance.CheckInComment = comment
		existingAttendance.UpdatedAt = time.Now()
		err := s.repo.Update(existingAttendance, gameID)
		return false, err
	}

	gameAttendance := &models.GameAttendance{
		UserID:         userId,
		CheckedIn:      true,
		Status:         status,
		CheckInComment: comment,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	err = s.repo.Add(gameAttendance, gameID)
	return true, err
}

func (s *Service) Remove(gameID, userID string) error {
	return s.repo.Remove(gameID, userID)
}
