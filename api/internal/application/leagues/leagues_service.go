package leagues

import (
	"context"
	"errors"

	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence/repositories"
	"github.com/google/uuid"
)

type LeagueService struct {
	leagueRepo repositories.LeagueRepository
}

func NewLeagueService(leagueRepo *repositories.LeagueRepository) *LeagueService {
	return &LeagueService{
		leagueRepo: *leagueRepo,
	}
}

func (ls *LeagueService) GetLeague(leagueId string) (*models.League, error) {
	return ls.leagueRepo.GetById(leagueId)
}

func (ls *LeagueService) CreateLeague(ctx context.Context, name string) (*models.League, error) {
	id, _ := uuid.NewV7()

	league := &models.League{
		ID:       id.String(),
		Name:     name,
		IsActive: true,
	}

	if err := ls.leagueRepo.Save(league); err != nil {
		return nil, err
	}

	return league, nil
}

func (ls *LeagueService) AddMember(leagueId, userId string, role models.Role) error {
	league, err := ls.leagueRepo.GetById(leagueId)

	if err != nil {
		return err
	}

	if league == nil {
		return errors.New("league not found")
	}

	if err := league.AddMember(userId, role); err != nil {
		return err
	}

	return ls.leagueRepo.Save(league)
}

func (ls *LeagueService) AddGame(leagueId, userId, role string) error {
	return nil
}

func (ls *LeagueService) Delete(leagueId string) error {
	return ls.leagueRepo.Delete(leagueId)
}

func (ls *LeagueService) GetLeaguesForUser(userId string) ([]*models.League, error) {
	return ls.leagueRepo.FindLeaguesByUserID(userId)
}
