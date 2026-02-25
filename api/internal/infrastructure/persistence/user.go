package persistence

type User struct {
	Base
	Name         string `gorm:"type:varchar(100);" json:"name"`
	Email        string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	PasswordHash string `gorm:"type:varchar(255);" json:"-"`
	Roles        Roles  `gorm:"type:jsonb" json:"roles"`
}
