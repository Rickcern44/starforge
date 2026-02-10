package persistence

type PaymentAllocation struct {
	PaymentID     string `gorm:"type:uuid;primaryKey"`
	GameChargeID  string `gorm:"type:uuid;primaryKey"`
	AmountInCents int    `gorm:"not null;check:amount_in_cents > 0"`

	Payment    Payment    `gorm:"foreignKey:PaymentID;constraint:OnDelete:CASCADE"`
	GameCharge GameCharge `gorm:"foreignKey:GameChargeID;constraint:OnDelete:CASCADE"`
}
