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

type AttendanceStatus int

const (
	Yes AttendanceStatus = iota
	No
	Tentative
)

func (r AttendanceStatus) String() string {
	return [...]string{"Yes", "No", "Tentative"}[r]
}

func (r AttendanceStatus) Value() int {
	return int(r)
}
