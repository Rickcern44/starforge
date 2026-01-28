package handlers

import (
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *application.AuthService
}

func NewAuthHandler(service *application.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func RegisterAuthRoutes(rg *gin.RouterGroup, handler *AuthHandler) {
	authGroup := rg.Group("/auth")

	authGroup.POST("/login")

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
func (h AuthHandler) LoginHandler(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(request.Email, request.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "token_type": "Bearer"})
}
