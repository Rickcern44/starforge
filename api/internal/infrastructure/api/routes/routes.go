package routes

import (
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/dependencies"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	engine *gin.Engine,
	deps *dependencies.Dependencies,
) {
	api := engine.Group("/api")

	handlers.RegisterAuthRoutes(api, deps.AuthHandler)
	handlers.RegisterLeagueRoutes(api, deps.LeagueHandler)
	handlers.RegisterLeagueMemberHandlers(api, deps.LeagueMemberHandler)
	handlers.RegisterGameRoutes(api, deps.GameHandler)
}
