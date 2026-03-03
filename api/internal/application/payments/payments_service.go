package payments

import (
	"log/slog"
	"time"

	"github.com/bouncy/bouncy-api/internal/domain/interfaces"
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

	slog.Info("Attempting auto-allocation for payment", "paymentId", payment.ID, "userId", payment.UserID, "externalName", payment.ExternalName)

	var unpaidCharges []models.GameCharge
	var err error

	// Auto-allocate: Prefer UserID, fallback to ExternalName
	if payment.UserID != nil && *payment.UserID != "" {
		unpaidCharges, err = s.repo.ListUnpaidChargesByUser(*payment.UserID)
	} else if payment.ExternalName != "" {
		unpaidCharges, err = s.repo.ListUnpaidChargesByExternalName(payment.ExternalName)
	}

	if err != nil {
		slog.Error("Failed to list unpaid charges for auto-allocation", "paymentId", payment.ID, "error", err)
		return nil
	}

	slog.Info("Found unpaid charges for allocation", "count", len(unpaidCharges))

	if len(unpaidCharges) > 0 {
		s.allocateToCharges(payment, unpaidCharges)
	}

	return nil
}

func (s *Service) allocateToCharges(payment *models.Payment, charges []models.GameCharge) {
	remainingPayment := payment.AmountCents
	for _, charge := range charges {
		if remainingPayment <= 0 {
			break
		}

		// Calculate how much is left to pay on this charge
		alreadyPaid := 0
		for _, a := range charge.Allocations {
			alreadyPaid += a.AmountInCents
		}
		needed := charge.AmountCents - alreadyPaid

		slog.Info("Processing charge for allocation", "chargeId", charge.ID, "neededCents", needed, "remainingPaymentCents", remainingPayment)

		if needed <= 0 {
			continue
		}

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

func (s *Service) GetFinancialSummary(leagueID string) (*models.LeagueFinancialSummary, error) {
	return s.repo.GetFinancialSummary(leagueID)
}
