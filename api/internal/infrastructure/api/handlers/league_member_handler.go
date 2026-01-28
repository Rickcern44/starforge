package handlers

import (
	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/gin-gonic/gin"
)

type LeagueMemberHandler struct {
	service *leagues.MemberService
}

func NewLeagueMemberHandler(service *leagues.MemberService) *LeagueMemberHandler {
	return &LeagueMemberHandler{service: service}
}

func RegisterLeagueMemberHandlers(rg *gin.RouterGroup, service *leagues.MemberService) {
	group := rg.Group("/league")

	group.GET("/:id/members")
	group.POST("/:id/members")
}
