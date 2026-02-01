package dependencies

import (
	"net/http" // New import for http.Handler

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/handlers"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/middleware" // New import for middleware
)

type Dependencies struct {
	LeagueHandler       *handlers.LeagueHandler
	LeagueMemberHandler *handlers.LeagueMemberHandler
	GameHandler         *handlers.GameHandler
	AuthHandler         *handlers.AuthHandler
	AuthMiddleware      func(next http.Handler) http.Handler
}

func BuildDependencies(
	leagueService *leagues.LeagueService,
	authService *application.AuthService,
	leagueMemberService *leagues.LeagueMemberService,
	gameService *application.GameService,
	jwtService *application.JwtService,
) *Dependencies {
	return &Dependencies{
		LeagueHandler:       handlers.NewLeagueHandler(leagueService),
		LeagueMemberHandler: handlers.NewLeagueMemberHandler(leagueMemberService),
		GameHandler:         handlers.NewGameHandler(gameService),
		AuthHandler:         handlers.NewAuthHandler(authService),
		AuthMiddleware:      middleware.AuthMiddleware(jwtService),
	}
}
