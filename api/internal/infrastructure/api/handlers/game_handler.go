package handlers

import (
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	service *application.GameService
}

func NewGameHandler(service *application.GameService) *GameHandler {
	return &GameHandler{service: service}
}

func RegisterGameRoutes(rg *gin.RouterGroup, handler *GameHandler) {
	leagueGroup := rg.Group("/league")
	gameGroup := rg.Group("/game")

	leagueGroup.GET("/:leagueId/games", handler.ListGames)
	leagueGroup.POST("/:leagueId/games", handler.AddGame)

	gameGroup.GET("/:gameId", handler.GetGame)
	// Add an update here eventually
	gameGroup.DELETE("/:gameId", handler.CancelGame)
}

func (h *GameHandler) ListGames(c *gin.Context) {
	leagueId := c.Param("leagueId")

	games, err := h.service.GetGamesForLeague(leagueId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, games)
}

func (h *GameHandler) GetGame(c *gin.Context) {
	gameId := c.Param("gameId")

	game, err := h.service.GetGameById(gameId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, game)
}

type CreateGameRequest struct {
	location    string
	costInCents int
}

func (h *GameHandler) AddGame(c *gin.Context) {
	leagueId := c.Param("leagueId")

	var req CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	game := models.CreateGame(leagueId, req.location, req.costInCents)

	result, err := h.service.Create(game)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func (h *GameHandler) CancelGame(c *gin.Context) {
	gameId := c.Param("gameId")

	if err := h.service.CancelGame(gameId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{})
}
