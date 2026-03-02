package models

import "time"

type Payment struct {
	ID            string              `json:"id"`
	UserID        *string             `json:"userId"` // Nullable for unclaimed payments
	ExternalName  string              `json:"externalName"`
	LeagueID      string              `json:"leagueId"`
	AmountCents   int                 `json:"amountCents"`
	Method        PaymentMethod       `json:"method"`
	ReceivedAt    time.Time           `json:"receivedAt"`
	RecordedBy    string              `json:"recordedBy"`
	Reference     *string             `json:"reference"`
	Allocations   []PaymentAllocation `json:"allocations"`
}

func (p Payment) UnallocatedAmount() int {
	allocated := 0
	for _, a := range p.Allocations {
		allocated += a.AmountInCents
	}
	return p.AmountCents - allocated
}

func CreatePayment(userID *string, externalName string, leagueID string, amount int, method PaymentMethod, reference *string) *Payment {
	return &Payment{
		UserID:        userID,
		ExternalName:  externalName,
		LeagueID:      leagueID,
		AmountCents:   amount,
		Method:        method,
		ReceivedAt:    time.Now(),
		Reference:     reference,
		Allocations:   make([]PaymentAllocation, 0),
	}
}
