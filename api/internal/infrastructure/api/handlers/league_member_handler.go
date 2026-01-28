package handlers

import (
	"github.com/bouncy/bouncy-api/internal/application/leagues"
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

func (h *LeagueMemberHandler) AddMember(c *gin.Context) {}

func (h *LeagueMemberHandler) RemoveMember(c *gin.Context) {}

func (h *LeagueMemberHandler) UpdateRole(c *gin.Context) {}
