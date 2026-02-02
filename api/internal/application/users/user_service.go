package users

import (
	"github.com/bouncy/bouncy-api/internal/application/interfaces"
	"github.com/bouncy/bouncy-api/internal/domain/models"
)

type Service struct {
	userRepo interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) GetUserByID(id string) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *Service) UpdateUserRoles(userID string, roles []string) error {
	return s.userRepo.UpdateUserRoles(userID, roles)
}
