package routes

import (
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/dependencies"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(
	router *chi.Mux,
	deps *dependencies.Dependencies,
) {
	router.Route("/api", func(r chi.Router) {
		handlers.RegisterAuthRoutes(r, deps.AuthHandler)
		handlers.RegisterLeagueRoutes(r, deps.LeagueHandler)
		handlers.RegisterLeagueMemberHandlers(r, deps.LeagueMemberHandler)
		handlers.RegisterGameRoutes(r, deps.GameHandler)
	})
}
