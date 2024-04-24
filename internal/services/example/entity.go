package example

import "github.com/google/uuid"

type (
	Example struct {
		ID   uuid.UUID
		Text string
	}
)
