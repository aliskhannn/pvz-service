package token

import (
	"github.com/aliskhannn/pvz-service/internal/infrastructure/jwt"
	"github.com/google/uuid"
)

type Generator interface {
	CreateToken(userId uuid.UUID, role string) (string, error)
	ValidateToken(tokenString string) (*jwt.Claims, error)
}
