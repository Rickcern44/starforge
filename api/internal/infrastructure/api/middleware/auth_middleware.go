package middleware

import (
	"context"

	"net/http"
	"strings"

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/infrastructure/api/handlers" // Assuming WriteJSON is here
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const UserIDContextKey contextKey = "userID"

// AuthMiddleware provides JWT authentication for routes.
func AuthMiddleware(jwtService *application.JwtService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handlers.WriteJSON(w, http.StatusUnauthorized, handlers.ErrorResponse{Error: "Authorization header required"})
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
				handlers.WriteJSON(w, http.StatusUnauthorized, handlers.ErrorResponse{Error: "Authorization header format must be Bearer {token}"})
				return
			}

			tokenString := headerParts[1]
			claims, err := jwtService.Validate(tokenString)
			if err != nil {
				handlers.WriteJSON(w, http.StatusUnauthorized, handlers.ErrorResponse{Error: "Invalid or expired token"})
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
