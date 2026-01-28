package models

import "time"

type GameAttendance struct {
	UserID    string
	CheckedIn bool
	CreatedAt time.Time
}
