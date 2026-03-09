package features

import (
	"github.com/bouncy/bouncy-api/internal/domain/interfaces"
	"github.com/bouncy/bouncy-api/internal/domain/models"
)

type FeatureFlagService struct {
	repo interfaces.FeatureFlagRepository
}

func NewFeatureFlagService(repo interfaces.FeatureFlagRepository) *FeatureFlagService {
	return &FeatureFlagService{repo: repo}
}

func (s *FeatureFlagService) GetAll() ([]models.FeatureFlag, error) {
	return s.repo.GetAll()
}

func (s *FeatureFlagService) IsEnabled(key string) bool {
	flag, err := s.repo.GetByKey(key)
	if err != nil {
		return false
	}
	return flag.Enabled
}

func (s *FeatureFlagService) Update(key string, enabled bool) error {
	return s.repo.Update(key, enabled)
}

func (s *FeatureFlagService) Create(flag *models.FeatureFlag) error {
	return s.repo.Create(flag)
}

func (s *FeatureFlagService) Delete(key string) error {
	return s.repo.Delete(key)
}
