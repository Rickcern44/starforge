package container

import (
	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/application/game_attendances"
	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/bouncy/bouncy-api/internal/application/payments"
	"github.com/bouncy/bouncy-api/internal/application/users"
	"github.com/bouncy/bouncy-api/internal/domain/interfaces"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/dependencies"
	"github.com/bouncy/bouncy-api/internal/infrastructure/config"
	"github.com/bouncy/bouncy-api/internal/infrastructure/email"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence/repositories"
	"gorm.io/gorm"
)

type AppContainer struct {
	DB                    *gorm.DB
	Settings              *config.Config
	LeagueService         *leagues.LeagueService
	LeagueMemberService   *leagues.LeagueMemberService
	GameService           *application.GameService
	UserService           *users.Service
	GameAttendanceService *game_attendances.Service
	PaymentsService       *payments.Service
	JwtService            *application.JwtService
	AuthService           *application.AuthService
	EmailService          interfaces.EmailService
}

func NewAppContainer(db *gorm.DB, settings *config.Config) *AppContainer {
	leagueRepo := repositories.NewLeagueRepository(db)
	leagueMemberRepo := repositories.NewLeagueMemberRepository(db)
	gameRepo := repositories.NewGameRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	gameAttendanceRepo := repositories.NewGameAttendanceRepository(db)
	paymentsRepo := repositories.NewPaymentsRepository(db)

	leagueService := leagues.NewLeagueService(leagueRepo)
	leagueMemberService := leagues.NewLeagueMemberService(leagueMemberRepo)
	gameService := application.NewGameService(gameRepo)
	userService := users.NewUserService(authRepo)
	gameAttendanceService := game_attendances.NewGameAttendanceService(gameAttendanceRepo, gameRepo, paymentsRepo)
	paymentsService := payments.NewPaymentsService(paymentsRepo)

	emailService := email.NewMockEmailService()

	jwtService := application.NewJwtService(settings.Auth.JwtSecret)
	authService := application.NewAuthService(jwtService, authRepo, emailService, paymentsService, leagueMemberRepo)

	return &AppContainer{
		DB:                    db,
		Settings:              settings,
		LeagueService:         leagueService,
		LeagueMemberService:   leagueMemberService,
		GameService:           gameService,
		UserService:           userService,
		GameAttendanceService: gameAttendanceService,
		PaymentsService:       paymentsService,
		JwtService:            jwtService,
		AuthService:           authService,
		EmailService:          emailService,
	}
}

func (c *AppContainer) ToDependencies() *dependencies.Dependencies {
	return dependencies.BuildDependencies(
		c.LeagueService,
		c.AuthService,
		c.LeagueMemberService,
		c.GameService,
		c.JwtService,
		c.UserService,
		c.GameAttendanceService,
		c.PaymentsService,
	)
}
