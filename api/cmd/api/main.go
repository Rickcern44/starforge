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
		log.Fatal(err)
	}

	chiServer := server.NewServer()
	dbServer := database.NewDatabaseService(settings)
	if err := dbServer.Connect(); err != nil {
		log.Fatal(err)
	}

	if settings.IsDevelopment {
		_ = dbServer.UpdateDatabase()
		_ = dbServer.Seed()
	}


	app := container.NewAppContainer(dbServer.Database, settings)
	routes.RegisterRoutes(chiServer.Router(), app.ToDependencies())

	addr := fmt.Sprintf(":%v", settings.Server.Port)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine (so we can wait for shutdown)
	go func() {
		log.Printf("Starting Server on %s", addr)
		if err := chiServer.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-shutdown

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := chiServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %+v", err)
	}

	log.Println("Server exited properly")
}
