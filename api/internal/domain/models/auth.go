package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Roles        []string  `json:"roles"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Invitation struct {
	Token       string     `json:"token"`
	Email       string     `json:"email"`
	LeagueID    string     `json:"leagueId"`
	InvitedBy   string     `json:"invitedBy"`
	SpecialRole string     `json:"specialRole,omitempty"` // Role to grant upon registration (e.g., league_creator)
	ExpiresAt   time.Time  `json:"expiresAt"`
	UsedAt      *time.Time `json:"usedAt"`
}

func (i Invitation) IsValid() bool {
	return i.UsedAt == nil && i.ExpiresAt.After(time.Now())
}
