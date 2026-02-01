package repositories

import (
	"errors"

	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence/mappers"
	"gorm.io/gorm"
)

type LeagueRepository struct {
	db *gorm.DB
}

func NewLeagueRepository(db *gorm.DB) *LeagueRepository {
	return &LeagueRepository{
		db: db,
	}
}

func (lr *LeagueRepository) GetById(leagueId string) (*models.League, error) {
	var model persistence.League
	err := lr.db.Preload("Members").Preload("Games").First(&model, "id = ?", leagueId).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	league := mappers.LeagueToDomain(model)

	return &league, nil
}

func (lr *LeagueRepository) Save(league *models.League) error {
	return lr.db.Transaction(func(tx *gorm.DB) error {
		model := mappers.LeagueFromDomain(*league)

		// Upsert league root
		if err := tx.Save(&model).Error; err != nil {
			return err
		}

		// Replace members (aggregate-owned)
		if err := tx.
			Where("league_id = ?", league.ID).
			Delete(&persistence.LeagueMember{}).
			Error; err != nil {
			return err
		}

		for _, m := range league.Members {
			member := persistence.LeagueMember{
				LeagueID: league.ID,
				UserID:   m.PlayerID,
				Role:     m.Role,
				JoinedAt: m.JoinedAt,
			}

			if err := tx.Create(&member).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
func (lr *LeagueRepository) Delete(id string) error {
	return lr.db.Delete(&persistence.League{}, "id = ?", id).Error
}
