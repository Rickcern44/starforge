package models

import "time"

type Game struct {
	ID        string
	LeagueID  string
	StartTime time.Time
	Location  string

	CostInCents int
	IsCanceled  bool

	Attendance []GameAttendance
	Charges    []GameCharge
}
