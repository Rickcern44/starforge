package leagues

import (
	"time"

	"github.com/bouncy/bouncy-api/internal/application/interfaces"
	"github.com/bouncy/bouncy-api/internal/domain/models"
)

type LeagueMemberService struct {
	repo interfaces.LeagueMemberRepository
}

func NewLeagueMemberService(repo interfaces.LeagueMemberRepository) *LeagueMemberService {
	return &LeagueMemberService{repo: repo}
}

func (s *LeagueMemberService) ListMembers(leagueId string) ([]models.LeagueMember, error) {
	return s.repo.ListByLeague(leagueId)
}

func (s *LeagueMemberService) AddMember(leagueID, addingUserId, userId string, role models.Role) error {
	//TODO: Add an admin check before allowing the user to add the member to the league

	return s.repo.Add(models.LeagueMember{
		UserID:   userId,
		Role:     role,
		JoinedAt: time.Time{},
	})
}

func (s *LeagueMemberService) UpdateRole(leagueId, userId string, role models.Role) error {
	//TODO: Add an admin check before allowing the user to add the member to the league
	return s.repo.UpdateRole(leagueId, userId, role)
}

func (s *LeagueMemberService) RemoveMember(leagueID, userId string) error {
	//TODO: Add an admin check before allowing the user to add the member to the league
	return s.repo.Remove(leagueID, userId)
}
