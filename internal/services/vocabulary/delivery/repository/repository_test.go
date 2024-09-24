package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/db/postgres"
	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"
	"github.com/av-ugolkov/lingua-evo/internal/services/user/delivery/repository"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var errCancelTx = errors.New("transaction canceled")

func TestGetVocabulariesWithMaxWords(t *testing.T) {
	ctx := context.Background()

	tp := postgres.NewTempPostgres(ctx, "../../../../..")
	defer tp.DropDB(ctx)

	if tp == nil {
		t.Fatal("can't init container for DB")
	}

	tr := transactor.NewTransactor(tp.PgxPool)
	repo := NewRepo(tr)
	userRepo := repository.NewRepo(tr)

	t.Run("empty vocabularies", func(t *testing.T) {
		vocabs, err := repo.GetVocabulariesWithMaxWords(ctx, 3, []uint8{uint8(access.Public), uint8(access.Subscribers)})
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

			vocabs, err := repo.GetVocabulariesWithMaxWords(ctx, 3, []uint8{uint8(access.Public), uint8(access.Subscribers)})
			if err != nil {
				return err
			}

			assert.Equal(t, 0, len(vocabs))

			return errCancelTx
		})

		assert.ErrorIs(t, err, errCancelTx)
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

			vocabs, err := repo.GetVocabulariesWithMaxWords(ctx, 3, []uint8{uint8(access.Public), uint8(access.Subscribers)})
			if err != nil {
				return err
			}

			assert.Equal(t, 3, len(vocabs))
			assert.Equal(t, uint(10), vocabs[0].WordsCount)
			assert.Equal(t, uint(9), vocabs[1].WordsCount)
			assert.Equal(t, uint(8), vocabs[2].WordsCount)

			return errCancelTx
		})

		assert.ErrorIs(t, err, errCancelTx)
	})
}

func TestGetVocabsWithCountWords(t *testing.T) {
	ctx := context.Background()

	tp := postgres.NewTempPostgres(ctx, "../../../../..")
	defer tp.DropDB(ctx)

	if tp == nil {
		t.Fatal("can't init container for DB")
	}

	tr := transactor.NewTransactor(tp.PgxPool)
	repo := NewRepo(tr)
	userRepo := repository.NewRepo(tr)

	uid, err := userRepo.AddUser(ctx, &entityUser.User{
		ID:           uuid.New(),
		Name:         "test_user",
		Email:        "test_user@email.com",
		Role:         runtime.User,
		PasswordHash: "qwerty",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err = userRepo.RemoveUser(ctx, uid)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("empty vocabularies", func(t *testing.T) {
		vocabs, err := repo.GetVocabsWithCountWords(ctx, uid, []uint8{uint8(access.Public), uint8(access.Subscribers)})
		if err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, 0, len(vocabs))
	})
	t.Run("get user vocabs", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
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
			vocabs, err := repo.GetVocabsWithCountWords(ctx, uid, []uint8{uint8(access.Public), uint8(access.Subscribers)})
			if err != nil {
				return err
			}

			assert.Equal(t, 10, len(vocabs))

			return errCancelTx
		})

		assert.ErrorIs(t, err, errCancelTx)
	})
}

