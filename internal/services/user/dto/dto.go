package dto

import "github.com/google/uuid"

type (
	GetIDRq struct {
		Value string `json:"value"`
	}

	UserIDRs struct {
		ID uuid.UUID `json:"id"`
	}
)
