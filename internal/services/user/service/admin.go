package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Service) CheckAdmin(ctx context.Context, uid uuid.UUID) (bool, error) {
	usr, err := s.GetUserByID(ctx, uid)
	if err != nil {
		return false, fmt.Errorf("user.Service.CheckAdmin: %w", err)
	}

	return usr.Role.IsAdmin(), nil
}
