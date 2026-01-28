package repositories

import (
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence/mappers"
	"gorm.io/gorm"
)

type LeagueMemberRepository struct {
	db *gorm.DB
}

func NewLeagueMemberRepository(db *gorm.DB) *LeagueMemberRepository {
	return &LeagueMemberRepository{db: db}
}

func (l LeagueMemberRepository) ListByLeague(leagueID string) ([]models.LeagueMember, error) {
	var rows []persistence.LeagueMember
	if err := l.db.Where("league_id = ?", leagueID).Find(&rows).Error; err != nil {
		return nil, err
	}

	var members []models.LeagueMember
	for _, row := range rows {
		members = append(members, mappers.LeagueMemberToDomain(row))
	}
	return members, nil
}

func (l LeagueMemberRepository) Add(member models.LeagueMember) error {
	leagueMemberDto := mappers.LeagueMemberToDto(member)
	return l.db.Create(&leagueMemberDto).Error
}

func (l LeagueMemberRepository) UpdateRole(leagueID, userID string, role models.Role) error {
	//TODO implement me
	panic("implement me")
}

func (l LeagueMemberRepository) Remove(leagueID, userID string) error {
	return l.db.
		Where("league_id = ? AND user_id = ?", leagueID, userID).
		Delete(&persistence.LeagueMember{}).
		Error
}

func (l LeagueMemberRepository) IsAdmin(leagueID, userID string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
