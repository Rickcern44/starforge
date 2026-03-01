package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/contract"
	"github.com/bouncy/bouncy-api/internal/infrastructure/auth"
	"github.com/bouncy/bouncy-api/internal/infrastructure/utils"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	service *application.AuthService
}

func NewAuthHandler(service *application.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func RegisterAuthRoutes(r chi.Router, handler *AuthHandler) {
	r.Post("/auth/login", handler.LoginHandler)
	r.Post("/auth/register", handler.RegistrationHandler)
}

func RegisterAdminAuthRoutes(r chi.Router, handler *AuthHandler) {
	r.Post("/admin/invite", handler.InviteHandler)
	r.Get("/admin/league/{leagueId}/invitations", handler.ListLeagueInvitationsHandler)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.Error("Failed to decode login request", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: "Invalid request body"})
		return
	}

	token, err := h.service.Login(request.Email, request.Password)
	if err != nil {
		slog.Error("Login failed", "email", request.Email, "error", err)
		if errors.Is(err, application.ErrUserNotFound) || errors.Is(err, application.ErrInvalidCredentials) {
			utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "Unable to complete login", Message: "Invalid email or password"})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: "An internal error occurred"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, LoginResponse{Token: token, TokenType: "Bearer"})
}

type registrationRequest struct {
	Token    string `json:"token"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registrationResponse struct {
	Message string `json:"message"`
}

func (h *AuthHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var request registrationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.Error("Failed to decode registration request", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: "Invalid request body"})
		return
	}

	// We now require an invitation token for registration
	if request.Token == "" {
		slog.Error("Registration failed: missing invitation token", "email", request.Email)
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: "Invitation token is required"})
		return
	}

	err := h.service.RegisterWithInvitation(request.Token, request.Name, request.Email, request.Password)
	if err != nil {
		slog.Error("Registration failed", "email", request.Email, "error", err)
		if errors.Is(err, application.ErrUserAlreadyExists) {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		if errors.Is(err, application.ErrInvalidInvitation) {
			utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: "Invalid or expired invitation"})
			return
		}

		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, registrationResponse{Message: "User registered successfully"})
}

type InviteRequest struct {
	Email    string `json:"email"`
	LeagueID string `json:"leagueId"`
}

func (h *AuthHandler) InviteHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.Claims)
	if !ok {
		slog.Error("Invite failed: invalid user context")
		utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "invalid user context"})
		return
	}

	var request InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.Error("Failed to decode invite request", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: "Invalid request body"})
		return
	}

	if err := h.service.InviteUser(request.Email, request.LeagueID, claims.UserId); err != nil {
		slog.Error("Invite user failed", "email", request.Email, "leagueId", request.LeagueID, "invitedBy", claims.UserId, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Invitation sent successfully"})
}

func (h *AuthHandler) ListLeagueInvitationsHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.Claims)
	if !ok {
		slog.Error("List invitations failed: invalid user context")
		utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "invalid user context"})
		return
	}

	leagueId := chi.URLParam(r, "leagueId")

	invites, err := h.service.GetLeagueInvitations(leagueId, claims.UserId)
	if err != nil {
		slog.Error("List league invitations failed", "leagueId", leagueId, "userId", claims.UserId, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, invites)
}
