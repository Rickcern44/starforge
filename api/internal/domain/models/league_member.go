package models

import "time"

type LeagueMember struct {
	UserID   string
	Role     Role
	JoinedAt time.Time
}

func CreateLeagueMember(userId string, role Role) *LeagueMember {
	return &LeagueMember{
		UserID:   userId,
		Role:     role,
		JoinedAt: time.Now(),
	}
}
