package entity

import "github.com/google/uuid"

type (
	Tag struct {
		ID   uuid.UUID
		Text string
	}
)
