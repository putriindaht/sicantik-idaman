package token

import (
	"sicantik-idaman/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(jwtKey string, claims *domain.JwtClaims) (string, error) {
	exp := time.Now().Add(24 * time.Hour)

	claims.RegisteredClaims = jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(exp)}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}
