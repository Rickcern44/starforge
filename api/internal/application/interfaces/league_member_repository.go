package interfaces

import (
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
)

type LeagueMemberRepository interface {
	ListByLeague(leagueID string) ([]models.LeagueMember, error)
	Add(member *persistence.LeagueMember) error
	UpdateRole(leagueID, userID string, role models.Role) error
	Remove(leagueID, userID string) error
	IsAdmin(leagueID, userID string) (bool, error)
}
