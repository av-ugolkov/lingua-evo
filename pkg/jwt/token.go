package jwt

import (
	"fmt"

	"lingua-evo/internal/config"
	entitySession "lingua-evo/internal/services/auth/entity"
	entityUser "lingua-evo/internal/services/user/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	emptyString = ""
)

type (
	UserClaims struct {
		ID uuid.UUID
		jwt.RegisteredClaims
	}
)

func NewJWTToken(u *entityUser.User, s *entitySession.Claims) (string, error) {
	userClaims := UserClaims{
		ID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        s.ID.String(),
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(s.ExpiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	t, err := token.SignedString([]byte(config.GetConfig().JWT.Secret))
	if err != nil {
		return emptyString, fmt.Errorf("pkg.jwt.token.NewJWTToken - can't signed token: %w", err)
	}
	return t, nil
}
