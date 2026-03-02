package repositories

import (
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

type PaymentsRepository struct {
	db *gorm.DB
}

func NewPaymentsRepository(db *gorm.DB) *PaymentsRepository {
	return &PaymentsRepository{db: db}
}

func (r *PaymentsRepository) ListByLeague(leagueID string) ([]models.Payment, error) {
	var payments []persistence.Payment
	if err := r.db.Preload("Allocations").Where("league_id = ?", leagueID).Find(&payments).Error; err != nil {
		return nil, err
	}

	var domainPayments []models.Payment
	for _, p := range payments {
		domainPayments = append(domainPayments, persistenceToDomainPayment(p))
	}
	return domainPayments, nil
}

func (r *PaymentsRepository) ListPaymentsByUser(userID string) ([]models.Payment, error) {
	var payments []persistence.Payment
	if err := r.db.Preload("Allocations").Where("user_id = ?", userID).Find(&payments).Error; err != nil {
		return nil, err
	}

	var domain []models.Payment
	for _, p := range payments {
		domain = append(domain, persistenceToDomainPayment(p))
	}
	return domain, nil
}

func (r *PaymentsRepository) ListPaymentsByExternalName(name string) ([]models.Payment, error) {
	var payments []persistence.Payment
	if err := r.db.Preload("Allocations").Where("external_name = ?", name).Find(&payments).Error; err != nil {
		return nil, err
	}

	var domain []models.Payment
	for _, p := range payments {
		domain = append(domain, persistenceToDomainPayment(p))
	}
	return domain, nil
}

func persistenceToDomainPayment(p persistence.Payment) models.Payment {
	allocations := make([]models.PaymentAllocation, len(p.Allocations))
	for i, a := range p.Allocations {
		allocations[i] = models.PaymentAllocation{
			PaymentID:     a.PaymentID,
			GameChargeID:  a.GameChargeID,
			AmountInCents: a.AmountCents,
		}
	}

	return models.Payment{
		ID:            p.ID,
		AmountCents:   p.AmountCents,
		UserID:        p.UserID,
		ExternalName:  p.ExternalName,
		LeagueID:      p.LeagueID,
		ReceivedAt:    p.ReceivedAt,
		Method:        models.PaymentMethod(p.Method),
		Reference:     p.Reference,
		Allocations:   allocations,
	}
}

func (r *PaymentsRepository) Add(payment *models.Payment) error {
	p := &persistence.Payment{
		Base:         persistence.Base{ID: payment.ID},
		AmountCents:  payment.AmountCents,
		UserID:       payment.UserID,
		ExternalName: payment.ExternalName,
		LeagueID:     payment.LeagueID,
		ReceivedAt:   payment.ReceivedAt,
		Method:       string(payment.Method),
		RecordedBy:   payment.RecordedBy,
		Reference:    payment.Reference,
	}
	return r.db.Create(p).Error
}

func (r *PaymentsRepository) AddAllocation(paymentID string, allocation *models.PaymentAllocation) error {
	a := &persistence.PaymentAllocation{
		PaymentID:    paymentID,
		GameChargeID: allocation.GameChargeID,
		AmountCents:  allocation.AmountInCents,
	}
	return r.db.Create(a).Error
}

func (r *PaymentsRepository) CreateCharge(charge *models.GameCharge) error {
	c := &persistence.GameCharge{
		Base:         persistence.Base{ID: charge.ID},
		GameID:       charge.GameID,
		UserID:       charge.UserID,
		ExternalName: charge.ExternalName,
		AmountCents:  charge.AmountCents,
	}
	return r.db.Create(c).Error
}

func (r *PaymentsRepository) ListChargesByUser(userID string) ([]models.GameCharge, error) {
	var charges []persistence.GameCharge
	if err := r.db.Preload("Allocations").Preload("Game").Where("user_id = ?", userID).Find(&charges).Error; err != nil {
		return nil, err
	}

	var domain []models.GameCharge
	for _, c := range charges {
		domain = append(domain, persistenceToDomainCharge(c))
	}
	return domain, nil
}

func (r *PaymentsRepository) ListUnpaidChargesByUser(userID string) ([]models.GameCharge, error) {
	var charges []persistence.GameCharge
	if err := r.db.Preload("Allocations").Preload("Game").Where("user_id = ?", userID).Order("created_at asc").Find(&charges).Error; err != nil {
		return nil, err
	}

	var domain []models.GameCharge
	for _, c := range charges {
		dc := persistenceToDomainCharge(c)
		if !dc.IsPaid() {
			domain = append(domain, dc)
		}
	}
	return domain, nil
}

func (r *PaymentsRepository) ListChargesByExternalName(name string) ([]models.GameCharge, error) {
	var charges []persistence.GameCharge
	if err := r.db.Preload("Allocations").Preload("Game").Where("external_name = ?", name).Find(&charges).Error; err != nil {
		return nil, err
	}

	var domain []models.GameCharge
	for _, c := range charges {
		domain = append(domain, persistenceToDomainCharge(c))
	}
	return domain, nil
}

func (r *PaymentsRepository) ListUnpaidChargesByExternalName(name string) ([]models.GameCharge, error) {
	var charges []persistence.GameCharge
	if err := r.db.Preload("Allocations").Preload("Game").Where("external_name = ?", name).Order("created_at asc").Find(&charges).Error; err != nil {
		return nil, err
	}

	var domain []models.GameCharge
	for _, c := range charges {
		dc := persistenceToDomainCharge(c)
		if !dc.IsPaid() {
			domain = append(domain, dc)
		}
	}
	return domain, nil
}

func persistenceToDomainCharge(c persistence.GameCharge) models.GameCharge {
	allocations := make([]models.PaymentAllocation, len(c.Allocations))
	for i, a := range c.Allocations {
		allocations[i] = models.PaymentAllocation{
			PaymentID:     a.PaymentID,
			GameChargeID:  a.GameChargeID,
			AmountInCents: a.AmountCents,
		}
	}

	var domainGame *models.Game
	if c.Game.ID != "" {
		// We use a simplified mapping here or call mappers.GameToDomain
		// For now let's just populate basic info to avoid circular dependency if mappers import repositories
		domainGame = &models.Game{
			ID:          c.Game.ID,
			LeagueID:    c.Game.LeagueID,
			StartTime:   c.Game.StartTime,
			Location:    c.Game.Location,
			CostInCents: c.Game.CostInCents,
			IsCanceled:  c.Game.IsCanceled,
		}
	}

	return models.GameCharge{
		ID:           c.ID,
		GameID:       c.GameID,
		UserID:       c.UserID,
		ExternalName: c.ExternalName,
		AmountCents:  c.AmountCents,
		CreatedAt:    c.CreatedAt,
		Allocations:  allocations,
		Game:         domainGame,
	}
}

func (r *PaymentsRepository) ClaimUnclaimedRecords(userID string, externalName string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Claim Charges
		if err := tx.Model(&persistence.GameCharge{}).
			Where("external_name = ? AND user_id IS NULL", externalName).
			Update("user_id", userID).Error; err != nil {
			return err
		}

		// 2. Claim Payments
		if err := tx.Model(&persistence.Payment{}).
			Where("external_name = ? AND user_id IS NULL", externalName).
			Update("user_id", userID).Error; err != nil {
			return err
		}

		return nil
	})
}
