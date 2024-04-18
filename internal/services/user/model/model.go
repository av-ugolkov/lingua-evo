package model

import (
	"github.com/av-ugolkov/lingua-evo/runtime"
	"github.com/google/uuid"
)

type (
	GetValueRq struct {
		Value string `json:"value"`
	}

	UserRs struct {
		ID    uuid.UUID    `json:"id"`
		Name  string       `json:"name"`
		Email string       `json:"email"`
		Role  runtime.Role `json:"role"`
	}
)
