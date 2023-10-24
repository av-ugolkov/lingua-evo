package dto

import "github.com/google/uuid"

type (
	CreateUserRq struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	CreateUserRs struct {
		UserID uuid.UUID `json:"user_id"`
	}

	GetIDRq struct {
		Value string `json:"value"`
	}

	UserIDRs struct {
		ID uuid.UUID `json:"id"`
	}
)
