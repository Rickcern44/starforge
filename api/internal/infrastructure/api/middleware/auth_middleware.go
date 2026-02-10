package middleware

import (
	"context"

	"net/http"
	"strings"

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/contract"
	"github.com/bouncy/bouncy-api/internal/infrastructure/auth"
	"github.com/bouncy/bouncy-api/internal/infrastructure/utils"
)

// AuthMiddleware provides JWT authentication for routes.
func AuthMiddleware(jwtService *application.JwtService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "Authorization header required"})
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
				utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "Authorization header format must be Bearer {token}"})
				return
			}

			tokenString := headerParts[1]
			claims, err := jwtService.Validate(tokenString)
			if err != nil {
				utils.WriteJSON(w, http.StatusUnauthorized, contract.ErrorResponse{Error: "Invalid or expired token"})
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), auth.ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
