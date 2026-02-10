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
	if err := r.db.Where("league_id = ?", leagueID).Find(&payments).Error; err != nil {
		return nil, err
	}

	var domainPayments []models.Payment
	for _, p := range payments {
		domainPayments = append(domainPayments, models.Payment{
			ID:            p.ID,
			AmountInCents: p.AmountInCents,
			UserID:        p.UserID,
			LeagueID:      p.LeagueID,
			ReceivedAt:    p.ReceivedAt,
		})
	}
	return domainPayments, nil
}

func (r *PaymentsRepository) Add(payment *models.Payment) error {
	p := &persistence.Payment{
		AmountInCents: payment.AmountInCents,
		UserID:        payment.UserID,
		LeagueID:      payment.LeagueID,
		ReceivedAt:    payment.ReceivedAt,
		Method:        string(payment.Method),
		RecordedBy:    payment.RecordedBy,
		Reference:     payment.Reference,
	}
	return r.db.Create(p).Error
}

func (r *PaymentsRepository) AddAllocation(paymentID string, allocation *models.PaymentAllocation) error {
	a := &persistence.PaymentAllocation{
		PaymentID:     paymentID,
		GameChargeID:  allocation.GameChargeID,
		AmountInCents: allocation.AmountInCents,
	}
	return r.db.Create(a).Error
}
