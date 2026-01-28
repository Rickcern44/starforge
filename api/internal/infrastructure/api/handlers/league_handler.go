package handlers

import (
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/gin-gonic/gin"
)

type LeagueHandler struct {
	service *leagues.LeagueService
}

func NewLeagueHandler(service *leagues.LeagueService) *LeagueHandler {
	return &LeagueHandler{service: service}
}

func RegisterLeagueRoutes(rg *gin.RouterGroup, handler *LeagueHandler) {
	leagues := rg.Group("/league")

	leagues.POST("", handler.CreateLeague)
	leagues.GET("/:id", handler.GetLeague)
	leagues.DELETE("/:id", handler.Delete)
}

func (h *LeagueHandler) GetLeague(c *gin.Context) {
	id := c.Param("id")

	game, err := h.service.GetLeague(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, game)
}

type createLeagueRequest struct {
	Name string `json:"name" binding:"required,min=3"`
}

func (h *LeagueHandler) CreateLeague(c *gin.Context) {
	var req createLeagueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	league, err := h.service.CreateLeague(
		c.Request.Context(),
		req.Name,
	)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, league)
}

type addMemberRequest struct {
	UserID string      `json:"userId"`
	Role   models.Role `json:"role"`
}

func (h *LeagueHandler) AddMember(c *gin.Context) {
	leagueId := c.Param("id")

	var req addMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AddMember(leagueId, req.UserID, req.Role)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (h *LeagueHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
