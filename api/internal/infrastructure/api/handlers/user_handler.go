package handlers

import (
	"encoding/json"
	"log/slog"
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

func RegisterAdminUserRoutes(r chi.Router, handler *UserHandler) {
	r.Patch("/users/{userId}/roles", handler.UpdateUserRoles)
}

// GetCurrentUser godoc
// @Summary Get the current user
// @Description Details of the currently logged-in user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User "Current user details"
// @Failure 401 {object} contract.ErrorResponse "Unauthorized"
// @Failure 404 {object} contract.ErrorResponse "User not found"
// @Router /api/users/me [get]
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.Claims)
	if !ok || claims == nil {
		slog.Error("Get current user failed: authentication required")
		utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "authentication required"})
		return
	}

	user, err := h.service.GetUserByID(claims.UserId)
	if err != nil {
		slog.Error("Get user by ID service failed", "userId", claims.UserId, "error", err)
		utils.WriteJSON(w, http.StatusNotFound, contract.ErrorResponse{Error: "user not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

// UpdateUserRolesRequest represents the request body for updating user roles.
type UpdateUserRolesRequest struct {
	Roles []string `json:"roles"`
}

// UpdateUserRoles godoc
// @Summary Update roles assigned to a user
// @Description Allows users with the admin role to update the roles of a selected user
// @Tags auth
// @Accept json
// @Produce json
// @Param userId path string true "ID of the user to update"
// @Param request body UpdateUserRolesRequest true "New roles for the user"
// @Security BearerAuth
// @Success 200 {object} object "User roles updated successfully"
// @Failure 400 {object} contract.ErrorResponse "Invalid request body"
// @Failure 404 {object} contract.ErrorResponse "User not found"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /api/users/{userId}/roles [patch]
func (h *UserHandler) UpdateUserRoles(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	var req UpdateUserRolesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode update user roles request", "userId", userID, "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: "Invalid request body"})
		return
	}

	if err := h.service.UpdateUserRoles(userID, req.Roles); err != nil {
		slog.Error("Update user roles service failed", "userId", userID, "roles", req.Roles, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: "Failed to update user roles"})
		return
	}

	w.WriteHeader(http.StatusOK)
}
