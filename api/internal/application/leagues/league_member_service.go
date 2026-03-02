package leagues

import (
	"github.com/bouncy/bouncy-api/internal/domain/interfaces"
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
	leagueMember := models.CreateLeagueMember(leagueID, userId, role)

	return s.repo.Add(leagueMember)
}

func (s *LeagueMemberService) UpdateRole(leagueId, userId string, role models.Role) error {
	//TODO: Add an admin check before allowing the user to add the member to the league
	return s.repo.UpdateRole(leagueId, userId, role)
}

func (s *LeagueMemberService) RemoveMember(leagueID, userId string) error {
	//TODO: Add an admin check before allowing the user to add the member to the league
	return s.repo.Remove(leagueID, userId)
}
