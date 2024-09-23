package repository

import (
	"context"
	"fmt"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"
	"github.com/av-ugolkov/lingua-evo/internal/services/user/delivery/repository"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"testing"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/db/postgres"
	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetVocabulariesWithMaxWords(t *testing.T) {
	ctx := context.Background()

	tp := postgres.NewTempPostgres(ctx, "../../../../..")
	defer tp.DropDB(ctx)

	if tp == nil {
		t.Fatal("can't init container for DB")
	}

	repo := NewRepo(tp.PgxPool)
	userRepo := repository.NewRepo(tp.PgxPool)
	tr := transactor.NewTransactor(tp.PgxPool)

	t.Run("empty vocabularies", func(t *testing.T) {
		vocabs, err := repo.GetVocabulariesWithMaxWords(ctx, 3, []uint8{1, 2})
		if err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, 0, len(vocabs))
	})
	t.Run("vocabularies are empty", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			uid, err := userRepo.AddUser(ctx, &entityUser.User{
				ID:           uuid.New(),
				Name:         "test_user",
				Email:        "test_user@email.com",
				Role:         runtime.User,
				PasswordHash: "qwerty",
			})
			if err != nil {
				return err
			}

			for i := 0; i < 10; i++ {
				_, err := repo.AddVocab(ctx, entity.Vocab{
					ID:            uuid.New(),
					UserID:        uid,
					Name:          fmt.Sprintf("test_%d", i),
					Access:        uint8(access.Subscribers),
					NativeLang:    "en",
					TranslateLang: "ru",
					Description:   "",
					Tags:          nil,
					CreatedAt:     time.Now().UTC(),
					UpdatedAt:     time.Now().UTC(),
				}, nil)
				if err != nil {
					return err
				}
			}

			vocabs, err := repo.GetVocabulariesWithMaxWords(ctx, 3, []uint8{1, 2})
			if err != nil {
				return err
			}

			assert.Equal(t, 0, len(vocabs))

			return nil
		})

		assert.NoError(t, err)
	})
	t.Run("get vocabularies with max count words", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			uid, err := userRepo.AddUser(ctx, &entityUser.User{
				ID:           uuid.New(),
				Name:         "test_user",
				Email:        "test_user@email.com",
				Role:         runtime.User,
				PasswordHash: "qwerty",
			})
			if err != nil {
				return err
			}

			for i := 0; i < 10; i++ {
				vid, err := repo.AddVocab(ctx, entity.Vocab{
					ID:            uuid.New(),
					UserID:        uid,
					Name:          fmt.Sprintf("test_%d", i),
					Access:        uint8(access.Subscribers),
					NativeLang:    "en",
					TranslateLang: "ru",
					Description:   "",
					Tags:          nil,
					CreatedAt:     time.Now().UTC(),
					UpdatedAt:     time.Now().UTC(),
				}, nil)
				if err != nil {
					return err
				}

				for j := 0; j < 10-i; j++ {
					_, err := repo.AddWord(ctx, entity.VocabWord{
						VocabID:       vid,
						ID:            uuid.New(),
						NativeID:      uuid.New(),
						TranslateIDs:  []uuid.UUID{},
						Pronunciation: "",
					})
					if err != nil {
						return err
					}
				}
			}

			vocabs, err := repo.GetVocabulariesWithMaxWords(ctx, 3, []uint8{1, 2})
			if err != nil {
				return err
			}

			assert.Equal(t, 3, len(vocabs))
			assert.Equal(t, uint(10), vocabs[0].WordsCount)
			assert.Equal(t, uint(9), vocabs[1].WordsCount)
			assert.Equal(t, uint(8), vocabs[2].WordsCount)

			return nil
		})

		assert.NoError(t, err)
	})
}
