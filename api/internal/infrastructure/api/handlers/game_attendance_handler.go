package handlers

import (
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application/game_attendances"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/contract"
	"github.com/bouncy/bouncy-api/internal/infrastructure/auth"
	"github.com/bouncy/bouncy-api/internal/infrastructure/utils"
	"github.com/go-chi/chi/v5"
)

type GameAttendanceHandler struct {
	service *game_attendances.Service
}

func NewGameAttendanceHandler(service *game_attendances.Service) *GameAttendanceHandler {
	return &GameAttendanceHandler{service: service}
}

// RegisterGameAttendanceRoutes registers the game attendance related routes.
func RegisterGameAttendanceRoutes(r chi.Router, handler *GameAttendanceHandler) {
	r.Post("/game/{gameId}/attendance", handler.AddAttendance)
	r.Delete("/game/{gameId}/attendance/{userId}", handler.RemoveAttendance)
}

// GameAttendanceResponse represents a generic success response for game attendance operations.
type GameAttendanceResponse struct{}

// AddAttendance handles adding a user's attendance to a game.
// @Summary Add user attendance to a game
// @Description Adds the currently authenticated user as an attendee to a specific game.
// @Tags game_attendance
// @Accept json
// @Produce json
// @Param gameId path string true "ID of the game"
// @Security BearerAuth
// @Success 201 {object} GameAttendanceResponse "User attendance added successfully"
// @Failure 401 {object} contract.ErrorResponse "Unauthorized or invalid user context"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /game/{gameId}/attendance [post]
func (h *GameAttendanceHandler) AddAttendance(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")
	userID, ok := r.Context().Value(auth.ClaimsContextKey).(string)
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "invalid user context"})
		return
	}

	if err := h.service.Add(gameId, userID); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// RemoveAttendance handles removing a user's attendance from a game.
// @Summary Remove user attendance from a game
// @Description Removes a specified user's attendance from a specific game.
// @Tags game_attendance
// @Accept json
// @Produce json
// @Param gameId path string true "ID of the game"
// @Param userId path string true "ID of the user to remove attendance for"
// @Security BearerAuth
// @Success 200 {object} GameAttendanceResponse "User attendance removed successfully"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /game/{gameId}/attendance/{userId} [delete]
func (h *GameAttendanceHandler) RemoveAttendance(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")
	userId := chi.URLParam(r, "userId")

	if err := h.service.Remove(gameId, userId); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
