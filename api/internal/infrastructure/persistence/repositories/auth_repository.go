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

func (r *AuthRepository) GetUserByEmail(email string) (*models.Player, error) {
	var user persistence.User // Use persistence model for query
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return mappers.UserToDomain(user), nil // Convert to domain model before returning
}

func (r *AuthRepository) CreateUser(user *models.Player) error {
	persistenceUser := mappers.UserToPersistence(user) // Convert domain model to persistence model
	return r.db.Create(persistenceUser).Error
}
