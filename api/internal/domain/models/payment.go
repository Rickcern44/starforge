package models

import "time"

type Payment struct {
	ID          string
	UserID      string
	AmountCents int
	Method      PaymentMethod
	ReceivedAt  time.Time
	RecordedBy  string  // Admin ID
	Reference   *string // Venmo note / spreadsheet row

	Allocations []PaymentAllocation
}

func (p Payment) UnallocatedAmount() int {
	allocated := 0
	for _, a := range p.Allocations {
		allocated += a.AmountCents
	}
	return p.AmountCents - allocated
}
