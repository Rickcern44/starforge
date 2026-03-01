package mappers

import (
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
)

func GameToDomain(m persistence.Game) *models.Game {
	attendance := make([]models.GameAttendance, len(m.Attendance))
	for i, a := range m.Attendance {
		attendance[i] = models.GameAttendance{
			UserID:         a.UserID,
			UserName:       a.User.Name, // Populate from preloaded User
			CheckedIn:      a.CheckedIn,
			Status:         models.AttendanceStatus(a.Status),
			CheckInComment: a.CheckInComments,
			CreatedAt:      a.CreatedAt,
			UpdatedAt:      a.UpdatedAt,
		}
	}

	charges := make([]models.GameCharge, len(m.Charges))
	for i, c := range m.Charges {
		allocations := make([]models.PaymentAllocation, len(c.Allocations))
		for j, a := range c.Allocations {
			allocations[j] = models.PaymentAllocation{
				PaymentID:     a.PaymentID,
				GameChargeID:  a.GameChargeID,
				AmountInCents: a.AmountInCents,
			}
		}

		charges[i] = models.GameCharge{
			ID:           c.ID,
			GameID:       c.GameID,
			UserID:       c.UserID,
			ExternalName: c.ExternalName,
			AmountCents:  c.AmountCents,
			CreatedAt:    c.CreatedAt,
			Allocations:  allocations,
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
		Charges:     charges,
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
