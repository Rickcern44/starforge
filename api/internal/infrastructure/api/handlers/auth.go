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
	r.Post("/auth/login", handler.LoginHandler)
	r.Post("/auth/register", handler.RegistrationHandler)
}

// LoginRequest represents the request body for user login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response body for a successful user login.
type LoginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}

// LoginHandler handles user login.
// @Summary User Login
// @Description Authenticates a user and returns a JWT token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse "Authentication successful"
// @Failure 400 {object} contract.ErrorResponse "Invalid request body"
// @Failure 401 {object} contract.ErrorResponse "Invalid email or password"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /v1/auth/login [post]
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: "Invalid request body"})
		return
	}

	token, err := h.service.Login(request.Email, request.Password)
	if err != nil {
		if errors.Is(err, application.ErrUserNotFound) || errors.Is(err, application.ErrInvalidCredentials) {
			utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "Unable to complete login", Message: "Invalid email or password"})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: "An internal error occurred"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, LoginResponse{Token: token, TokenType: "Bearer"})
}

// registrationRequest represents the request body for user registration.
type registrationRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// registrationResponse represents the response body for a successful user registration.
type registrationResponse struct {
	Message string `json:"message"`
}

// RegistrationHandler handles user registration.
// @Summary User Registration
// @Description Registers a new user with the provided details.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registrationRequest true "User registration details"
// @Success 201 {object} registrationResponse "User registered successfully"
// @Failure 400 {object} contract.ErrorResponse "Invalid request body or validation error"
// @Failure 409 {object} contract.ErrorResponse "User with given email already exists"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /v1/auth/register [post]
func (h *AuthHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var request registrationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: "Invalid request body"})
		return
	}

	err := h.service.Register(request.Name, request.Email, request.Password)
	if err != nil {
		if errors.Is(err, application.ErrUserAlreadyExists) {
			//utils.WriteJSON(w, http.StatusConflict, contract.ErrorResponse{Error: err.Error()})
			http.Error(w, "User already exists", http.StatusConflict)
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
