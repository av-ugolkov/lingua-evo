package model

import "github.com/google/uuid"

type (
	TagRs struct {
		ID   uuid.UUID `json:"id"`
		Text string    `json:"text"`
	}
)
