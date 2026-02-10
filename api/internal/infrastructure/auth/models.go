package auth

import "github.com/bouncy/bouncy-api/internal/infrastructure/utils"

type Claims struct {
	UserId string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
}

// ClaimsContextKey is the key to retrieve JWT claims from the request context.
const ClaimsContextKey utils.ContextKey = "jwtClaims"
