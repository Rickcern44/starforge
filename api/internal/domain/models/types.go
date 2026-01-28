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

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleUser      Role = "user"
	RoleTreasurer Role = "treasurer"
	RoleGuest     Role = "guest"
)

func (r Role) IsAdmin() bool {
	return r == RoleAdmin
}
