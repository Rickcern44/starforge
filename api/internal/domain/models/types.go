package models

type PaymentMethod string

const (
	PaymentMethodVenmo PaymentMethod = "venmo"
	PaymentMethodCash  PaymentMethod = "cash"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusComplete PaymentStatus = "confirmed"
)
