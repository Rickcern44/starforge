package persistence

type User struct {
	Base
	Name         string `gorm:"type:varchar(100);"`
	Email        string `gorm:"type:varchar(100);uniqueIndex"`
	PasswordHash string `gorm:"type:varchar(255);"`
}
