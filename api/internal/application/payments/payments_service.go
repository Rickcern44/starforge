package payments

import (
	"github.com/bouncy/bouncy-api/internal/application/interfaces"
	"github.com/bouncy/bouncy-api/internal/domain/models"
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

func (s *Service) Add(payment *models.Payment) error {
	return s.repo.Add(payment)
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
