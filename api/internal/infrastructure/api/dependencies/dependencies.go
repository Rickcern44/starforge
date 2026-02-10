package dependencies

import (
	"net/http" // New import for http.Handler

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/application/game_attendances"
	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/bouncy/bouncy-api/internal/application/payments"
	"github.com/bouncy/bouncy-api/internal/application/users"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/handlers"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/middleware" // New import for middleware
)

type Dependencies struct {
	LeagueHandler         *handlers.LeagueHandler
	LeagueMemberHandler   *handlers.LeagueMemberHandler
	GameHandler           *handlers.GameHandler
	AuthHandler           *handlers.AuthHandler
	UserHandler           *handlers.UserHandler
	GameAttendanceHandler *handlers.GameAttendanceHandler
	PaymentsHandler       *handlers.PaymentsHandler
	AuthMiddleware        func(next http.Handler) http.Handler
}

func BuildDependencies(
	leagueService *leagues.LeagueService,
	authService *application.AuthService,
	leagueMemberService *leagues.LeagueMemberService,
	gameService *application.GameService,
	jwtService *application.JwtService,
	userService *users.Service,
	gameAttendanceService *game_attendances.Service,
	paymentsService *payments.Service,
) *Dependencies {
	return &Dependencies{
		LeagueHandler:         handlers.NewLeagueHandler(leagueService),
		LeagueMemberHandler:   handlers.NewLeagueMemberHandler(leagueMemberService),
		GameHandler:           handlers.NewGameHandler(gameService),
		AuthHandler:           handlers.NewAuthHandler(authService),
		UserHandler:           handlers.NewUserHandler(userService),
		GameAttendanceHandler: handlers.NewGameAttendanceHandler(gameAttendanceService),
		PaymentsHandler:       handlers.NewPaymentsHandler(paymentsService),
		AuthMiddleware:        middleware.AuthMiddleware(jwtService),
	}
}
