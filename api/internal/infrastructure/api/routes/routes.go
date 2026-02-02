package routes

import (
	"net/http"

	"github.com/bouncy/bouncy-api/internal/infrastructure/api/dependencies"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/handlers"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(
	router *chi.Mux,
	deps *dependencies.Dependencies,
) {
	router.Route("/api", func(r chi.Router) {
		r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handlers.RegisterDocsEndpoints(r, "./.docs")

		handlers.RegisterAuthRoutes(r, deps.AuthHandler)

		r.Group(func(r chi.Router) {
			r.Use(deps.AuthMiddleware)

			handlers.RegisterLeagueRoutes(r, deps.LeagueHandler)
			handlers.RegisterLeagueMemberHandlers(r, deps.LeagueMemberHandler)
			handlers.RegisterGameRoutes(r, deps.GameHandler)
			handlers.RegisterUserRoutes(r, deps.UserHandler)
			handlers.RegisterGameAttendanceRoutes(r, deps.GameAttendanceHandler)

			r.Group(func(r chi.Router) {
				r.Use(middleware.RoleMiddleware("admin"))
				handlers.RegisterPaymentsRoutes(r, deps.PaymentsHandler)
				// Admin-protected user routes
				r.Group(func(r chi.Router) {
					handlers.RegisterAdminUserRoutes(r, deps.UserHandler)
				})
			})
		})
	})
}
