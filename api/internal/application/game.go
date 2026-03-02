package application

import (
	"time"

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

func (s *GameService) CreateRecurringGames(leagueId, location string, costInCents int, startTime time.Time, interval string, count int) ([]*models.Game, error) {
	var games []*models.Game
	currentTime := startTime

	for i := 0; i < count; i++ {
		game := models.CreateGameFromData(leagueId, location, costInCents, currentTime)
		result, err := s.repo.Create(game)
		if err != nil {
			return games, err
		}
		games = append(games, result)

		switch interval {
		case "Daily":
			currentTime = currentTime.AddDate(0, 0, 1)
		case "Weekly":
			currentTime = currentTime.AddDate(0, 0, 7)
		case "Bi-weekly":
			currentTime = currentTime.AddDate(0, 0, 14)
		case "Monthly":
			currentTime = currentTime.AddDate(0, 1, 0)
		default:
			return games, nil
		}
	}

	return games, nil
}

func (s *GameService) CancelGame(leagueId string) error {
	return s.repo.Cancel(leagueId)
}
