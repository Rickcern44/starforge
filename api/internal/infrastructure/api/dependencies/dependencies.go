package dependencies

import (
	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/handlers"
)

type Dependencies struct {
	LeagueHandler *handlers.LeagueHandler
	AuthHandler   *handlers.AuthHandler
}

func BuildDependencies(
	leagueService *leagues.LeagueService,
	authService *application.AuthService,
) *Dependencies {
	return &Dependencies{
		LeagueHandler: handlers.NewLeagueHandler(leagueService),
		AuthHandler:   handlers.NewAuthHandler(authService),
	}
}
