package runtime

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type contextUserIDKey struct{}

var (
	errUserIDNotFound = errors.New("user id not found")
)

func SetUserIDInContext(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, &contextUserIDKey{}, userID)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(&contextUserIDKey{}).(uuid.UUID)
	if !ok {
		return uuid.Nil, errUserIDNotFound
	}
	return userID, nil
}
