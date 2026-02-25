package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Roles        []string  `json:"roles"`
	CreatedAt    time.Time `json:"createdAt"`
}
