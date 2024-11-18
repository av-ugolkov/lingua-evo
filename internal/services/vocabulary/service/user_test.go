package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/av-ugolkov/lingua-evo/internal/db/postgres"
	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	"github.com/av-ugolkov/lingua-evo/internal/services/subscribers"
	subscribersRepo "github.com/av-ugolkov/lingua-evo/internal/services/subscribers/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/tag"
	"github.com/av-ugolkov/lingua-evo/internal/services/tag/repository"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"
	userRepo "github.com/av-ugolkov/lingua-evo/internal/services/user/repository"
	user "github.com/av-ugolkov/lingua-evo/internal/services/user/service"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	vocabRepo "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/repository"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	tr             *transactor.Transactor
	userTestSvc    *user.Service
	subscrbTestSvc *subscribers.Service
	vocabSvc       *Service

	usr *entityUser.User
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	tp := postgres.NewTempPostgres(ctx, "../../../..")
	defer tp.DropDB(ctx)

	tr := transactor.NewTransactor(tp.PgxPool)
	userTestSvc = user.NewService(tr, userRepo.NewRepo(tr), nil, nil)
	subscrbTestSvc = subscribers.NewService(subscribersRepo.NewRepo(tr))

	var err error
	usr, err = userTestSvc.GetUserByNickname(ctx, "admin")
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		os.Exit(1)
	}
	tagSvc := tag.NewService(repository.NewRepo(tr))

	vocabSvc = NewService(tr, vocabRepo.NewRepo(tr), nil, nil, tagSvc, subscrbTestSvc, nil)
}

func TestService_UserGetVocabularies(t *testing.T) {
	ctx := context.Background()

	t.Run("empty vocab", func(t *testing.T) {
		var (
			page          = 1
			itemsPerPage  = 5
			typeSort      = 1
			order         = 0
			search        = runtime.EmptyString
			nativeLang    = "en"
			translateLang = "ru"
		)
		vocabs, count, err := vocabSvc.UserGetVocabularies(ctx, usr.ID, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang)
		if err != nil {
			assert.Error(t, err)
			return
		}
		if !assert.Equal(t, 0, count) {
			t.Errorf("UserGetVocabularies() got = %v, want %v", count, 0)
		}
		if !assert.Equal(t, 0, len(vocabs)) {
			t.Errorf("UserGetVocabularies() got = %v, want %v", len(vocabs), 5)
		}
	})

	t.Run("get only user vocabs", func(t *testing.T) {
		var (
			expectCount  = 10
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
		)

		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			_, err := addVocabs(ctx, vocabSvc, usr.ID, expectCount)
			if err != nil {
				return err
			}

			sortVocab, count, err := vocabSvc.UserGetVocabularies(ctx, usr.ID, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang)
			if err != nil {
				return err
			}
			if !assert.Equal(t, expectCount, count) {
				t.Errorf("UserGetVocabularies() got = %v, want %v", count, expectCount)
			}
			if !assert.Equal(t, expectVocabs, len(sortVocab)) {
				t.Errorf("UserGetVocabularies() got = %v, want %v", len(sortVocab), expectVocabs)
			}

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})

	t.Run("get user and subscribers vocabs", func(t *testing.T) {
		var (
			expectCount  = 14
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
		)

		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			_, err := addVocabs(ctx, vocabSvc, usr.ID, 5)
			if err != nil {
				return err
			}

			for i := 0; i < 3; i++ {
				uid, err := userTestSvc.AddUser(ctx, entityUser.User{
					ID:       uuid.New(),
					Nickname: fmt.Sprintf("user_%d", i),
					Email:    fmt.Sprintf("user_%d@user_%d.com", i, i),
					Role:     runtime.User,
				}, fmt.Sprintf("password_%d", i))
				if err != nil {
					return err
				}

				_, err = addVocabs(ctx, vocabSvc, uid, 3)
				if err != nil {
					return err
				}

				err = subscrbTestSvc.Subscribe(ctx, usr.ID, uid)
				if err != nil {
					return err
				}
			}

			sortVocab, count, err := vocabSvc.UserGetVocabularies(ctx, usr.ID, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang)
			if err != nil {
				return err
			}
			if !assert.Equal(t, expectCount, count) {
				t.Errorf("UserGetVocabularies() got = %v, want %v", count, expectCount)
			}
			if !assert.Equal(t, expectVocabs, len(sortVocab)) {
				t.Errorf("UserGetVocabularies() got = %v, want %v", len(sortVocab), expectVocabs)
			}

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
}

func addVocabs(ctx context.Context, vocabSvc *Service, uid uuid.UUID, count int) ([]entity.Vocab, error) {
	vocabs := make([]entity.Vocab, 0, count)
	for j := 0; j < count; j++ {
		vocab, err := vocabSvc.UserAddVocabulary(ctx, entity.Vocab{
			UserID:        uid,
			Name:          fmt.Sprintf("vocab_%d", j),
			NativeLang:    "en",
			TranslateLang: "ru",
			Access:        1,
			Tags:          []tag.Tag{},
		})
		if err != nil {
			return nil, fmt.Errorf("user.createUsers - UserAddVocabulary: %w", err)
		}

		vocabs = append(vocabs, vocab)
	}

	return vocabs, nil
}
