package interfaces

import (
	"time"

	"github.com/bouncy/bouncy-api/internal/domain/models"
)

type LeagueRepository interface {
	GetById(league string) (*models.League, error)
	Save(league *models.League) error
	Delete(league string) error
}

type GameRepository interface {
	ListGamesByLeague(leagueId string) ([]*models.Game, error)
	Create(game *models.Game) (*models.Game, error)
	GetById(gameId string) (*models.Game, error)
	Update(gameId string, game *models.Game) (*models.Game, error)
	Cancel(gameId string) error
}

type PaymentsRepository interface {
	ListByLeague(leagueID string) ([]models.Payment, error)
	ListPaymentsByUser(userID string) ([]models.Payment, error)
	ListPaymentsByExternalName(name string) ([]models.Payment, error)
	Add(payment *models.Payment) error
	AddAllocation(paymentID string, allocation *models.PaymentAllocation) error
	CreateCharge(charge *models.GameCharge) error
	ListChargesByUser(userID string) ([]models.GameCharge, error)
	ListUnpaidChargesByUser(userID string) ([]models.GameCharge, error)
	ListChargesByExternalName(name string) ([]models.GameCharge, error)
	ListUnpaidChargesByExternalName(name string) ([]models.GameCharge, error)
	ClaimUnclaimedRecords(userID string, externalName string) error
	GetFinancialSummary(leagueID string) (*models.LeagueFinancialSummary, error)
}

type GameAttendanceRepository interface {
	Add(attendance *models.GameAttendance, gameID string) error
	Remove(gameID, userID string) error
	FindByGameAndUser(gameID string, userID string) (*models.GameAttendance, error)
	Update(attendance *models.GameAttendance, gameID string) error
}

type LeagueMemberRepository interface {
	ListByLeague(leagueID string) ([]models.LeagueMember, error)
	Add(member *models.LeagueMember) error
	UpdateRole(leagueID, userID string, role models.Role) error
	Remove(leagueID, userID string) error
	IsAdmin(leagueID, userID string) (bool, error)
}

type UserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	FindByName(name string) (*models.User, error)
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	UpdateUserRoles(userID string, roles []string) error

	// Invitation methods
	CreateInvitation(invitation *models.Invitation) error
	GetInvitationByToken(token string) (*models.Invitation, error)
	GetInvitationsByLeague(leagueID string) ([]models.Invitation, error)
	MarkInvitationAsUsed(token string, usedAt time.Time) error
}

type FeatureFlagRepository interface {
	GetAll() ([]models.FeatureFlag, error)
	GetByKey(key string) (*models.FeatureFlag, error)
	Update(key string, enabled bool) error
	Create(flag *models.FeatureFlag) error
	Delete(key string) error
}
