package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/contract"
	"github.com/bouncy/bouncy-api/internal/infrastructure/utils"
	"github.com/go-chi/chi/v5"
)

type LeagueHandler struct {
	service *leagues.LeagueService
}

func NewLeagueHandler(service *leagues.LeagueService) *LeagueHandler {
	return &LeagueHandler{service: service}
}

// RegisterLeagueRoutes registers the league related routes.
func RegisterLeagueRoutes(r chi.Router, handler *LeagueHandler) {
	r.Post("/league", handler.CreateLeague)
	r.Get("/league/{leagueId}", handler.GetLeague)
	r.Delete("/league/{leagueId}", handler.Delete)
}

// GetLeague handles getting a league by its ID.
// @Summary Get a league by its ID
// @Description Retrieves details of a single league using its league ID.
// @Tags leagues
// @Accept json
// @Produce json
// @Param leagueId path string true "ID of the league"
// @Security BearerAuth
// @Success 200 {object} models.League "League details"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /league/{leagueId} [get]
func (h *LeagueHandler) GetLeague(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "leagueId")

	game, err := h.service.GetLeague(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, game)
}

// createLeagueRequest represents the request body for creating a new league.
type createLeagueRequest struct {
	Name string `json:"name"`
}

// CreateLeague handles creating a new league.
// @Summary Create a new league
// @Description Creates a new league with the provided name.
// @Tags leagues
// @Accept json
// @Produce json
// @Param request body createLeagueRequest true "League creation details"
// @Security BearerAuth
// @Success 201 {object} models.League "League created successfully"
// @Failure 400 {object} contract.ErrorResponse "Invalid request body"
// @Failure 409 {object} contract.ErrorResponse "League with given name already exists"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /league [post]
func (h *LeagueHandler) CreateLeague(w http.ResponseWriter, r *http.Request) {
	var req createLeagueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: err.Error()})
		return
	}

	league, err := h.service.CreateLeague(
		r.Context(),
		req.Name,
	)
	if err != nil {
		utils.WriteJSON(w, http.StatusConflict, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, league)
}

// Delete handles deleting a league.
// @Summary Delete a league
// @Description Deletes a league with the given ID.
// @Tags leagues
// @Accept json
// @Produce json
// @Param leagueId path string true "ID of the league to delete"
// @Security BearerAuth
// @Success 200 {object} object "League deleted successfully"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /league/{leagueId} [delete]
func (h *LeagueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "leagueId")

	if err := h.service.Delete(id); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
