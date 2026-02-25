package models

import "time"

type GameCharge struct {
	ID          string              `json:"id"`
	GameID      string              `json:"gameId"`
	UserID      string              `json:"userId"`
	AmountCents int                 `json:"amountCents"`
	CreatedAt   time.Time           `json:"createdAt"`
	Allocations []PaymentAllocation `json:"allocations"`
}

func (c GameCharge) IsPaid() bool {
	paid := 0
	for _, a := range c.Allocations {
		paid += a.AmountInCents
	}
	return paid >= c.AmountCents
}
