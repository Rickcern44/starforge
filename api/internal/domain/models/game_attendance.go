package models

import "time"

type GameAttendance struct {
	UserID         string           `json:"userId"`
	UserName       string           `json:"userName"` // Added for display
	CheckedIn      bool             `json:"checkedIn"`
	Status         AttendanceStatus `json:"status"`
	CheckInComment string           `json:"checkInComment"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
}
