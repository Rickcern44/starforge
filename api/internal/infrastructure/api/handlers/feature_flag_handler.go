package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application/features"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/contract"
	"github.com/bouncy/bouncy-api/internal/infrastructure/utils"
	"github.com/go-chi/chi/v5"
)

type FeatureFlagHandler struct {
	service *features.FeatureFlagService
}

func NewFeatureFlagHandler(service *features.FeatureFlagService) *FeatureFlagHandler {
	return &FeatureFlagHandler{service: service}
}

func RegisterFeatureFlagRoutes(r chi.Router, handler *FeatureFlagHandler) {
	r.Get("/admin/features", handler.GetAllFeatureFlags)
	r.Patch("/admin/features/{key}", handler.UpdateFeatureFlag)
}

func (h *FeatureFlagHandler) GetAllFeatureFlags(w http.ResponseWriter, r *http.Request) {
	flags, err := h.service.GetAll()
	if err != nil {
		slog.Error("Get all feature flags service failed", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: "failed to get feature flags"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, flags)
}

type UpdateFeatureFlagRequest struct {
	Enabled bool `json:"enabled"`
}

func (h *FeatureFlagHandler) UpdateFeatureFlag(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	var req UpdateFeatureFlagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode update feature flag request", "key", key, "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: "Invalid request body"})
		return
	}

	if err := h.service.Update(key, req.Enabled); err != nil {
		slog.Error("Update feature flag service failed", "key", key, "enabled", req.Enabled, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: "Failed to update feature flag"})
		return
	}

	w.WriteHeader(http.StatusOK)
}
