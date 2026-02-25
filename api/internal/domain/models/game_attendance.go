package models

import "time"

type GameAttendance struct {
	UserID         string           `json:"userId"`
	CheckedIn      bool             `json:"checkedIn"`
	Status         AttendanceStatus `json:"status"`
	CheckInComment string           `json:"checkInComment"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
}
