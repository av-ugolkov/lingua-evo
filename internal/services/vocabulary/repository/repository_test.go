package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/db/postgres"
	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"
	userRepository "github.com/av-ugolkov/lingua-evo/internal/services/user/repository"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	tr        *transactor.Transactor
	vocabRepo *VocabRepo
	userRepo  *userRepository.UserRepo
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	tp := postgres.NewTempPostgres(ctx, "../../../..")
	defer tp.DropDB(ctx)

	tr = transactor.NewTransactor(tp.PgxPool)
	vocabRepo = NewRepo(tr)
	userRepo = userRepository.NewRepo(tr)

	code := m.Run()
	os.Exit(code)
}

func TestGetVocabulariesWithMaxWords(t *testing.T) {
	ctx := context.Background()

	t.Run("empty vocabularies", func(t *testing.T) {
		vocabs, err := vocabRepo.GetVocabulariesWithMaxWords(ctx, []uint8{uint8(access.Public), uint8(access.Subscribers)}, 3)
		if err != nil {
			assert.Error(t, err)
			return
		}
		assert.Equal(t, 0, len(vocabs))
	})
	t.Run("vocabularies are empty", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			uid, err := userRepo.AddUser(ctx, &entityUser.User{
				ID:       uuid.New(),
				Nickname: "test_user",
				Email:    "test_user@email.com",
				Role:     runtime.User,
			}, "qwerty")
			if err != nil {
				return err
			}

			for i := 0; i < 10; i++ {
				_, err := vocabRepo.AddVocab(ctx, entity.Vocab{
					ID:            uuid.New(),
					UserID:        uid,
					Name:          fmt.Sprintf("test_%d", i),
					Access:        uint8(access.Subscribers),
					NativeLang:    "en",
					TranslateLang: "ru",
					Description:   runtime.EmptyString,
					CreatedAt:     time.Now().UTC(),
					UpdatedAt:     time.Now().UTC(),
				})
				if err != nil {
					return err
				}
			}

			vocabs, err := vocabRepo.GetVocabulariesWithMaxWords(ctx, []uint8{uint8(access.Public), uint8(access.Subscribers)}, 3)
			if err != nil {
				return err
			}

			assert.Equal(t, 0, len(vocabs))

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
	t.Run("get vocabularies with max count words", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			uid, err := userRepo.AddUser(ctx, &entityUser.User{
				ID:       uuid.New(),
				Nickname: "test_user",
				Email:    "test_user@email.com",
				Role:     runtime.User,
			}, "qwerty")
			if err != nil {
				return err
			}

			for i := 0; i < 10; i++ {
				vid, err := vocabRepo.AddVocab(ctx, entity.Vocab{
					ID:            uuid.New(),
					UserID:        uid,
					Name:          fmt.Sprintf("test_%d", i),
					Access:        uint8(access.Subscribers),
					NativeLang:    "en",
					TranslateLang: "ru",
					Description:   runtime.EmptyString,
					CreatedAt:     time.Now().UTC(),
					UpdatedAt:     time.Now().UTC(),
				})
				if err != nil {
					return err
				}

				for j := 0; j < 10-i; j++ {
					_, err := vocabRepo.AddWord(ctx, entity.VocabWord{
						VocabID:       vid,
						ID:            uuid.New(),
						NativeID:      uuid.New(),
						TranslateIDs:  []uuid.UUID{},
						Pronunciation: runtime.EmptyString,
					})
					if err != nil {
						return err
					}
				}
			}

			vocabs, err := vocabRepo.GetVocabulariesWithMaxWords(ctx, []uint8{uint8(access.Public), uint8(access.Subscribers)}, 3)
			if err != nil {
				return err
			}

			assert.Equal(t, 3, len(vocabs))
			assert.Equal(t, uint(10), vocabs[0].WordsCount)
			assert.Equal(t, uint(9), vocabs[1].WordsCount)
			assert.Equal(t, uint(8), vocabs[2].WordsCount)

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
}

func TestGetVocabsWithCountWords(t *testing.T) {
	ctx := context.Background()

	owner, err := userRepo.AddUser(ctx, &entityUser.User{
		ID:       uuid.New(),
		Nickname: "test_user",
		Email:    "test_user@email.com",
		Role:     runtime.User,
	}, "qwerty")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err = userRepo.RemoveUser(ctx, owner)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("empty vocabularies", func(t *testing.T) {
		vocabs, err := vocabRepo.GetVocabsWithCountWords(ctx, uuid.Nil, owner, []uint8{uint8(access.Public), uint8(access.Subscribers)})
		if err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, 0, len(vocabs))
	})
	t.Run("get user vocabs", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			for i := 0; i < 10; i++ {
				_, err := vocabRepo.AddVocab(ctx, entity.Vocab{
					ID:            uuid.New(),
					UserID:        owner,
					Name:          fmt.Sprintf("test_%d", i),
					Access:        uint8(access.Subscribers),
					NativeLang:    "en",
					TranslateLang: "ru",
					Description:   runtime.EmptyString,
					CreatedAt:     time.Now().UTC(),
					UpdatedAt:     time.Now().UTC(),
				})
				if err != nil {
					return err
				}
			}
			vocabs, err := vocabRepo.GetVocabsWithCountWords(ctx, uuid.Nil, owner, []uint8{uint8(access.Public), uint8(access.Subscribers)})
			if err != nil {
				return err
			}

			assert.Equal(t, 10, len(vocabs))

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
}

func TestGetVocabulariesRecommended(t *testing.T) {
	ctx := context.Background()

	uid, err := userRepo.AddUser(ctx, &entityUser.User{
		ID:       uuid.New(),
		Nickname: "test_user",
		Email:    "test_user@email.com",
		Role:     runtime.User,
	}, "qwerty")
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
		recommendedVocabs, err := vocabRepo.GetVocabulariesRecommended(ctx, uid, []uint8{uint8(access.Public), uint8(access.Subscribers)}, 3)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 0, len(recommendedVocabs))
	})
	t.Run("get recommended vocabularies when user don't have vocabs", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			for k := 0; k < 5; k++ {
				nuid, err := userRepo.AddUser(ctx, &entityUser.User{
					ID:       uuid.New(),
					Nickname: fmt.Sprintf("test_user_%d", k),
					Email:    fmt.Sprintf("test_user_%d@email.com", k),
					Role:     runtime.User,
				}, "qwerty")
				if err != nil {
					t.Fatal(err)
				}
				for i := 0; i < 10; i++ {
					_, err := vocabRepo.AddVocab(ctx, entity.Vocab{
						ID:            uuid.New(),
						UserID:        nuid,
						Name:          fmt.Sprintf("test_%d", i),
						Access:        uint8(access.Subscribers),
						NativeLang:    "en",
						TranslateLang: "ru",
						Description:   runtime.EmptyString,
						CreatedAt:     time.Now().UTC(),
						UpdatedAt:     time.Now().UTC(),
					})
					if err != nil {
						return err
					}
				}
			}
			vocabs, err := vocabRepo.GetVocabulariesRecommended(ctx, uid, []uint8{uint8(access.Public), uint8(access.Subscribers)}, 3)
			if err != nil {
				return err
			}

			assert.Equal(t, 0, len(vocabs))

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
	t.Run("get recommended vocabularies when user have vocabs with different language", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			_, err := vocabRepo.AddVocab(ctx, entity.Vocab{
				ID:            uuid.New(),
				UserID:        uid,
				Name:          "test",
				Access:        uint8(access.Subscribers),
				NativeLang:    "fi",
				TranslateLang: "ru",
				Description:   runtime.EmptyString,
				CreatedAt:     time.Now().UTC(),
				UpdatedAt:     time.Now().UTC(),
			})
			if err != nil {
				return err
			}
			for k := 0; k < 5; k++ {
				nuid, err := userRepo.AddUser(ctx, &entityUser.User{
					ID:       uuid.New(),
					Nickname: fmt.Sprintf("test_user_%d", k),
					Email:    fmt.Sprintf("test_user_%d@email.com", k),
					Role:     runtime.User,
				}, "qwerty")
				if err != nil {
					t.Fatal(err)
				}
				for i := 0; i < 10; i++ {
					_, err := vocabRepo.AddVocab(ctx, entity.Vocab{
						ID:            uuid.New(),
						UserID:        nuid,
						Name:          fmt.Sprintf("test_%d", i),
						Access:        uint8(access.Subscribers),
						NativeLang:    "en",
						TranslateLang: "ru",
						Description:   runtime.EmptyString,
						CreatedAt:     time.Now().UTC(),
						UpdatedAt:     time.Now().UTC(),
					})
					if err != nil {
						return err
					}
				}
			}
			vocabs, err := vocabRepo.GetVocabulariesRecommended(ctx, uid, []uint8{uint8(access.Public), uint8(access.Subscribers)}, 3)
			if err != nil {
				return err
			}

			assert.Equal(t, 0, len(vocabs))

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
	t.Run("get recommended vocabularies", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			_, err := vocabRepo.AddVocab(ctx, entity.Vocab{
				ID:            uuid.New(),
				UserID:        uid,
				Name:          "test",
				Access:        uint8(access.Subscribers),
				NativeLang:    "en",
				TranslateLang: "ru",
				Description:   runtime.EmptyString,
				CreatedAt:     time.Now().UTC(),
				UpdatedAt:     time.Now().UTC(),
			})
			if err != nil {
				return err
			}
			for k := 0; k < 10; k++ {
				nuid, err := userRepo.AddUser(ctx, &entityUser.User{
					ID:       uuid.New(),
					Nickname: fmt.Sprintf("test_user_%d", k),
					Email:    fmt.Sprintf("test_user_%d@email.com", k),
					Role:     runtime.User,
				}, "qwerty")
				if err != nil {
					return err
				}
				vid, err := vocabRepo.AddVocab(ctx, entity.Vocab{
					ID:            uuid.New(),
					UserID:        nuid,
					Name:          fmt.Sprintf("test_%d", k),
					Access:        uint8(access.Subscribers),
					NativeLang:    "en",
					TranslateLang: "ru",
					Description:   runtime.EmptyString,
					CreatedAt:     time.Now().UTC(),
					UpdatedAt:     time.Now().UTC(),
				})
				if err != nil {
					return err
				}

				for j := 0; j < 10-k; j++ {
					_, err := vocabRepo.AddWord(ctx, entity.VocabWord{
						VocabID:       vid,
						ID:            uuid.New(),
						NativeID:      uuid.New(),
						TranslateIDs:  []uuid.UUID{},
						Pronunciation: runtime.EmptyString,
					})
					if err != nil {
						return err
					}
				}
			}
			vocabs, err := vocabRepo.GetVocabulariesRecommended(ctx, uid, []uint8{uint8(access.Public), uint8(access.Subscribers)}, 3)
			if err != nil {
				return err
			}

			assert.Equal(t, 3, len(vocabs))
			assert.Equal(t, uint(10), vocabs[0].WordsCount)
			assert.Equal(t, uint(9), vocabs[1].WordsCount)
			assert.Equal(t, uint(8), vocabs[2].WordsCount)

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
}
