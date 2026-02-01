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

func RegisterPaymentsRoutes(r chi.Router, handler *PaymentsHandler) {
	r.Get("/league/{leagueId}/payments", handler.ListPayments)
	r.Post("/league/{leagueId}/payments", handler.AddPayment)
	r.Post("/payments/{paymentId}/allocations", handler.AddAllocation)
}

func (h *PaymentsHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")

	payments, err := h.service.ListByLeague(leagueId)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, payments)
}

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
