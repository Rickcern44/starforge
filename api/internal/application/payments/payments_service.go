package payments

import (
	"log/slog"
	"time"

	"github.com/bouncy/bouncy-api/internal/application/interfaces"
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/google/uuid"
)

type Service struct {
	repo interfaces.PaymentsRepository
}

func NewPaymentsService(repo interfaces.PaymentsRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListByLeague(leagueID string) ([]models.Payment, error) {
	return s.repo.ListByLeague(leagueID)
}

func (s *Service) ListPaymentsByUser(userID string) ([]models.Payment, error) {
	return s.repo.ListPaymentsByUser(userID)
}

func (s *Service) ListPaymentsByExternalName(name string) ([]models.Payment, error) {
	return s.repo.ListPaymentsByExternalName(name)
}

func (s *Service) Add(payment *models.Payment) error {
	if payment.ID == "" {
		payment.ID = uuid.NewString()
	}
	if payment.ReceivedAt.IsZero() {
		payment.ReceivedAt = time.Now()
	}

	if err := s.repo.Add(payment); err != nil {
		return err
	}

	// Auto-allocate if the payment is linked to a user
	if payment.UserID != nil {
		unpaidCharges, err := s.repo.ListUnpaidChargesByUser(*payment.UserID)
		if err != nil {
			slog.Error("Failed to list unpaid charges for auto-allocation", "userId", *payment.UserID, "error", err)
			return nil // We don't fail the payment add if allocation fails
		}

		remainingPayment := payment.AmountInCents
		for _, charge := range unpaidCharges {
			if remainingPayment <= 0 {
				break
			}

			// Calculate how much is left to pay on this charge
			alreadyPaid := 0
			for _, a := range charge.Allocations {
				alreadyPaid += a.AmountInCents
			}
			needed := charge.AmountCents - alreadyPaid

			allocationAmount := needed
			if remainingPayment < needed {
				allocationAmount = remainingPayment
			}

			allocation := &models.PaymentAllocation{
				PaymentID:     payment.ID,
				GameChargeID:  charge.ID,
				AmountInCents: allocationAmount,
			}

			if err := s.repo.AddAllocation(payment.ID, allocation); err != nil {
				slog.Error("Failed to auto-allocate payment", "paymentId", payment.ID, "chargeId", charge.ID, "error", err)
				continue
			}

			remainingPayment -= allocationAmount
		}
	}

	return nil
}

func (s *Service) AddAllocation(paymentID string, allocation *models.PaymentAllocation) error {
	return s.repo.AddAllocation(paymentID, allocation)
}

func (s *Service) CreateCharge(charge *models.GameCharge) error {
	return s.repo.CreateCharge(charge)
}

func (s *Service) ListChargesByUser(userID string) ([]models.GameCharge, error) {
	return s.repo.ListChargesByUser(userID)
}

func (s *Service) ListChargesByExternalName(name string) ([]models.GameCharge, error) {
	return s.repo.ListChargesByExternalName(name)
}

func (s *Service) ClaimUnclaimedRecords(userID string, externalName string) error {
	return s.repo.ClaimUnclaimedRecords(userID, externalName)
}
