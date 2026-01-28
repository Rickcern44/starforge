package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/dependencies"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/routes"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/server"
	"github.com/bouncy/bouncy-api/internal/infrastructure/database"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence/repositories"
	"gorm.io/gorm"
)

// @title 		Bouncy API
// @version 	1.0

func main() {
	ginServer := server.NewServer()
	dbServer := database.NewDatabaseService()

	if err := dbServer.Connect(); err != nil {
		log.Fatal(err)
	}

	_ = dbServer.UpdateDatabase()

	deps := BuildApplication(dbServer.Database)
	routes.RegisterRoutes(ginServer.Engine(), deps)

	// Start server in a goroutine
	go func() {
		log.Println("Starting Server on :3000")
		if err := ginServer.Start(":3000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	if err := ginServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited properly")
}

func BuildApplication(db *gorm.DB) *dependencies.Dependencies {

	leagueRepo := repositories.NewLeagueRepository(db)

	leagueService := leagues.NewLeagueService(leagueRepo)
	// Auth services
	jwtService := application.NewJwtService(os.Getenv("JWT_TOKEN"))
	authService := application.NewAuthService(jwtService, time.Hour*24)

	return dependencies.BuildDependencies(leagueService, authService)
}
