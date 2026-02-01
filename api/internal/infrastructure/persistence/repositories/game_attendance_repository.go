package repositories

import (
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

type GameAttendanceRepository struct {
	db *gorm.DB
}

func NewGameAttendanceRepository(db *gorm.DB) *GameAttendanceRepository {
	return &GameAttendanceRepository{db: db}
}

func (r *GameAttendanceRepository) Add(gameID, userID string) error {
	attendance := persistence.GameAttendance{
		GameID: gameID,
		UserID: userID,
	}
	return r.db.Create(&attendance).Error
}

func (r *GameAttendanceRepository) Remove(gameID, userID string) error {
	return r.db.Where("game_id = ? AND user_id = ?", gameID, userID).Delete(&persistence.GameAttendance{}).Error
}
