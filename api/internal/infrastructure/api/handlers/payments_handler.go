package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application/payments"
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/contract"
	"github.com/bouncy/bouncy-api/internal/infrastructure/utils"
	"github.com/go-chi/chi/v5"
)

type PaymentsHandler struct {
	service *payments.Service
}

func NewPaymentsHandler(service *payments.Service) *PaymentsHandler {
	return &PaymentsHandler{service: service}
}

// RegisterPaymentsRoutes registers the payment related routes.
func RegisterPaymentsRoutes(r chi.Router, handler *PaymentsHandler) {
	r.Get("/league/{leagueId}/payments", handler.ListPayments)
	r.Post("/league/{leagueId}/payments", handler.AddPayment)
	r.Post("/payments/{paymentId}/allocations", handler.AddAllocation)
}

// ListPayments handles listing all payments for a specific league.
// @Summary List all payments for a specific league
// @Description Retrieves a list of all payments associated with the given league ID.
// @Tags payments
// @Accept json
// @Produce json
// @Param leagueId path string true "ID of the league"
// @Security BearerAuth
// @Success 200 {array} models.Payment "List of payments"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /league/{leagueId}/payments [get]
func (h *PaymentsHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")

	payments, err := h.service.ListByLeague(leagueId)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, payments)
}

// AddPayment handles adding a new payment to a league.
// @Summary Add a new payment to a league
// @Description Creates a new payment entry for the specified league.
// @Tags payments
// @Accept json
// @Produce json
// @Param leagueId path string true "ID of the league"
// @Param request body models.Payment true "Payment creation details"
// @Security BearerAuth
// @Success 201 {object} object "Payment created successfully"
// @Failure 400 {object} contract.ErrorResponse "Invalid request body"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /league/{leagueId}/payments [post]
func (h *PaymentsHandler) AddPayment(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")

	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: err.Error()})
		return
	}
	payment.LeagueID = leagueId

	if err := h.service.Add(&payment); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// AddAllocation handles adding an allocation to a payment.
// @Summary Add an allocation to a payment
// @Description Adds a new allocation to an existing payment.
// @Tags payments
// @Accept json
// @Produce json
// @Param paymentId path string true "ID of the payment"
// @Param request body models.PaymentAllocation true "Payment allocation details"
// @Security BearerAuth
// @Success 201 {object} object "Allocation created successfully"
// @Failure 400 {object} contract.ErrorResponse "Invalid request body"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /payments/{paymentId}/allocations [post]
func (h *PaymentsHandler) AddAllocation(w http.ResponseWriter, r *http.Request) {
	paymentId := chi.URLParam(r, "paymentId")

	var allocation models.PaymentAllocation
	if err := json.NewDecoder(r.Body).Decode(&allocation); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.AddAllocation(paymentId, &allocation); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
