package services

import (
	"lingua-evo/internal/delivery/repository"
	"lingua-evo/pkg/logging"
)

type Lingua struct {
	logger *logging.Logger
	db     *repository.Database
}

func NewLinguaService(logger *logging.Logger, db *repository.Database) *Lingua {
	return &Lingua{
		logger: logger,
		db:     db,
	}
}
