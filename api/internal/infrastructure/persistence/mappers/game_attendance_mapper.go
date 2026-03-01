package mappers

import (
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
)

func GameAttendanceToDomain(p persistence.GameAttendance) *models.GameAttendance {
	userName := p.UserID
	if p.User.Name != "" {
		userName = p.User.Name
	}

	return &models.GameAttendance{
		UserID:         p.UserID,
		UserName:       userName,
		CheckedIn:      p.CheckedIn,
		Status:         models.AttendanceStatus(p.Status),
		CheckInComment: p.CheckInComments,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}

func GameAttendanceToDto(m *models.GameAttendance, gameId string) *persistence.GameAttendance {
	return &persistence.GameAttendance{
		GameID:          gameId,
		UserID:          m.UserID,
		CheckedIn:       m.CheckedIn,
		Status:          m.Status.Value(),
		CheckInComments: m.CheckInComment,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
