package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/bouncy/bouncy-api/internal/application/features"
	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/contract"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/middleware"
	"github.com/bouncy/bouncy-api/internal/infrastructure/auth"
	"github.com/bouncy/bouncy-api/internal/infrastructure/utils"
)

type LeagueHandler struct {
	service        *leagues.LeagueService
	featureService *features.FeatureFlagService
}

func NewLeagueHandler(service *leagues.LeagueService, featureService *features.FeatureFlagService) *LeagueHandler {
	return &LeagueHandler{
		service:        service,
		featureService: featureService,
	}
}

// RegisterLeagueRoutes registers the league related routes.
func RegisterLeagueRoutes(r chi.Router, handler *LeagueHandler) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RoleMiddleware("league_creator"))
		r.Post("/league", handler.CreateLeague)
	})
	r.Get("/league/{leagueId}", handler.GetLeague)
	r.Get("/me/leagues", handler.GetLeaguesForUser)
	r.Delete("/league/{leagueId}", handler.Delete)
}

// GetLeaguesForUser handles getting all leagues for the authenticated user.
// @Summary Get all leagues for the authenticated user
// @Description Retrieves all leagues for the currently authenticated user.
// @Tags leagues
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.League "A list of leagues"
// @Failure 401 {object} contract.ErrorResponse "Unauthorized"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /api/v1/me/leagues [get]
func (h *LeagueHandler) GetLeaguesForUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.Claims)
	if !ok || claims == nil {
		slog.Error("Get leagues for user failed: authentication required")
		utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "authentication required"})
		return
	}

	userLeagues, err := h.service.GetLeaguesForUser(claims.UserId)
	if err != nil {
		slog.Error("Get leagues for user service failed", "userId", claims.UserId, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, userLeagues)
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
// @Router /api/v1/league/{leagueId} [get]
func (h *LeagueHandler) GetLeague(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "leagueId")

	game, err := h.service.GetLeague(id)
	if err != nil {
		slog.Error("Get league failed", "leagueId", id, "error", err)
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
// @Router /api/v1/league [post]
func (h *LeagueHandler) CreateLeague(w http.ResponseWriter, r *http.Request) {
	if !h.featureService.IsEnabled("league_creation") {
		utils.WriteJSON(w, http.StatusForbidden, contract.ErrorResponse{Error: "League creation is currently disabled"})
		return
	}

	var req createLeagueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode create league request", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: err.Error()})
		return
	}

	claims, _ := r.Context().Value(auth.ClaimsContextKey).(*auth.Claims)

	league, err := h.service.CreateLeague(
		r.Context(),
		req.Name,
		claims.UserId,
	)
	if err != nil {
		slog.Error("Create league service failed", "name", req.Name, "error", err)
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
// @Router /api/v1/league/{leagueId} [delete]
func (h *LeagueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "leagueId")

	if err := h.service.Delete(id); err != nil {
		slog.Error("Delete league failed", "leagueId", id, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
