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
	Token     string
	Email     string
	LeagueID  string
	InvitedBy string
	ExpiresAt time.Time
	UsedAt    *time.Time
}

func (i Invitation) IsValid() bool {
	return i.UsedAt == nil && i.ExpiresAt.After(time.Now())
}
