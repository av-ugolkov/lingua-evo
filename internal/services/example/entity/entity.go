package entity

import "github.com/google/uuid"

type Example struct {
	Id        uuid.UUID
	Native    string
	Translate string
}
