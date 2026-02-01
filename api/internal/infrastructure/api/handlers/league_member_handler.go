package handlers

import (
	"encoding/json"
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

func RegisterLeagueMemberHandlers(r chi.Router, handler *LeagueMemberHandler) {
	r.Get("/league/{leagueId}/members", handler.ListMembers)
	r.Post("/league/{leagueId}/members", handler.AddMember)
	r.Patch("/league/{leagueId}/members/{memberId}", handler.UpdateRole)
	r.Delete("/league/{leagueId}/members/{memberId}", handler.RemoveMember)
}

func (h *LeagueMemberHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")

	members, err := h.service.ListMembers(leagueId)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, members)
}

type AddMemberRequest struct {
	UserId       string      `json:"userId"`
	AddingUserId string      `json:"addingUserId"`
	Role         models.Role `json:"role"`
}

func (h *LeagueMemberHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	var request AddMemberRequest
	leagueId := chi.URLParam(r, "leagueId")

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.AddMember(leagueId, request.AddingUserId, request.UserId, request.Role); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *LeagueMemberHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")
	memberId := chi.URLParam(r, "memberId")
	// For actual implementation, you'd likely get the `removingUserId` from context

	if err := h.service.RemoveMember(leagueId, memberId); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

type UpdateMemberRoleRequest struct {
	Role models.Role `json:"role"`
}

func (h *LeagueMemberHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	leagueId := chi.URLParam(r, "leagueId")
	memberId := chi.URLParam(r, "memberId")
	// For actual implementation, you'd likely get the `updatingUserId` from context

	var request UpdateMemberRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, contract.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.UpdateRole(leagueId, memberId, request.Role); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
