package application

import (
	"github.com/bouncy/bouncy-api/internal/domain/interfaces"
	"github.com/bouncy/bouncy-api/internal/domain/models"
)

type GameService struct {
	repo interfaces.GameRepository
}

func NewGameService(repo interfaces.GameRepository) *GameService {
	return &GameService{repo: repo}
}

func (s *GameService) GetGamesForLeague(leagueId string) ([]*models.Game, error) {
	return s.repo.ListGamesByLeague(leagueId)
}

func (s *GameService) GetGameById(gameId string) (*models.Game, error) {
	return s.repo.GetById(gameId)
}

func (s *GameService) UpdateGame(game *models.Game) (*models.Game, error) {
	return s.repo.Update(game.ID, game)
}

func (s *GameService) Create(game *models.Game) (*models.Game, error) {
	return s.repo.Create(game)
}

func (s *GameService) CancelGame(leagueId string) error {
	return s.repo.Cancel(leagueId)
}
