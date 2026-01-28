package mappers

import (
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
)

func LeagueMemberToDomain(m persistence.LeagueMember) models.LeagueMember {
	return models.LeagueMember{
		UserID:   m.UserID,
		Role:     m.Role,
		JoinedAt: m.JoinedAt,
	}
}

func LeagueMemberToDto(m models.LeagueMember) persistence.LeagueMember {
	return persistence.LeagueMember{
		UserID:   m.UserID,
		Role:     m.Role,
		JoinedAt: m.JoinedAt,
	}
}
