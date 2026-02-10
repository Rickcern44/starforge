package models

import "time"

type LeagueMember struct {
	LeagueID string
	PlayerID string
	Role     Role
	JoinedAt time.Time
}

func CreateLeagueMember(leagueID, playerID string, role Role) *LeagueMember {
	return &LeagueMember{
		LeagueID: leagueID,
		PlayerID: playerID,
		Role:     role,
		JoinedAt: time.Now(),
	}
}
