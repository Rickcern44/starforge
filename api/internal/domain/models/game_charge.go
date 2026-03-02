package models

import (
	"time"

	"github.com/google/uuid"
)

type GameCharge struct {
	ID           string              `json:"id"`
	GameID       string              `json:"gameId"`
	UserID       *string             `json:"userId"`       // Nullable for unclaimed charges
	ExternalName string              `json:"externalName"` // Name from spreadsheet if UserID is null
	AmountCents  int                 `json:"amountCents"`
	CreatedAt    time.Time           `json:"createdAt"`
	Allocations  []PaymentAllocation `json:"allocations"`
	Game         *Game               `json:"game,omitempty"`
}

func (c GameCharge) IsPaid() bool {
	paid := 0
	for _, a := range c.Allocations {
		paid += a.AmountInCents
	}
	return paid >= c.AmountCents
}

func CreateGameCharge(gameID string, userID *string, externalName string, amount int) *GameCharge {
	id, _ := uuid.NewV7()
	return &GameCharge{
		ID:           id.String(),
		GameID:       gameID,
		UserID:       userID,
		ExternalName: externalName,
		AmountCents:  amount,
		CreatedAt:    time.Now(),
		Allocations:  make([]PaymentAllocation, 0),
	}
}
