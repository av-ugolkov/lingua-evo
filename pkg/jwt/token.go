package jwt

import (
	"fmt"
	"log/slog"

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
		ID uuid.UUID `json:"id"`
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

func ValidateToken(tokenStr string, secret string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return false, fmt.Errorf("pkg.jwt.token.ValidateToken - can't parse token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return false, fmt.Errorf("pkg.jwt.token.ValidateToken - invalid token")
	}

	slog.Info("claims: %v", claims)

	return true, nil
}
