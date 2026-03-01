package repositories

import (
	"time"

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
	var user persistence.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return mappers.UserToDomain(user), nil
}

func (r *AuthRepository) FindByName(name string) (*models.User, error) {
	var user persistence.User
	if err := r.db.Where("name = ?", name).First(&user).Error; err != nil {
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

func (r *AuthRepository) CreateInvitation(inv *models.Invitation) error {
	p := &persistence.Invitation{
		Token:     inv.Token,
		Email:     inv.Email,
		LeagueID:  inv.LeagueID,
		InvitedBy: inv.InvitedBy,
		ExpiresAt: inv.ExpiresAt,
	}
	return r.db.Create(p).Error
}

func (r *AuthRepository) GetInvitationByToken(token string) (*models.Invitation, error) {
	var p persistence.Invitation
	if err := r.db.Where("token = ?", token).First(&p).Error; err != nil {
		return nil, err
	}
	return &models.Invitation{
		Token:     p.Token,
		Email:     p.Email,
		LeagueID:  p.LeagueID,
		InvitedBy: p.InvitedBy,
		ExpiresAt: p.ExpiresAt,
		UsedAt:    p.UsedAt,
	}, nil
}

func (r *AuthRepository) GetInvitationsByLeague(leagueID string) ([]models.Invitation, error) {
	var pInvites []persistence.Invitation
	if err := r.db.Where("league_id = ?", leagueID).Find(&pInvites).Error; err != nil {
		return nil, err
	}

	invites := make([]models.Invitation, len(pInvites))
	for i, p := range pInvites {
		invites[i] = models.Invitation{
			Token:     p.Token,
			Email:     p.Email,
			LeagueID:  p.LeagueID,
			InvitedBy: p.InvitedBy,
			ExpiresAt: p.ExpiresAt,
			UsedAt:    p.UsedAt,
		}
	}
	return invites, nil
}

func (r *AuthRepository) MarkInvitationAsUsed(token string, usedAt time.Time) error {
	return r.db.Model(&persistence.Invitation{}).Where("token = ?", token).Update("used_at", usedAt).Error
}
