package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/go-chi/chi/v5"
)

type LeagueHandler struct {
	service *leagues.LeagueService
}

func NewLeagueHandler(service *leagues.LeagueService) *LeagueHandler {
	return &LeagueHandler{service: service}
}

func RegisterLeagueRoutes(r chi.Router, handler *LeagueHandler) {
	r.Post("/league", handler.CreateLeague)
	r.Get("/league/{leagueId}", handler.GetLeague)
	r.Delete("/league/{leagueId}", handler.Delete)
}

func (h *LeagueHandler) GetLeague(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "leagueId")

	game, err := h.service.GetLeague(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, game)
}

type createLeagueRequest struct {
	Name string `json:"name"`
}

func (h *LeagueHandler) CreateLeague(w http.ResponseWriter, r *http.Request) {
	var req createLeagueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	league, err := h.service.CreateLeague(
		r.Context(),
		req.Name,
	)
	if err != nil {
		writeJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, league)
}

func (h *LeagueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "leagueId")

	if err := h.service.Delete(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
