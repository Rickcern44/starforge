package dependencies

import (
	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/handlers"
)

type Dependencies struct {
	LeagueHandler       *handlers.LeagueHandler
	LeagueMemberHandler *handlers.LeagueMemberHandler
	GameHandler         *handlers.GameHandler
	AuthHandler         *handlers.AuthHandler
}

func BuildDependencies(
	leagueService *leagues.LeagueService,
	authService *application.AuthService,
	leagueMemberService *leagues.LeagueMemberService,
	gameService *application.GameService,
) *Dependencies {
	return &Dependencies{
		LeagueHandler:       handlers.NewLeagueHandler(leagueService),
		LeagueMemberHandler: handlers.NewLeagueMemberHandler(leagueMemberService),
		GameHandler:         handlers.NewGameHandler(gameService),
		AuthHandler:         handlers.NewAuthHandler(authService),
	}
}
