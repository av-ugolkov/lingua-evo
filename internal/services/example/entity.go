package example

import "github.com/google/uuid"

type (
	Example struct {
		Id   uuid.UUID
		Text string
	}
)
