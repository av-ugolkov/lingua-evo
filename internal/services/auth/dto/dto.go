package dto

import "github.com/google/uuid"

type (
	CreateSessionRq struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}

	CreateSessionRs struct {
		AccessToken  string    `json:"access_token"`
		RefreshToken uuid.UUID `json:"refresh_token"`
	}
)
