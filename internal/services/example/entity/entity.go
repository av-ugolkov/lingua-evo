package entity

import "github.com/google/uuid"

type Example struct {
	Id        uuid.UUID
	Original  string
	Translate string
}
