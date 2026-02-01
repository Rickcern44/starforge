package middleware

import (
	"net/http"

	"github.com/bouncy/bouncy-api/internal/infrastructure/api/contract"
	"github.com/bouncy/bouncy-api/internal/infrastructure/auth"
	"github.com/bouncy/bouncy-api/internal/infrastructure/utils"
)

// RoleMiddleware checks if the authenticated user has the required role.
func RoleMiddleware(requiredRole string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(auth.ClaimsContextKey).(*auth.Claims)
			if !ok || claims == nil {
				utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "authentication required"})
				return
			}

			if !hasRole(claims.Roles, requiredRole) {
				utils.WriteJSON(w, http.StatusForbidden, contract.ErrorResponse{Error: "insufficient permissions"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func hasRole(userRoles []string, requiredRole string) bool {
	for _, role := range userRoles {
		if role == requiredRole {
			return true
		}
	}
	return false
}
