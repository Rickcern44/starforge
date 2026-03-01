package mappers

import (
	"testing"
	"time"

	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
)

func TestUserMapper(t *testing.T) {
	now := time.Now()
	domainUser := &models.User{
		ID:           "user-123",
		Email:        "test@example.com",
		Name:         "John Doe",
		PasswordHash: "hashed-pw",
		Roles:        []string{"admin", "user"},
		CreatedAt:    now,
	}

	// 1. To Persistence
	pUser := UserToPersistence(domainUser)
	if pUser.ID != domainUser.ID || pUser.Email != domainUser.Email || len(pUser.Roles) != 2 {
		t.Errorf("ToPersistence failed: got %+v", pUser)
	}

	// 2. To Domain
	dUser := UserToDomain(*pUser)
	if dUser.ID != domainUser.ID || dUser.Name != domainUser.Name || len(dUser.Roles) != 2 {
		t.Errorf("ToDomain failed: got %+v", dUser)
	}
}

func TestGameMapper(t *testing.T) {
	now := time.Now()
	pGame := persistence.Game{
		Base:        persistence.Base{ID: "game-1"},
		LeagueID:    "league-1",
		StartTime:   now,
		Location:    "Gym",
		CostInCents: 1000,
		IsCanceled:  false,
		Attendance: []persistence.GameAttendance{
			{UserID: "u1", CheckedIn: true, Status: 0},
		},
	}

	// 1. To Domain
	dGame := GameToDomain(pGame)
	if dGame.ID != pGame.ID || len(dGame.Attendance) != 1 || dGame.Attendance[0].UserID != "u1" {
		t.Errorf("GameToDomain failed: got %+v", dGame)
	}

	// 2. To Persistence
	pGame2 := GameToPersistence(dGame)
	if pGame2.ID != dGame.ID || pGame2.Location != dGame.Location {
		t.Errorf("GameToPersistence failed: got %+v", pGame2)
	}
}

func TestLeagueMapper(t *testing.T) {
	pLeague := persistence.League{
		Base:     persistence.Base{ID: "league-1"},
		Name:     "Monday Night",
		IsActive: true,
		Members: []persistence.LeagueMember{
			{UserID: "u1", LeagueID: "league-1", Role: "player"},
		},
		Games: []persistence.Game{
			{Base: persistence.Base{ID: "g1"}, Location: "Field"},
		},
	}

	// 1. To Domain
	dLeague := LeagueToDomain(pLeague)
	if dLeague.ID != pLeague.ID || len(dLeague.Members) != 1 || len(dLeague.Games) != 1 {
		t.Errorf("LeagueToDomain failed: got %+v", dLeague)
	}

	// 2. From Domain
	pLeague2 := LeagueFromDomain(dLeague)
	if pLeague2.ID != dLeague.ID || pLeague2.Name != dLeague.Name {
		t.Errorf("LeagueFromDomain failed: got %+v", pLeague2)
	}
}
