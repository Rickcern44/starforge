package models

type PaymentAllocation struct {
	PaymentID     string `json:"paymentId"`
	GameChargeID  string `json:"gameChargeId"`
	AmountInCents int    `json:"amountInCents"`
}
