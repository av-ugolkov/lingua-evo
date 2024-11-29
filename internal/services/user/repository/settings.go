package repository

import (
	"context"
	"fmt"

	"github.com/av-ugolkov/lingua-evo/runtime"
	"github.com/google/uuid"
)

func (r *UserRepo) GetPswHash(ctx context.Context, uid uuid.UUID) (string, error) {
	const query = `SELECT password_hash FROM users WHERE id=$1`

	var pswHash string
	err := r.tr.QueryRow(ctx, query, uid).Scan(&pswHash)
	if err != nil {
		return runtime.EmptyString, fmt.Errorf("user.repository.UserRepo.GetPswHash: %w", err)
	}

	return pswHash, nil
}

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

func (r *UserRepo) UpdateNickname(ctx context.Context, uid uuid.UUID, newNickname string) (err error) {
	query := `UPDATE users SET nickname = $1 WHERE id = $2`

	_, err = r.tr.Exec(ctx, query, newNickname, uid)
	if err != nil {
		return fmt.Errorf("user.repository.UserRepo.UpdateNickname: %w", err)
	}

	return nil
}
