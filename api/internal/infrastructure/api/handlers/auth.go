package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

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

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	token, err := h.service.Login(request.Email, request.Password)
	if err != nil {
		if errors.Is(err, application.ErrUserNotFound) || errors.Is(err, application.ErrInvalidCredentials) {
			WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Invalid email or password"})
			return
		}
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "An internal error occurred"})
		return
	}

	WriteJSON(w, http.StatusOK, LoginResponse{Token: token, TokenType: "Bearer"})
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
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	err := h.service.Register(request.Name, request.Email, request.Password)
	if err != nil {
		if errors.Is(err, application.ErrUserAlreadyExists) {
			WriteJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
			return
		}
		// Catch validation errors
		if err.Error() == "invalid email address" || err.Error() == "password must be at least 8 characters" {
			WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		// For any other error, return a generic internal server error
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "An internal error occurred"})
		return
	}

	WriteJSON(w, http.StatusCreated, registrationResponse{Message: "User registered successfully"})
}
