package notifications

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_SetVocabNotification(t *testing.T) {
	t.Parallel()
	t.Run("GetNotification Error SetVocabNotification", func(t *testing.T) {
		var (
			ctx = context.Background()
			uid = uuid.New()
			vid = uuid.New()
		)
		repoNotificationMock := new(mockRepoNotification)
		repoNotificationMock.On("GetVocabNotification", ctx, uid, vid).Return(false, errors.New("some error"))
		repoNotificationMock.On("SetVocabNotification", ctx, uid, vid).Return(nil)
		repoNotificationMock.On("DeleteVocabNotification", ctx, uid, vid).Return(nil)

		s := &Service{repo: repoNotificationMock}

		got, err := s.SetVocabNotification(ctx, uid, vid)
		assert.Error(t, err)
		assert.Equal(t, false, got)
	})
	t.Run("SetNotification Error SetVocabNotification", func(t *testing.T) {
		var (
			ctx = context.Background()
			uid = uuid.New()
			vid = uuid.New()
		)
		repoNotificationMock := new(mockRepoNotification)
		repoNotificationMock.On("GetVocabNotification", ctx, uid, vid).Return(false, nil)
		repoNotificationMock.On("SetVocabNotification", ctx, uid, vid).Return(errors.New("some error"))
		repoNotificationMock.On("DeleteVocabNotification", ctx, uid, vid).Return(nil)

		s := &Service{repo: repoNotificationMock}

		got, err := s.SetVocabNotification(ctx, uid, vid)
		assert.Error(t, err)
		assert.Equal(t, false, got)
	})
	t.Run("SetNotification SetVocabNotification", func(t *testing.T) {
		var (
			ctx = context.Background()
			uid = uuid.New()
			vid = uuid.New()
		)
		repoNotificationMock := new(mockRepoNotification)
		repoNotificationMock.On("GetVocabNotification", ctx, uid, vid).Return(false, nil)
		repoNotificationMock.On("SetVocabNotification", ctx, uid, vid).Return(nil)
		repoNotificationMock.On("DeleteVocabNotification", ctx, uid, vid).Return(nil)

		s := &Service{repo: repoNotificationMock}

		got, err := s.SetVocabNotification(ctx, uid, vid)
		assert.NoError(t, err)
		assert.Equal(t, true, got)
	})
	t.Run("DeleteNotification Error SetVocabNotification", func(t *testing.T) {
		var (
			ctx = context.Background()
			uid = uuid.New()
			vid = uuid.New()
		)
		repoNotificationMock := new(mockRepoNotification)
		repoNotificationMock.On("GetVocabNotification", ctx, uid, vid).Return(true, nil)
		repoNotificationMock.On("SetVocabNotification", ctx, uid, vid).Return(nil)
		repoNotificationMock.On("DeleteVocabNotification", ctx, uid, vid).Return(errors.New("some error"))

		s := &Service{repo: repoNotificationMock}

		got, err := s.SetVocabNotification(ctx, uid, vid)
		assert.Error(t, err)
		assert.Equal(t, false, got)
	})
	t.Run("DeleteNotification SetVocabNotification", func(t *testing.T) {
		var (
			ctx = context.Background()
			uid = uuid.New()
			vid = uuid.New()
		)
		repoNotificationMock := new(mockRepoNotification)
		repoNotificationMock.On("GetVocabNotification", ctx, uid, vid).Return(true, nil)
		repoNotificationMock.On("SetVocabNotification", ctx, uid, vid).Return(nil)
		repoNotificationMock.On("DeleteVocabNotification", ctx, uid, vid).Return(nil)

		s := &Service{repo: repoNotificationMock}

		got, err := s.SetVocabNotification(ctx, uid, vid)
		assert.NoError(t, err)
		assert.Equal(t, false, got)
	})
}
