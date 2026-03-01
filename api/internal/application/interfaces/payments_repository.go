package interfaces

import "github.com/bouncy/bouncy-api/internal/domain/models"

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
	ClaimUnclaimedRecords(userID string, externalName string) error
}
