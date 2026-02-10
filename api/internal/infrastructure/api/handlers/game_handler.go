package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/contract"
	"github.com/bouncy/bouncy-api/internal/infrastructure/utils"
	"github.com/go-chi/chi/v5"
)

type GameHandler struct {
	service *application.GameService
}

func NewGameHandler(service *application.GameService) *GameHandler {
	return &GameHandler{service: service}
}

// RegisterGameRoutes registers the game related routes.
func RegisterGameRoutes(r chi.Router, handler *GameHandler) {
	r.Get("/league/{leagueId}/games", handler.ListGames)
	r.Post("/league/{leagueId}/games", handler.AddGame)
	r.Get("/game/{gameId}", handler.GetGame)
	r.Put("/game/{gameId}", handler.UpdateGame)
	r.Delete("/game/{gameId}", handler.CancelGame)
}

// UpdateGame handles updating an existing game.
// @Summary Update an existing game
// @Description Updates the details of an existing game.
// @Tags games
// @Accept json
// @Produce json
// @Param gameId path string true "ID of the game to update"
// @Param request body models.Game true "Game object that needs to be updated"
// @Security BearerAuth
// @Success 200 {object} models.Game "Game updated successfully"
// @Failure 400 {object} contract.ErrorResponse "Invalid request body"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /game/{gameId} [put]
func (h *GameHandler) UpdateGame(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")

	var game models.Game
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: err.Error()})
		return
	}

	game.ID = gameId

	updatedGame, err := h.service.UpdateGame(&game)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedGame)
}

// ListGames handles listing all games for a specific league.
// @Summary List all games for a specific league
// @Description Retrieves a list of all games associated with the given league ID.
// @Tags games
// @Accept json
// @Produce json
// @Param leagueId path string true "ID of the league"
// @Security BearerAuth
// @Success 200 {array} models.Game "List of games"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /league/{leagueId}/games [get]
func (h *GameHandler) ListGames(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")

	games, err := h.service.GetGamesForLeague(leagueId)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, games)
}

// GetGame handles getting a game by its ID.
// @Summary Get a game by its ID
// @Description Retrieves details of a single game using its game ID.
// @Tags games
// @Accept json
// @Produce json
// @Param gameId path string true "ID of the game"
// @Security BearerAuth
// @Success 200 {object} models.Game "Game details"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /game/{gameId} [get]
func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")

	game, err := h.service.GetGameById(gameId)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, game)
}

// CreateGameRequest represents the request body for creating a new game.
type CreateGameRequest struct {
	Location    string `json:"location"`
	CostInCents int    `json:"costInCents"`
}

// AddGame handles adding a new game to a league.
// @Summary Add a new game to a league
// @Description Creates a new game entry for the specified league.
// @Tags games
// @Accept json
// @Produce json
// @Param leagueId path string true "ID of the league"
// @Param request body CreateGameRequest true "Game creation details"
// @Security BearerAuth
// @Success 200 {object} models.Game "Game created successfully"
// @Failure 400 {object} contract.ErrorResponse "Invalid request body"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /league/{leagueId}/games [post]
func (h *GameHandler) AddGame(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")

	var req CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: err.Error()})
		return
	}

	game := models.CreateGame(leagueId, req.Location, req.CostInCents)

	result, err := h.service.Create(game)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, result)
}

// CancelGame handles cancelling a game.
// @Summary Cancel a game
// @Description Marks a game as cancelled.
// @Tags games
// @Accept json
// @Produce json
// @Param gameId path string true "ID of the game to cancel"
// @Security BearerAuth
// @Success 200 {object} object "Game cancelled successfully"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /game/{gameId} [delete]
func (h *GameHandler) CancelGame(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")

	if err := h.service.CancelGame(gameId); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{})
}
