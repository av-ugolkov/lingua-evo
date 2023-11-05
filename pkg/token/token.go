package token

import (
	"errors"
	"fmt"

	"lingua-evo/internal/config"
	entitySession "lingua-evo/internal/services/auth/entity"

	"github.com/golang-jwt/jwt/v5"
)

const (
	emptyString = ""
)

type (
	UserClaims struct {
		Email           string
		HashFingerprint string

		jwt.RegisteredClaims
	}
)

func NewJWTToken(s *entitySession.Claims) (string, error) {
	userClaims := UserClaims{
		Email:           s.Email,
		HashFingerprint: s.HashFingerprint,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        s.ID.String(),
			Subject:   s.UserID.String(),
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(s.ExpiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	secret := config.GetConfig().JWT.Secret
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return emptyString, fmt.Errorf("pkg.jwt.token.NewJWTToken - can't signed token: %w", err)
	}
	return t, nil
}

func ValidateJWT(tokenStr string, secret string) (*UserClaims, error) {
	var claims UserClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if !token.Valid {
		return nil, fmt.Errorf("pkg.jwt.token.ValidateToken - invalid token")
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil, fmt.Errorf("pkg.jwt.token.ValidateToken - malformed token")
	} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return nil, fmt.Errorf("pkg.jwt.token.ValidateToken - invalid signature")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return nil, fmt.Errorf("pkg.jwt.token.ValidateToken - expired token")
	} else if err != nil {
		return nil, fmt.Errorf("pkg.jwt.token.ValidateToken: %w", err)
	}

	return &claims, nil
}
