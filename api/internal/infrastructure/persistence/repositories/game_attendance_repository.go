package repositories

import (
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence/mappers"
	"gorm.io/gorm"
)

type GameAttendanceRepository struct {
	db *gorm.DB
}

func NewGameAttendanceRepository(db *gorm.DB) *GameAttendanceRepository {
	return &GameAttendanceRepository{db: db}
}

func (r *GameAttendanceRepository) Add(attendance *models.GameAttendance, gameId string) error {
	attendanceDto := mappers.GameAttendanceToDto(attendance, gameId)
	return r.db.Create(&attendanceDto).Error
}

func (r *GameAttendanceRepository) Remove(gameID, userID string) error {
	return r.db.Where("game_id = ? AND user_id = ?", gameID, userID).Delete(&persistence.GameAttendance{}).Error
}
