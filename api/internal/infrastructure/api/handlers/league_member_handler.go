package handlers

import (
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/gin-gonic/gin"
)

type LeagueMemberHandler struct {
	service *leagues.LeagueMemberService
}

func NewLeagueMemberHandler(service *leagues.LeagueMemberService) *LeagueMemberHandler {
	return &LeagueMemberHandler{service: service}
}

func RegisterLeagueMemberHandlers(rg *gin.RouterGroup, handler *LeagueMemberHandler) {
	group := rg.Group("/league")

	group.GET("/:leagueId/members", handler.ListMembers)
	group.POST("/:leagueId/members", handler.AddMember)
	group.PATCH("/:leagueId/members/:memberId", handler.UpdateRole)
	group.DELETE("/:leagueId/members/:memberId", handler.RemoveMember)
}

func (h *LeagueMemberHandler) ListMembers(c *gin.Context) {}

type AddMemberRequest struct {
	userId       string      `json:"userId"`
	addingUserId string      `json:"addingUserId"`
	role         models.Role `json:"role"`
}

func (h *LeagueMemberHandler) AddMember(c *gin.Context) {
	var request AddMemberRequest
	leagueId := c.Param("leagueId")

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddMember(leagueId, request.addingUserId, request.userId, request.role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

func (h *LeagueMemberHandler) RemoveMember(c *gin.Context) {}

func (h *LeagueMemberHandler) UpdateRole(c *gin.Context) {}
