package repositories

import (
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

type FeatureFlagRepository struct {
	db *gorm.DB
}

func NewFeatureFlagRepository(db *gorm.DB) *FeatureFlagRepository {
	return &FeatureFlagRepository{db: db}
}

func (r *FeatureFlagRepository) GetAll() ([]models.FeatureFlag, error) {
	var pFlags []persistence.FeatureFlag
	if err := r.db.Find(&pFlags).Error; err != nil {
		return nil, err
	}

	flags := make([]models.FeatureFlag, len(pFlags))
	for i, p := range pFlags {
		flags[i] = models.FeatureFlag{
			Key:         p.Key,
			Name:        p.Name,
			Description: p.Description,
			Enabled:     p.Enabled,
		}
	}
	return flags, nil
}

func (r *FeatureFlagRepository) GetByKey(key string) (*models.FeatureFlag, error) {
	var p persistence.FeatureFlag
	if err := r.db.Where("key = ?", key).First(&p).Error; err != nil {
		return nil, err
	}

	return &models.FeatureFlag{
		Key:         p.Key,
		Name:        p.Name,
		Description: p.Description,
		Enabled:     p.Enabled,
	}, nil
}

func (r *FeatureFlagRepository) Update(key string, enabled bool) error {
	return r.db.Model(&persistence.FeatureFlag{}).Where("key = ?", key).Update("enabled", enabled).Error
}

func (r *FeatureFlagRepository) Create(flag *models.FeatureFlag) error {
	p := persistence.FeatureFlag{
		Key:         flag.Key,
		Name:        flag.Name,
		Description: flag.Description,
		Enabled:     flag.Enabled,
	}
	return r.db.Create(&p).Error
}

func (r *FeatureFlagRepository) Delete(key string) error {
	return r.db.Where("key = ?", key).Delete(&persistence.FeatureFlag{}).Error
}
