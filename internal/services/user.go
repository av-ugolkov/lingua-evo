package services

import (
	"context"

	"github.com/google/uuid"

	"lingua-evo/internal/delivery/repository"
)

func (l *Lingua) AddUser(ctx context.Context, u *repository.User) (uuid.UUID, error) {
	uid, err := l.db.AddUser(ctx, u)
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func (l *Lingua) EditUser(ctx context.Context, u *repository.User) error {
	return nil
}

func (l *Lingua) FindUser(ctx context.Context, username string) (uuid.UUID, error) {
	uid, err := l.db.FindUser(ctx, username)
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func (l *Lingua) FindEmail(ctx context.Context, email string) (uuid.UUID, error) {
	uid, err := l.db.FindUserByEmail(ctx, email)
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func (l *Lingua) RemoveUser(ctx context.Context, u *repository.User) error {
	return nil
}
