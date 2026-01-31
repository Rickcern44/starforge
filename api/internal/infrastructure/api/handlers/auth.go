package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application"
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
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// LoginHandler  godoc
// @Summary      Login
// @Description  Authenticate user and return JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body handlers.LoginRequest true "Login request"
// @Success      200 {object} handlers.LoginResponse
// @Failure      401 {object} ErrorResponse
// @Router       /api/auth/login [post]
func (h AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	token, err := h.service.Login(request.Email, request.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{Token: token, TokenType: "Bearer"})
}
