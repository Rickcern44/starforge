package persistence

import "time"

type Payment struct {
	Base

	UserID        *string `gorm:"type:uuid;index"` // Nullable for unclaimed payments
	ExternalName  string  `gorm:"type:varchar(255)"`
	LeagueID      string  `gorm:"type:uuid;not null;index"`
	AmountInCents int     `gorm:"not null;check:amount_in_cents > 0"`
	Method        string  `gorm:"type:varchar(50);not null"`
	ReceivedAt    time.Time
	RecordedBy    string  `gorm:"type:uuid;not null"`
	Reference     *string `gorm:"type:varchar(255)"`

	Allocations []PaymentAllocation `gorm:"foreignKey:PaymentID"`
}
