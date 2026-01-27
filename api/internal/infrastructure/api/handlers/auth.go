package handlers

import "github.com/gin-gonic/gin"

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (a AuthHandler) LoginHandler(c *gin.Context) {

}
