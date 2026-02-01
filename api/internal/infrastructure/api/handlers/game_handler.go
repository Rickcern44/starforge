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

func RegisterGameRoutes(r chi.Router, handler *GameHandler) {
	r.Get("/league/{leagueId}/games", handler.ListGames)
	r.Post("/league/{leagueId}/games", handler.AddGame)
	r.Get("/game/{gameId}", handler.GetGame)
	r.Put("/game/{gameId}", handler.UpdateGame)
	r.Delete("/game/{gameId}", handler.CancelGame)
}

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

func (h *GameHandler) ListGames(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")

	games, err := h.service.GetGamesForLeague(leagueId)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, games)
}

func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")

	game, err := h.service.GetGameById(gameId)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, game)
}

type CreateGameRequest struct {
	Location    string `json:"location"`
	CostInCents int    `json:"costInCents"`
}

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

func (h *GameHandler) CancelGame(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")

	if err := h.service.CancelGame(gameId); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{})
}
