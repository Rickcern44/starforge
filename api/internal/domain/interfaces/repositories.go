package interfaces

import "github.com/bouncy/bouncy-api/internal/domain/models"

type LeagueRepository interface {
	GetById(league string) (*models.League, error)
	Save(league *models.League) error
	Delete(league string) error
}

type GameRepository interface {
	ListGamesByLeague(leagueId string) ([]*models.Game, error)
	Create(game *models.Game) (*models.Game, error)
	GetById(gameId string) (*models.Game, error)
	Update(gameId string, game *models.Game) (*models.Game, error)
	Cancel(gameId string) error
}
