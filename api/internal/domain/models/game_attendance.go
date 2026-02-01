package models

import "time"

type GameAttendance struct {
	UserID         string
	CheckedIn      bool
	CheckInComment string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
