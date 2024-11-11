package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *UserRepo) UpdatePsw(ctx context.Context, uid uuid.UUID, hashPsw string) (err error) {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2`

	_, err = r.tr.Exec(ctx, query, hashPsw, uid)
	if err != nil {
		return fmt.Errorf("user.repository.UserRepo.UpdatePsw: %w", err)
	}

	return nil
}

func (r *UserRepo) UpdateEmail(ctx context.Context, uid uuid.UUID, newEmail string) (err error) {
	query := `UPDATE users SET email = $1 WHERE id = $2`

	_, err = r.tr.Exec(ctx, query, newEmail, uid)
	if err != nil {
		return fmt.Errorf("user.repository.UserRepo.UpdateEmail: %w", err)
	}

	return nil
}
