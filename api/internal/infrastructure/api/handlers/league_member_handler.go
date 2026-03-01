package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bouncy/bouncy-api/internal/application/leagues"
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/contract"
	"github.com/bouncy/bouncy-api/internal/infrastructure/utils"
	"github.com/go-chi/chi/v5"
)

type LeagueMemberHandler struct {
	service *leagues.LeagueMemberService
}

func NewLeagueMemberHandler(service *leagues.LeagueMemberService) *LeagueMemberHandler {
	return &LeagueMemberHandler{service: service}
}

// RegisterLeagueMemberHandlers registers the league member related routes.
func RegisterLeagueMemberHandlers(r chi.Router, handler *LeagueMemberHandler) {
	r.Get("/league/{leagueId}/members", handler.ListMembers)
	r.Post("/league/{leagueId}/members", handler.AddMember)
	r.Patch("/league/{leagueId}/members/{memberId}", handler.UpdateRole)
	r.Delete("/league/{leagueId}/members/{memberId}", handler.RemoveMember)
}

// ListMembers handles listing all members of a league.
// @Summary List all members of a league
// @Description Retrieves a list of all members associated with the given league ID.
// @Tags league_members
// @Accept json
// @Produce json
// @Param leagueId path string true "ID of the league"
// @Security BearerAuth
// @Success 200 {array} models.LeagueMember "List of league members"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /league/{leagueId}/members [get]
func (h *LeagueMemberHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")

	members, err := h.service.ListMembers(leagueId)
	if err != nil {
		slog.Error("List members failed", "leagueId", leagueId, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, members)
}

// AddMemberRequest represents the request body for adding a new member to a league.
type AddMemberRequest struct {
	UserId       string      `json:"userId"`
	AddingUserId string      `json:"addingUserId"`
	Role         models.Role `json:"role"`
}

// AddMember handles adding a new member to a league.
// @Summary Add a new member to a league
// @Description Adds a new member to the specified league.
// @Tags league_members
// @Accept json
// @Produce json
// @Param leagueId path string true "ID of the league"
// @Param request body AddMemberRequest true "Member details to add"
// @Security BearerAuth
// @Success 201 {object} object "Member added successfully"
// @Failure 400 {object} contract.ErrorResponse "Invalid request body"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /league/{leagueId}/members [post]
func (h *LeagueMemberHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	var request AddMemberRequest
	leagueId := chi.URLParam(r, "leagueId")

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.Error("Failed to decode add member request", "leagueId", leagueId, "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.AddMember(leagueId, request.AddingUserId, request.UserId, request.Role); err != nil {
		slog.Error("Add member service failed", "leagueId", leagueId, "userId", request.UserId, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// RemoveMember handles removing a member from a league.
// @Summary Remove a member from a league
// @Description Removes a member from the specified league.
// @Tags league_members
// @Accept json
// @Produce json
// @Param leagueId path string true "ID of the league"
// @Param memberId path string true "ID of the member to remove"
// @Security BearerAuth
// @Success 200 {object} object "Member removed successfully"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /league/{leagueId}/members/{memberId} [delete]
func (h *LeagueMemberHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")
	memberId := chi.URLParam(r, "memberId")
	// For actual implementation, you'd likely get the `removingUserId` from context

	if err := h.service.RemoveMember(leagueId, memberId); err != nil {
		slog.Error("Remove member failed", "leagueId", leagueId, "memberId", memberId, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateMemberRoleRequest represents the request body for updating a member's role in a league.
type UpdateMemberRoleRequest struct {
	Role models.Role `json:"role"`
}

// UpdateRole handles updating a league member's role.
// @Summary Update a league member's role
// @Description Updates the role of a specific league member.
// @Tags league_members
// @Accept json
// @Produce json
// @Param leagueId path string true "ID of the league"
// @Param memberId path string true "ID of the member to update"
// @Param request body UpdateMemberRoleRequest true "New role for the member"
// @Security BearerAuth
// @Success 200 {object} object "Member role updated successfully"
// @Failure 400 {object} contract.ErrorResponse "Invalid request body"
// @Failure 500 {object} contract.ErrorResponse "Internal server error"
// @Router /league/{leagueId}/members/{memberId} [patch]
func (h *LeagueMemberHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")
	memberId := chi.URLParam(r, "memberId")
	// For actual implementation, you'd likely get the `updatingUserId` from context

	var request UpdateMemberRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.Error("Failed to decode update role request", "leagueId", leagueId, "memberId", memberId, "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.UpdateRole(leagueId, memberId, request.Role); err != nil {
		slog.Error("Update member role failed", "leagueId", leagueId, "memberId", memberId, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
