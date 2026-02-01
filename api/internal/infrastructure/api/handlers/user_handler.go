package handlers

import (
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application/users"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/contract"
	"github.com/bouncy/bouncy-api/internal/infrastructure/auth"
	"github.com/bouncy/bouncy-api/internal/infrastructure/utils"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service *users.Service
}

func NewUserHandler(service *users.Service) *UserHandler {
	return &UserHandler{service: service}
}

func RegisterUserRoutes(r chi.Router, handler *UserHandler) {
	r.Get("/users/me", handler.GetCurrentUser)
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.Claims)
	if !ok || claims == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "authentication required"})
		return
	}

	user, err := h.service.GetUserByID(claims.UserId)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, contract.ErrorResponse{Error: "user not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}
