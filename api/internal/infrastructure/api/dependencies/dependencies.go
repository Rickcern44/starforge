package dependencies

import (
	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/handlers"
)

type Dependencies struct {
	LeagueHandler *handlers.LeagueHandler
}

func BuildDependencies(
	leagueService *application.LeagueService,
) *Dependencies {
	return &Dependencies{
		LeagueHandler: handlers.NewLeagueHandler(leagueService),
	}
}
