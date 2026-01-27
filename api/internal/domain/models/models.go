package models

import (
	"errors"
	"time"
)

type League struct {
	ID       string
	Name     string
	IsActive bool

	Members []LeagueMember
	Games   []Game
}

type LeagueMember struct {
	UserID   string
	Role     string
	JoinedAt time.Time
}

type Game struct {
	ID        string
	LeagueID  string
	StartTime time.Time
	Location  string

	CostInCents int
	IsCanceled  bool

	Attendance []GameAttendance
	Payments   []GamePayment
}

type GameAttendance struct {
	UserID    string
	CheckedIn bool
	CreatedAt time.Time
}

type GamePayment struct {
	ID          string
	UserID      string
	AmountCents int
	Method      PaymentMethod
	Status      PaymentStatus
	PaidAt      *time.Time
	ConfirmedBy *string
}

func (l *League) AddMember(userId, role string) error {
	for _, m := range l.Members {
		if m.UserID == userId {
			return errors.New("user already in league")
		}
	}

	l.Members = append(l.Members, LeagueMember{
		UserID: userId,
		Role:   role,
	})

	return nil
}
