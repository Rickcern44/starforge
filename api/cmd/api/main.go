package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bouncy/bouncy-api/internal/infrastructure/api/routes"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/server"
	"github.com/bouncy/bouncy-api/internal/infrastructure/config"
	"github.com/bouncy/bouncy-api/internal/infrastructure/container"
	"github.com/bouncy/bouncy-api/internal/infrastructure/database"
)

// @title           Bouncy API
// @version         0.1.0
// @description     Internal API
// @BasePath        /api/v1
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:3000
func main() {
	settings, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting app", "env", settings.AppEnv, "isDev", settings.IsDevelopment)

	chiServer := server.NewServer()
	dbServer := database.NewDatabaseService(settings)
	
	slog.Info("Connecting to database...")
	if err := dbServer.Connect(); err != nil {
		slog.Error("Database connection failed", "error", err)
		os.Exit(1)
	}

	if settings.IsDevelopment {
		slog.Info("Development mode: running migrations and seeders")
		_ = dbServer.UpdateDatabase()
		_ = dbServer.Seed()
	}

	app := container.NewAppContainer(dbServer.Database, settings)
	routes.RegisterRoutes(chiServer.Router(), app.ToDependencies())

	addr := fmt.Sprintf("0.0.0.0:%v", settings.Server.Port)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine (so we can wait for shutdown)
	go func() {
		slog.Info("Server is listening", "addr", addr)
		if err := chiServer.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	<-shutdown

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := chiServer.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited properly")
}
