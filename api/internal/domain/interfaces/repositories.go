package interfaces

import "github.com/bouncy/bouncy-api/internal/domain/models"

type LeagueRepository interface {
	GetById(league string) (*models.League, error)
	Save(league *models.League) error
	Delete(league string) error
}

type GameRepository interface {
	GetById(id string) (*models.Game, error)
	Save(game *models.Game) error
	Delete(id string) error
}
