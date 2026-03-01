package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application/payments"
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/contract"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/middleware"
	"github.com/bouncy/bouncy-api/internal/infrastructure/auth"
	"github.com/bouncy/bouncy-api/internal/infrastructure/utils"
	"github.com/go-chi/chi/v5"
)

type PaymentsHandler struct {
	service *payments.Service
}

func NewPaymentsHandler(service *payments.Service) *PaymentsHandler {
	return &PaymentsHandler{service: service}
}

// RegisterUserPaymentsRoutes registers routes for a user to see their own data.
func RegisterUserPaymentsRoutes(r chi.Router, handler *PaymentsHandler) {
	r.Get("/me/payments", handler.GetCurrentUserPayments)
	r.Get("/me/charges", handler.GetCurrentUserCharges)
}

// RegisterPaymentsRoutes registers the payment related routes for league administration.
func RegisterPaymentsRoutes(r chi.Router, handler *PaymentsHandler) {
	r.Get("/league/{leagueId}/payments", handler.ListPayments)
	r.Post("/league/{leagueId}/payments", handler.AddPayment)
	r.Post("/payments/{paymentId}/allocations", handler.AddAllocation)

	// Admin routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.RoleMiddleware(string(models.RoleAdmin)))
		r.Post("/admin/payments/claim", handler.ClaimRecords)
	})
}

func (h *PaymentsHandler) GetCurrentUserPayments(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.Claims)
	if !ok {
		slog.Error("Get current user payments failed: invalid user context")
		utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "invalid user context"})
		return
	}

	payments, err := h.service.ListPaymentsByUser(claims.UserId)
	if err != nil {
		slog.Error("Get user payments service failed", "userId", claims.UserId, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, payments)
}

func (h *PaymentsHandler) GetCurrentUserCharges(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.Claims)
	if !ok {
		slog.Error("Get current user charges failed: invalid user context")
		utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "invalid user context"})
		return
	}

	charges, err := h.service.ListChargesByUser(claims.UserId)
	if err != nil {
		slog.Error("Get user charges service failed", "userId", claims.UserId, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, charges)
}

// ListPayments handles listing all payments for a specific league.
func (h *PaymentsHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")

	payments, err := h.service.ListByLeague(leagueId)
	if err != nil {
		slog.Error("List league payments failed", "leagueId", leagueId, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, payments)
}

// AddPayment handles adding a new payment to a league.
func (h *PaymentsHandler) AddPayment(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.Claims)
	if !ok {
		slog.Error("Add payment failed: invalid user context")
		utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "invalid user context"})
		return
	}

	leagueId := chi.URLParam(r, "leagueId")

	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		slog.Error("Failed to decode add payment request", "leagueId", leagueId, "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: err.Error()})
		return
	}
	payment.LeagueID = leagueId
	payment.RecordedBy = claims.UserId

	if err := h.service.Add(&payment); err != nil {
		slog.Error("Add payment service failed", "leagueId", leagueId, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// AddAllocation handles adding an allocation to a payment.
func (h *PaymentsHandler) AddAllocation(w http.ResponseWriter, r *http.Request) {
	paymentId := chi.URLParam(r, "paymentId")

	var allocation models.PaymentAllocation
	if err := json.NewDecoder(r.Body).Decode(&allocation); err != nil {
		slog.Error("Failed to decode add allocation request", "paymentId", paymentId, "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.AddAllocation(paymentId, &allocation); err != nil {
		slog.Error("Add allocation service failed", "paymentId", paymentId, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type ClaimRecordsRequest struct {
	UserID       string `json:"userId"`
	ExternalName string `json:"externalName"`
}

// ClaimRecords handles claiming unclaimed payments and charges for a user.
// @Summary Claim unclaimed payments and charges for a user
// @Description Links all unclaimed payments and charges with a specific external name to a user ID.
// @Tags payments
// @Accept json
// @Produce json
// @Param request body ClaimRecordsRequest true "Claim details"
// @Security BearerAuth
// @Success 200 {object} object "Records claimed successfully"
// @Failure 400 {object} contract.ErrorResponse "Invalid request body"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /admin/payments/claim [post]
func (h *PaymentsHandler) ClaimRecords(w http.ResponseWriter, r *http.Request) {
	var req ClaimRecordsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode claim records request", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.ClaimUnclaimedRecords(req.UserID, req.ExternalName); err != nil {
		slog.Error("Claim records service failed", "userId", req.UserID, "externalName", req.ExternalName, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Records claimed successfully"})
}
