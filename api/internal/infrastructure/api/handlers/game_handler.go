package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/domain/models"
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
	r.Delete("/game/{gameId}", handler.CancelGame)
}

func (h *GameHandler) ListGames(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")

	games, err := h.service.GetGamesForLeague(leagueId)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, games)
}

func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")

	game, err := h.service.GetGameById(gameId)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, game)
}

type CreateGameRequest struct {
	Location    string `json:"location"`
	CostInCents int    `json:"costInCents"`
}

func (h *GameHandler) AddGame(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")

	var req CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	game := models.CreateGame(leagueId, req.Location, req.CostInCents)

	result, err := h.service.Create(game)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *GameHandler) CancelGame(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")

	if err := h.service.CancelGame(gameId); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{})
}
