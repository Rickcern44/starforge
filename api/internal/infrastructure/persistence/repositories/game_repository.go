package repositories

import (
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence/mappers"
	"gorm.io/gorm"
)

type GameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{db: db}
}

func (r GameRepository) ListGamesByLeague(leagueId string) ([]*models.Game, error) {
	var rows []persistence.Game
	if err := r.db.Preload("Attendance.User").Where("league_id = ?", leagueId).Find(&rows).Error; err != nil {
		return nil, err
	}

	var games []*models.Game
	for _, row := range rows {
		games = append(games, mappers.GameToDomain(row))
	}
	return games, nil
}

func (r GameRepository) Create(game *models.Game) (*models.Game, error) {
	pGame := mappers.GameToPersistence(game)

	if err := r.db.Create(pGame).Error; err != nil {
		return nil, err
	}

	return mappers.GameToDomain(*pGame), nil
}

func (r GameRepository) GetById(gameId string) (*models.Game, error) {
	var row persistence.Game
	if err := r.db.Preload("Attendance.User").Where("id = ?", gameId).First(&row).Error; err != nil {
		return nil, err
	}

	return mappers.GameToDomain(row), nil
}

func (r GameRepository) Update(gameId string, game *models.Game) (*models.Game, error) {
	pGame := mappers.GameToPersistence(game)

	if err := r.db.Model(&persistence.Game{}).Where("id = ?", gameId).Updates(pGame).Error; err != nil {
		return nil, err
	}

	return game, nil
}

func (r GameRepository) Cancel(gameId string) error {
	return r.db.Delete(&persistence.Game{}, "id = ?", gameId).Error
}
