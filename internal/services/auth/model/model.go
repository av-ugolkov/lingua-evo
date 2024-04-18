package model

import "github.com/google/uuid"

type (
	CreateSessionRq struct {
		User        string `json:"user"`
		Password    string `json:"password"`
		Fingerprint string `json:"fingerprint"`
	}

	CreateUserRq struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Code     int    `json:"code"`
	}

	CreateCodeRq struct {
		Email string `json:"email"`
	}

	CreateSessionRs struct {
		AccessToken string `json:"access_token"`
	}

	CreateUserRs struct {
		UserID uuid.UUID `json:"user_id"`
	}
)
