package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	emptyString = ""
)

type (
	UserClaims struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		jwt.RegisteredClaims
	}
)

func NewToken(secretKey string) (string, error) {
	userClaims := UserClaims{
		Username: "me",
		Password: "pass",
		Email:    "email@will.be.here",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "uuid_here",
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	err := userClaims.validate()
	if err != nil {
		return emptyString, fmt.Errorf("jwt_manager.GetRoken - incorrect claims: %w", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	s, err := token.SignedString(secretKey)
	if err != nil {
		return emptyString, fmt.Errorf("jwt_manager.GetRoken - can't signed token: %w", err)
	}
	return s, nil
}

func (c *UserClaims) validate() error {
	//TODO validate
	return nil
}
