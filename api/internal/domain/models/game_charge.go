package models

import "time"

type GameCharge struct {
	ID          string
	GameID      string
	UserID      string
	AmountCents int
	CreatedAt   time.Time
	Allocations []PaymentAllocation
}

func (c GameCharge) IsPaid() bool {
	paid := 0
	for _, a := range c.Allocations {
		paid += a.AmountInCents
	}
	return paid >= c.AmountCents
}