func TestGetVocabulariesRecommended(t *testing.T) {
	ctx := context.Background()

	tp := postgres.NewTempPostgres(ctx, "../../../../..")
	defer tp.DropDB(ctx)

	if tp == nil {
		t.Fatal("can't init container for DB")
	}

	tr := transactor.NewTransactor(tp.PgxPool)
	repo := NewRepo(tr)
	userRepo := repository.NewRepo(tr)

	uid, err := userRepo.AddUser(ctx, &entityUser.User{
		ID:           uuid.New(),
		Name:         "test_user",
		Email:        "test_user@email.com",
		Role:         runtime.User,
		PasswordHash: "qwerty",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err = userRepo.RemoveUser(ctx, uid)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("empty vocabularies", func(t *testing.T) {
		recommendedVocabs, err := repo.GetVocabulariesRecommended(ctx, uid, []uint8{uint8(access.Public), uint8(access.Subscribers)}, 3)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 0, len(recommendedVocabs))
	})
	t.Run("get recommended vocabularies when user don't have vocabs", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			for k := 0; k < 5; k++ {
				nuid, err := userRepo.AddUser(ctx, &entityUser.User{
					ID:           uuid.New(),
					Name:         fmt.Sprintf("test_user_%d", k),
					Email:        fmt.Sprintf("test_user_%d@email.com", k),
					Role:         runtime.User,
					PasswordHash: "qwerty",
				})
				if err != nil {
					t.Fatal(err)
				}
				for i := 0; i < 10; i++ {
					_, err := repo.AddVocab(ctx, entity.Vocab{
						ID:            uuid.New(),
						UserID:        nuid,
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
			}
			vocabs, err := repo.GetVocabulariesRecommended(ctx, uid, []uint8{uint8(access.Public), uint8(access.Subscribers)}, 3)
			if err != nil {
				return err
			}

			assert.Equal(t, 0, len(vocabs))

			return errCancelTx
		})

		assert.ErrorIs(t, err, errCancelTx)
	})
	t.Run("get recommended vocabularies when user have vocabs with different language", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			_, err := repo.AddVocab(ctx, entity.Vocab{
				ID:            uuid.New(),
				UserID:        uid,
				Name:          "test",
				Access:        uint8(access.Subscribers),
				NativeLang:    "fi",
				TranslateLang: "ru",
				Description:   "",
				Tags:          nil,
				CreatedAt:     time.Now().UTC(),
				UpdatedAt:     time.Now().UTC(),
			}, nil)
			if err != nil {
				return err
			}
			for k := 0; k < 5; k++ {
				nuid, err := userRepo.AddUser(ctx, &entityUser.User{
					ID:           uuid.New(),
					Name:         fmt.Sprintf("test_user_%d", k),
					Email:        fmt.Sprintf("test_user_%d@email.com", k),
					Role:         runtime.User,
					PasswordHash: "qwerty",
				})
				if err != nil {
					t.Fatal(err)
				}
				for i := 0; i < 10; i++ {
					_, err := repo.AddVocab(ctx, entity.Vocab{
						ID:            uuid.New(),
						UserID:        nuid,
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
			}
			vocabs, err := repo.GetVocabulariesRecommended(ctx, uid, []uint8{uint8(access.Public), uint8(access.Subscribers)}, 3)
			if err != nil {
				return err
			}

			assert.Equal(t, 0, len(vocabs))

			return errCancelTx
		})

		assert.ErrorIs(t, err, errCancelTx)
	})
	t.Run("get recommended vocabularies", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			_, err := repo.AddVocab(ctx, entity.Vocab{
				ID:            uuid.New(),
				UserID:        uid,
				Name:          "test",
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
			for k := 0; k < 10; k++ {
				nuid, err := userRepo.AddUser(ctx, &entityUser.User{
					ID:           uuid.New(),
					Name:         fmt.Sprintf("test_user_%d", k),
					Email:        fmt.Sprintf("test_user_%d@email.com", k),
					Role:         runtime.User,
					PasswordHash: "qwerty",
				})
				if err != nil {
					return err
				}
				vid, err := repo.AddVocab(ctx, entity.Vocab{
					ID:            uuid.New(),
					UserID:        nuid,
					Name:          fmt.Sprintf("test_%d", k),
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

				for j := 0; j < 10-k; j++ {
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
			vocabs, err := repo.GetVocabulariesRecommended(ctx, uid, []uint8{uint8(access.Public), uint8(access.Subscribers)}, 3)
			if err != nil {
				return err
			}

			assert.Equal(t, 3, len(vocabs))
			assert.Equal(t, uint(10), vocabs[0].WordsCount)
			assert.Equal(t, uint(9), vocabs[1].WordsCount)
			assert.Equal(t, uint(8), vocabs[2].WordsCount)

			return errCancelTx
		})

		assert.ErrorIs(t, err, errCancelTx)
	})
}
