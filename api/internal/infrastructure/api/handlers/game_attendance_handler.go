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

func RegisterGameAttendanceRoutes(r chi.Router, handler *GameAttendanceHandler) {
	r.Post("/game/{gameId}/attendance", handler.AddAttendance)
	r.Delete("/game/{gameId}/attendance/{userId}", handler.RemoveAttendance)
}

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

func (h *GameAttendanceHandler) RemoveAttendance(w http.ResponseWriter, r *http.Request) {
	gameId := chi.URLParam(r, "gameId")
	userId := chi.URLParam(r, "userId")

	if err := h.service.Remove(gameId, userId); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
