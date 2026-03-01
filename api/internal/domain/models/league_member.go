package models

import "time"

type LeagueMember struct {
	LeagueID   string    `json:"leagueId"`
	PlayerID   string    `json:"playerId"`
	PlayerName string    `json:"playerName"`
	Role       Role      `json:"role"`
	JoinedAt   time.Time `json:"joinedAt"`
}

func CreateLeagueMember(leagueID, playerID string, role Role) *LeagueMember {
	return &LeagueMember{
		LeagueID: leagueID,
		PlayerID: playerID,
		Role:     role,
		JoinedAt: time.Now(),
	}
}
