package persistence

import "time"

type Payment struct {
	Base

	UserID      string `gorm:"type:uuid;not null;index"`
	AmountCents int    `gorm:"not null;check:amount_cents > 0"`
	Method      string `gorm:"type:varchar(50);not null"`
	ReceivedAt  time.Time
	RecordedBy  string  `gorm:"type:uuid;not null"`
	Reference   *string `gorm:"type:varchar(255)"`

	Allocations []PaymentAllocation `gorm:"foreignKey:PaymentID"`
}
