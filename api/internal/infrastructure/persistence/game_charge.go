package persistence

type GameCharge struct {
	Base

	GameID       string  `gorm:"type:uuid;not null;index" json:"gameId"`
	UserID       *string `gorm:"type:uuid;index" json:"userId"`
	ExternalName string  `gorm:"type:varchar(255)" json:"externalName"`
	AmountCents  int     `gorm:"not null;check:amount_cents > 0" json:"amountCents"`

	Allocations []PaymentAllocation `gorm:"foreignKey:GameChargeID" json:"allocations"`
	Game        Game                `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE" json:"game"`
}
