package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/application/game_attendances"
	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/bouncy/bouncy-api/internal/application/payments"
	"github.com/bouncy/bouncy-api/internal/application/users"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/dependencies"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/routes"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/server"
	"github.com/bouncy/bouncy-api/internal/infrastructure/config"
	"github.com/bouncy/bouncy-api/internal/infrastructure/database"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence/repositories"
	"gorm.io/gorm"
)

// @title 		Bouncy API
// @version 	1.0

func main() {
	var settings *config.Config
	settings, err := config.LoadConfig(os.Getenv("APP_CONFIG_PATH"))

	if err != nil {
		log.Fatal(err)
	}

	chiServer := server.NewServer()
	dbServer := database.NewDatabaseService(settings)

	if err := dbServer.Connect(); err != nil {
		log.Fatal(err)
	}

	_ = dbServer.UpdateDatabase()

	deps := BuildApplication(dbServer.Database, settings)
	routes.RegisterRoutes(chiServer.Router(), deps)

	// Start server in a goroutine
	go func() {
		port := fmt.Sprintf(":%v", settings.Server.Port)
		log.Printf("Starting Server on %s", port)
		if err := chiServer.Start(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := chiServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited properly")
}

func BuildApplication(db *gorm.DB, settings *config.Config) *dependencies.Dependencies {

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
	gameAttendanceService := game_attendances.NewGameAttendanceService(gameAttendanceRepo)
	paymentsService := payments.NewPaymentsService(paymentsRepo)
	// Auth services
	jwtService := application.NewJwtService(settings.Auth.JwtSecret)
	authService := application.NewAuthService(jwtService, authRepo)

	return dependencies.BuildDependencies(leagueService, authService, leagueMemberService, gameService, jwtService, userService, gameAttendanceService, paymentsService)
}
