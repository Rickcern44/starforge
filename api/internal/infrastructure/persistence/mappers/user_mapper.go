package mappers

import (
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
)

func UserToDomain(p persistence.User) *models.Player {
	return &models.Player{
		ID:           p.ID,
		Email:        p.Email,
		Name:         p.Name,
		PasswordHash: p.PasswordHash,
		CreatedAt:    p.CreatedAt,
	}
}

func UserToPersistence(d *models.Player) *persistence.User {
	return &persistence.User{
		Base: persistence.Base{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
		},
		Email:        d.Email,
		Name:         d.Name,
		PasswordHash: d.PasswordHash,
	}
}
