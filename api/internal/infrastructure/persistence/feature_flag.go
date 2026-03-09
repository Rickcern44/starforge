package persistence

type FeatureFlag struct {
	Base
	Key         string `gorm:"uniqueIndex;type:varchar(100);not null"`
	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Enabled     bool   `gorm:"default:false"`
}
