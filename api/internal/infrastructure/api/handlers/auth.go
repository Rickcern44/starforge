package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/contract"
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
	r.Post("/v1/auth/login", handler.LoginHandler)
	r.Post("/v1/auth/register", handler.RegistrationHandler)
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
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: "Invalid request body"})
		return
	}

	token, err := h.service.Login(request.Email, request.Password)
	if err != nil {
		if errors.Is(err, application.ErrUserNotFound) || errors.Is(err, application.ErrInvalidCredentials) {
			utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "Invalid email or password"})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: "An internal error occurred"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, LoginResponse{Token: token, TokenType: "Bearer"})
}

type registrationRequest struct {
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
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: "Invalid request body"})
		return
	}

	err := h.service.Register(request.Name, request.Email, request.Password)
	if err != nil {
		if errors.Is(err, application.ErrUserAlreadyExists) {
			utils.WriteJSON(w, http.StatusConflict, contract.ErrorResponse{Error: err.Error()})
			return
		}
		// Catch validation errors
		if err.Error() == "invalid email address" || err.Error() == "password must be at least 8 characters" {
			utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: err.Error()})
			return
		}

		// For any other error, return a generic internal server error
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: "An internal error occurred"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, registrationResponse{Message: "User registered successfully"})
}
