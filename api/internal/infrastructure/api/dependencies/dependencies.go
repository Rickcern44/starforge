package dependencies

import (
	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/handlers"
)

type Dependencies struct {
	LeagueHandler       *handlers.LeagueHandler
	LeagueMemberHandler *handlers.LeagueMemberHandler
	AuthHandler         *handlers.AuthHandler
}

func BuildDependencies(
	leagueService *leagues.LeagueService,
	authService *application.AuthService,
	leagueMemberService *leagues.LeagueMemberService,
) *Dependencies {
	return &Dependencies{
		LeagueHandler:       handlers.NewLeagueHandler(leagueService),
		LeagueMemberHandler: handlers.NewLeagueMemberHandler(leagueMemberService),
		AuthHandler:         handlers.NewAuthHandler(authService),
	}
}
