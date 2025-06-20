package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtClaims struct {
	UserId uuid.UUID `json:"user_id"`
	TeamId uuid.UUID `json:"team_id"`
	Name   string    `json:"name"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}
