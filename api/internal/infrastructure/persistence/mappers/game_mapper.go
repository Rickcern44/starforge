package mappers

import (
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
)

func GameToDomain(m persistence.Game) *models.Game {
	attendance := make([]models.GameAttendance, len(m.Attendance))
	for i, a := range m.Attendance {
		attendance[i] = models.GameAttendance{
			UserID:    a.UserID,
			CheckedIn: a.CheckedIn,
			CreatedAt: a.CreatedAt,
		}
	}

	return &models.Game{
		ID:          m.ID,
		LeagueID:    m.LeagueID,
		StartTime:   m.StartTime,
		Location:    m.Location,
		CostInCents: m.CostInCents,
		IsCanceled:  m.IsCanceled,
		Attendance:  attendance,
	}
}

func GameToPersistence(d *models.Game) *persistence.Game {
	return &persistence.Game{
		Base: persistence.Base{
			ID: d.ID,
		},
		LeagueID:    d.LeagueID,
		StartTime:   d.StartTime,
		Location:    d.Location,
		CostInCents: d.CostInCents,
		IsCanceled:  d.IsCanceled,
	}
}

func GameToDto(d *models.Game) *models.Game {
	return &models.Game{
		ID:          d.ID,
		LeagueID:    d.LeagueID,
		StartTime:   d.StartTime,
		Location:    d.Location,
		CostInCents: d.CostInCents,
		IsCanceled:  d.IsCanceled,
	}
}
