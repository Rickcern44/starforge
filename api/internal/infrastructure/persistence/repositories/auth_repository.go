package repositories

import (
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence/mappers"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	var user persistence.User // Use persistence model for query
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return mappers.UserToDomain(user), nil
}

func (r *AuthRepository) CreateUser(user *models.User) error {
	persistenceUser := mappers.UserToPersistence(user)
	return r.db.Create(persistenceUser).Error
}

func (r *AuthRepository) GetUserByID(id string) (*models.User, error) {
	var user persistence.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return mappers.UserToDomain(user), nil
}

func (r *AuthRepository) UpdateUserRoles(userID string, roles []string) error {
	var user persistence.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	user.Roles = roles
	return r.db.Save(&user).Error
}
