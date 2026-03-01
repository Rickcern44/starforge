package models

import "time"

type Payment struct {
	ID            string
	UserID        *string // Nullable for unclaimed payments
	ExternalName  string  // Name from spreadsheet if UserID is null
	LeagueID      string
	AmountInCents int
	Method        PaymentMethod
	ReceivedAt    time.Time
	RecordedBy    string  // Admin ID
	Reference     *string // Venmo note / spreadsheet row

	Allocations []PaymentAllocation
}

func (p Payment) UnallocatedAmount() int {
	allocated := 0
	for _, a := range p.Allocations {
		allocated += a.AmountInCents
	}
	return p.AmountInCents - allocated
}

func CreatePayment(userID *string, externalName string, leagueID string, amount int, method PaymentMethod, reference *string) *Payment {
	return &Payment{
		UserID:        userID,
		ExternalName:  externalName,
		LeagueID:      leagueID,
		AmountInCents: amount,
		Method:        method,
		ReceivedAt:    time.Now(),
		Reference:     reference,
		Allocations:   make([]PaymentAllocation, 0),
	}
}
