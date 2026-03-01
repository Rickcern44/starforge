package mappers

import (
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
)

func LeagueToDomain(m persistence.League) models.League {
	members := make([]models.LeagueMember, len(m.Members))
	for i, member := range m.Members {
		members[i] = models.LeagueMember{
			PlayerID:   member.UserID,
			PlayerName: member.User.Name,
			LeagueID:   member.LeagueID,
			Role:       member.Role,
			JoinedAt:   member.JoinedAt,
		}
	}

	games := make([]*models.Game, len(m.Games))
	for i, game := range m.Games {
		games[i] = GameToDomain(game)
	}

	return models.League{
		ID:       m.ID,
		Name:     m.Name,
		IsActive: m.IsActive,
		Members:  members,
		Games:    games,
	}
}

func LeagueFromDomain(d models.League) persistence.League {
	return persistence.League{
		Base:     persistence.Base{ID: d.ID},
		Name:     d.Name,
		IsActive: d.IsActive,
	}
}
