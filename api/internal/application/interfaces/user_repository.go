package interfaces

import "github.com/bouncy/bouncy-api/internal/domain/models"

type UserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	FindByName(name string) (*models.User, error)
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	UpdateUserRoles(userID string, roles []string) error
}
