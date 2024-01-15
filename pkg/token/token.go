package token

import (
	"errors"
	"fmt"

	"lingua-evo/internal/config"
	entityAuth "lingua-evo/internal/services/auth"
	entityUser "lingua-evo/internal/services/user"
	"lingua-evo/runtime"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type (
	UserClaims struct {
		UserID uuid.UUID `json:"user_id"`
		jwt.RegisteredClaims
	}
)

func NewJWTToken(u *entityUser.User, s *entityAuth.Claims) (string, error) {
	userClaims := UserClaims{
		UserID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        s.ID.String(),
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(s.ExpiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	t, err := token.SignedString([]byte(config.GetConfig().JWT.Secret))
	if err != nil {
		return runtime.EmptyString, fmt.Errorf("pkg.jwt.token.NewJWTToken - can't signed token: %w", err)
	}
	return t, nil
}

func ValidateJWT(tokenStr string, secret string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, fmt.Errorf("pkg.jwt.token.ValidateToken - token expired")
	} else if err != nil {
		return nil, fmt.Errorf("pkg.jwt.token.ValidateToken - can't parse token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("pkg.jwt.token.ValidateToken - invalid token")
	}

	return claims, nil
}
