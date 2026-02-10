package interfaces

import "github.com/bouncy/bouncy-api/internal/domain/models"

type GameAttendanceRepository interface {
	Add(attendance *models.GameAttendance, gameID string) error
	Remove(gameID, userID string) error
}
