package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type (
	UserClaims struct {
		UserID uuid.UUID `json:"user_id"`
		jwt.RegisteredClaims
	}
)

func NewJWTToken(uid, sid uuid.UUID, expiresAt time.Time) (string, error) {
	userClaims := UserClaims{
		UserID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        sid.String(),
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(&jwt.SigningMethodRSA{}, userClaims)
	t, err := token.SignedString([]byte(config.GetConfig().JWT.Secret))
	if err != nil {
		return runtime.EmptyString, fmt.Errorf("pkg.jwt.token.NewJWTToken - can't signed token: %w", err)
	}
	return t, nil
}

func ValidateJWT(tokenStr string, secret string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, fmt.Errorf("pkg.jwt.token.ValidateToken - token expired")
	} else if err != nil {
		return nil, fmt.Errorf("pkg.jwt.token.ValidateToken - can't parse token [%s]: %w", tokenStr, err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("pkg.jwt.token.ValidateToken - invalid token")
	}

	return claims, nil
}
