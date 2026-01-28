package persistence

type GameCharge struct {
	Base

	GameID      string `gorm:"type:uuid;not null;index"`
	UserID      string `gorm:"type:uuid;not null;index"`
	AmountCents int    `gorm:"not null;check:amount_cents > 0"`

	Allocations []PaymentAllocation `gorm:"foreignKey:GameChargeID"`
}
