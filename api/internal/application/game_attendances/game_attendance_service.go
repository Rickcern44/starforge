package game_attendances

import (
	"log/slog"
	"time"

	"github.com/bouncy/bouncy-api/internal/domain/interfaces"
	"github.com/bouncy/bouncy-api/internal/domain/models"
)

type Service struct {
	repo         interfaces.GameAttendanceRepository
	gameRepo     interfaces.GameRepository
	paymentsRepo interfaces.PaymentsRepository
}

func NewGameAttendanceService(repo interfaces.GameAttendanceRepository, gameRepo interfaces.GameRepository, paymentsRepo interfaces.PaymentsRepository) *Service {
	return &Service{
		repo:         repo,
		gameRepo:     gameRepo,
		paymentsRepo: paymentsRepo,
	}
}

func (s *Service) Add(status models.AttendanceStatus, gameID, userId, comment string) (bool, error) {
	game, err := s.gameRepo.GetById(gameID)
	if err != nil {
		return false, err
	}

	existingAttendance, err := s.repo.FindByGameAndUser(gameID, userId)
	if err != nil {
		return false, err
	}

	isNew := existingAttendance == nil

	if existingAttendance != nil {
		existingAttendance.Status = status
		existingAttendance.CheckInComment = comment
		existingAttendance.UpdatedAt = time.Now()
		if err := s.repo.Update(existingAttendance, gameID); err != nil {
			return false, err
		}
	} else {
		gameAttendance := &models.GameAttendance{
			UserID:         userId,
			CheckedIn:      true,
			Status:         status,
			CheckInComment: comment,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := s.repo.Add(gameAttendance, gameID); err != nil {
			return false, err
		}
	}

	// Create a charge if the status is 'Yes' (0) and cost is > 0
	if status == models.Yes && game.CostInCents > 0 {
		// Check if charge already exists
		existingCharges, _ := s.paymentsRepo.ListChargesByUser(userId)
		hasCharge := false
		for _, c := range existingCharges {
			if c.GameID == gameID {
				hasCharge = true
				break
			}
		}

		if !hasCharge {
			charge := models.CreateGameCharge(gameID, &userId, "", game.CostInCents)
			if err := s.paymentsRepo.CreateCharge(charge); err != nil {
				slog.Error("Failed to create automatic charge for RSVP", "gameId", gameID, "userId", userId, "error", err)
				// We don't fail the RSVP if charge creation fails
				return isNew, nil
			}
			slog.Info("Automatically created charge for RSVP", "gameId", gameID, "userId", userId, "amount", game.CostInCents)
		}
	}

	return isNew, nil
}

func (s *Service) Remove(gameID, userID string) error {
	return s.repo.Remove(gameID, userID)
}
