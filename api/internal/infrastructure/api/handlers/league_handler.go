package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/gin-gonic/gin"
)

type LeagueHandler struct {
	service *application.LeagueService
}

func NewLeagueHandler(service *application.LeagueService) *LeagueHandler {
	return &LeagueHandler{service: service}
}

func RegisterLeagueRoutes(rg *gin.RouterGroup, handler *LeagueHandler) {
	leagues := rg.Group("/league")

	leagues.POST("", handler.CreateLeague)
	leagues.GET("/:id", handler.GetLeague)
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
	UserID string `json:"userId"`
	Role   string `json:"role"`
}

func (h *LeagueHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	leagueID := r.PathValue("id")

	var req addMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.AddMember(leagueID, req.UserID, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
