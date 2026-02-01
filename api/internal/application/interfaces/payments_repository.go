package interfaces

import "github.com/bouncy/bouncy-api/internal/domain/models"

type PaymentsRepository interface {
	ListByLeague(leagueID string) ([]models.Payment, error)
	Add(payment *models.Payment) error
	AddAllocation(paymentID string, allocation *models.PaymentAllocation) error
}
