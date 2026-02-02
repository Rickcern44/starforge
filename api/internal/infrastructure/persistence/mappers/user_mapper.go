package mappers

import (
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
)

func UserToDomain(p persistence.User) *models.User {
	return &models.User{
		ID:           p.ID,
		Email:        p.Email,
		Name:         p.Name,
		PasswordHash: p.PasswordHash,
		Roles:        []string(p.Roles), // Explicitly cast
		CreatedAt:    p.CreatedAt,
	}
}

func UserToPersistence(d *models.User) *persistence.User {
	return &persistence.User{
		Base: persistence.Base{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
		},
		Email:        d.Email,
		Name:         d.Name,
		PasswordHash: d.PasswordHash,
		Roles:        persistence.Roles(d.Roles), // Explicitly cast
	}
}
