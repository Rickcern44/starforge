package models

import "time"

type LeagueMember struct {
	UserID   string
	Role     Role
	JoinedAt time.Time
}
