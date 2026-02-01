package interfaces

import "github.com/bouncy/bouncy-api/internal/domain/models"

type UserRepository interface {
	GetUserByEmail(email string) (*models.Player, error)
	CreateUser(user *models.Player) error
}
