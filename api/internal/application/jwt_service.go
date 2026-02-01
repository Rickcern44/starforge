package application

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bouncy/bouncy-api/internal/infrastructure/auth"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	secret []byte
}

func NewJwtService(secret string) *JwtService {
	return &JwtService{
		secret: []byte(secret),
	}
}

func (s JwtService) GenerateToken(claims *auth.Claims) (string, error) {
	mappedClaims := jwt.MapClaims{
		"sub":   claims.UserId,
		"email": claims.Email,
		"roles": claims.Roles,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mappedClaims)
	return token.SignedString(s.secret)
}

func (s JwtService) Validate(tokenString string) (*auth.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	var roles []string
	if rawRoles, ok := claimsMap["roles"].([]interface{}); ok {
		for _, r := range rawRoles {
			if roleStr, isString := r.(string); isString {
				roles = append(roles, roleStr)
			}
		}
	} else if roleStr, ok := claimsMap["roles"].(string); ok {
		// Fallback for comma-separated string if it was somehow stored this way
		roles = strings.Split(roleStr, ",")
	}

	return &auth.Claims{
		UserId: claimsMap["sub"].(string),
		Email:  claimsMap["email"].(string),
		Roles:  roles, // Assign the correctly parsed roles
	}, nil
}
