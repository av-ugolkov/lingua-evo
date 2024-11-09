package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/av-ugolkov/lingua-evo/internal/db/postgres"
	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entityDict "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	dictRepo "github.com/av-ugolkov/lingua-evo/internal/services/dictionary/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/example"
	exampleRepo "github.com/av-ugolkov/lingua-evo/internal/services/example/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/language"
	"github.com/av-ugolkov/lingua-evo/internal/services/language/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/tag"
	tagRepo "github.com/av-ugolkov/lingua-evo/internal/services/tag/repository"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"
	userRepo "github.com/av-ugolkov/lingua-evo/internal/services/user/repository"
	user "github.com/av-ugolkov/lingua-evo/internal/services/user/service"
	vocabEntity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	vocabRepo "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/repository"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetVocabularies(t *testing.T) {
	ctx := context.Background()

	tp := postgres.NewTempPostgres(ctx, "../../../..")
	defer tp.DropDB(ctx)

	if tp == nil {
		t.Fatal("can't init container for DB")
	}

	tr := transactor.NewTransactor(tp.PgxPool)
	userSvc := user.NewService(tr, userRepo.NewRepo(tr), nil, nil)

	usr, err := userSvc.GetUserByName(ctx, "admin")
	if err != nil {
		t.Fatal(err)
	}

	tagSvc := tag.NewService(tagRepo.NewRepo(tr))
	langSvc := language.NewService(repository.NewRepo(tr))
	dictSvc := entityDict.NewService(dictRepo.NewRepo(tr), langSvc)
	exampleSvc := example.NewService(exampleRepo.NewRepo(tr))
	mockEventsSvc := new(mockEventsSvc)
	mockEventsSvc.On("AsyncAddEvent", mock.Anything, mock.Anything).Return(nil)
	vocabSvc := NewService(tr, vocabRepo.NewRepo(tr), userSvc, exampleSvc, dictSvc, tagSvc, nil, mockEventsSvc)

	t.Run("empty vocab", func(t *testing.T) {
		var (
			page          = 1
			itemsPerPage  = 5
			typeSort      = 1
			order         = 0
			search        = runtime.EmptyString
			nativeLang    = "en"
			translateLang = "ru"
			maxWords      = 5
		)
		vocabs, count, err := vocabSvc.GetVocabularies(ctx, usr.ID, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang, maxWords)
		if err != nil {
			assert.Error(t, err)
			return
		}
		if !assert.Equal(t, 0, count) {
			t.Errorf("GetVocabularies() got = %v, want %v", count, 0)
		}
		if !assert.Equal(t, 0, len(vocabs)) {
			t.Errorf("GetVocabularies() got = %v, want %v", len(vocabs), 5)
		}
	})

	t.Run("add 4 vocabs", func(t *testing.T) {
		var (
			expectCount  = 4
			expectVocabs = 4
		)
		var (
			page          = 1
			itemsPerPage  = 5
			typeSort      = 1
			order         = 0
			search        = runtime.EmptyString
			nativeLang    = "en"
			translateLang = "ru"
			maxWords      = 5
		)

		err = tr.CreateTransaction(ctx, func(ctx context.Context) error {
			for i := 0; i < 2; i++ {
				uid, err := userSvc.AddUser(ctx, entityUser.UserCreate{
					ID:       uuid.New(),
					Name:     fmt.Sprintf("user_%d", i),
					Password: fmt.Sprintf("password_%d", i),
					Email:    fmt.Sprintf("user_%d@user_%d.com", i, i),
					Role:     runtime.User,
					Code:     0,
				})
				if err != nil {
					return err
				}

				_, err = addVocabs(ctx, vocabSvc, uid, 2)
				if err != nil {
					return err
				}
			}

			vocabs, count, err := vocabSvc.GetVocabularies(ctx, usr.ID, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang, maxWords)
			if err != nil {
				return err
			}
			if !assert.Equal(t, expectCount, count) {
				t.Errorf("GetVocabularies() got = %v, want %v", count, 0)
			}
			if !assert.Equal(t, expectVocabs, len(vocabs)) {
				t.Errorf("GetVocabularies() got = %v, want %v", len(vocabs), 5)
			}

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})

	t.Run("add 9 vocabs", func(t *testing.T) {
		var (
			expectCount  = 9
			expectVocabs = 5
		)
		var (
			page          = 1
			itemsPerPage  = 5
			typeSort      = 1
			order         = 0
			search        = runtime.EmptyString
			nativeLang    = "en"
			translateLang = "ru"
			maxWords      = 5
		)

		err = tr.CreateTransaction(ctx, func(ctx context.Context) error {
			for i := 0; i < 3; i++ {
				uid, err := userSvc.AddUser(ctx, entityUser.UserCreate{
					ID:       uuid.New(),
					Name:     fmt.Sprintf("user_%d", i),
					Password: fmt.Sprintf("password_%d", i),
					Email:    fmt.Sprintf("user_%d@user_%d.com", i, i),
					Role:     runtime.User,
					Code:     0,
				})
				if err != nil {
					return err
				}

				_, err = addVocabs(ctx, vocabSvc, uid, 3)
				if err != nil {
					return err
				}
			}

			vocabs, count, err := vocabSvc.GetVocabularies(ctx, usr.ID, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang, maxWords)
			if err != nil {
				return err
			}
			if !assert.Equal(t, expectCount, count) {
				t.Errorf("GetVocabularies() got = %v, want %v", count, 0)
			}
			if !assert.Equal(t, expectVocabs, len(vocabs)) {
				t.Errorf("GetVocabularies() got = %v, want %v", len(vocabs), 5)
			}

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})

	t.Run("add vocabs with words", func(t *testing.T) {
		var (
			expectCount  = 9
			expectVocabs = 5
			expectWords  = 3
		)
		var (
			page          = 1
			itemsPerPage  = 5
			typeSort      = 1
			order         = 0
			search        = runtime.EmptyString
			nativeLang    = "en"
			translateLang = "ru"
			maxWords      = 5
		)

		err = tr.CreateTransaction(ctx, func(ctx context.Context) error {
			for i := 0; i < 3; i++ {
				uid, err := userSvc.AddUser(ctx, entityUser.UserCreate{
					ID:       uuid.New(),
					Name:     fmt.Sprintf("user_%d", i),
					Password: fmt.Sprintf("password_%d", i),
					Email:    fmt.Sprintf("user_%d@user_%d.com", i, i),
					Role:     runtime.User,
					Code:     0,
				})
				if err != nil {
					return err
				}

				vocabs, err := addVocabs(ctx, vocabSvc, uid, 3)
				if err != nil {
					return err
				}

				for _, vocab := range vocabs {
					for i := 0; i < 3; i++ {
						_, err := vocabSvc.AddWord(ctx, uid, vocabEntity.VocabWordData{
							VocabID: vocab.ID,
							Native: entityDict.DictWord{
								Text: fmt.Sprintf("text_%d_%s", i, vocab.Name),
							},
						})
						if err != nil {
							return err
						}
					}
				}
			}

			vocabs, count, err := vocabSvc.GetVocabularies(ctx, usr.ID, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang, maxWords)
			if err != nil {
				return err
			}
			if !assert.Equal(t, expectCount, count) {
				t.Errorf("GetVocabularies() got = %v, want %v", count, expectCount)
			}
			if !assert.Equal(t, expectVocabs, len(vocabs)) {
				t.Errorf("GetVocabularies() got = %v, want %v", len(vocabs), expectVocabs)
			}
			for _, vocab := range vocabs {
				if !assert.Equal(t, uint(expectWords), vocab.WordsCount) {
					t.Errorf("GetVocabularies() got = %v, want %v", vocab.WordsCount, expectWords)
				}

				if !assert.Equal(t, expectWords, len(vocab.Words)) {
					t.Errorf("GetVocabularies() got = %v, want %v", len(vocab.Words), expectWords)
				}
			}

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})

	t.Run("add vocabs with words 2", func(t *testing.T) {
		var (
			expectCount         = 9
			expectVocabs        = 5
			expectWords    uint = 10
			expectGetWords      = 5
		)
		var (
			page          = 1
			itemsPerPage  = 5
			typeSort      = 1
			order         = 0
			search        = runtime.EmptyString
			nativeLang    = "en"
			translateLang = "ru"
			maxWords      = 5
		)

		err = tr.CreateTransaction(ctx, func(ctx context.Context) error {
			for i := 0; i < 3; i++ {
				uid, err := userSvc.AddUser(ctx, entityUser.UserCreate{
					ID:       uuid.New(),
					Name:     fmt.Sprintf("user_%d", i),
					Password: fmt.Sprintf("password_%d", i),
					Email:    fmt.Sprintf("user_%d@user_%d.com", i, i),
					Role:     runtime.User,
					Code:     0,
				})
				if err != nil {
					return err
				}

				vocabs, err := addVocabs(ctx, vocabSvc, uid, 3)
				if err != nil {
					return err
				}

				for _, vocab := range vocabs {
					for i := 0; i < 10; i++ {
						_, err := vocabSvc.AddWord(ctx, uid, vocabEntity.VocabWordData{
							VocabID: vocab.ID,
							Native: entityDict.DictWord{
								Text: fmt.Sprintf("text_%d_%s", i, vocab.Name),
							},
						})
						if err != nil {
							return err
						}
					}
				}
			}

			vocabs, count, err := vocabSvc.GetVocabularies(ctx, usr.ID, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang, maxWords)
			if err != nil {
				return err
			}
			if !assert.Equal(t, expectCount, count) {
				t.Errorf("GetVocabularies() got = %v, want %v", count, expectCount)
			}
			if !assert.Equal(t, expectVocabs, len(vocabs)) {
				t.Errorf("GetVocabularies() got = %v, want %v", len(vocabs), expectVocabs)
			}
			for _, vocab := range vocabs {
				if !assert.Equal(t, expectWords, vocab.WordsCount) {
					t.Errorf("GetVocabularies() got = %v, want %v", vocab.WordsCount, expectWords)
				}

				if !assert.Equal(t, expectGetWords, len(vocab.Words)) {
					t.Errorf("GetVocabularies() got = %v, want %v", len(vocab.Words), expectGetWords)
				}
			}

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
}
