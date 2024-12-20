package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/db/postgres"
	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	"github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	dictRepo "github.com/av-ugolkov/lingua-evo/internal/services/dictionary/repository"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/events"
	eventsRepo "github.com/av-ugolkov/lingua-evo/internal/services/events/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/example"
	exampleRepo "github.com/av-ugolkov/lingua-evo/internal/services/example/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/language"
	langRepo "github.com/av-ugolkov/lingua-evo/internal/services/language/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/notifications"
	repoNotif "github.com/av-ugolkov/lingua-evo/internal/services/notifications/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/tag"
	tagRepo "github.com/av-ugolkov/lingua-evo/internal/services/tag/repository"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"
	userRepo "github.com/av-ugolkov/lingua-evo/internal/services/user/repository"
	user "github.com/av-ugolkov/lingua-evo/internal/services/user/service"
	entityVocab "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	vocabRepo "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/repository"
	vocabService "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/service"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	tr            *transactor.Transactor
	userSvc       *user.Service
	vocabularySvc *vocabService.Service
	langSvc       *language.Service
	dictSvc       *dictionary.Service
	notifSvc      *notifications.Service
	tagSvc        *tag.Service
	exampleSvc    *example.Service
	eventsSvc     *Service

	usr *entityUser.User
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	tp := postgres.NewTempPostgres(ctx, "../../../..")

	tr = transactor.NewTransactor(tp.PgxPool)

	langSvc = language.NewService(langRepo.NewRepo(tr))
	dictSvc = dictionary.NewService(dictRepo.NewRepo(tr), langSvc)
	notifSvc = notifications.NewService(repoNotif.NewRepo(tr))
	tagSvc = tag.NewService(tagRepo.NewRepo(tr))
	eventsSvc = NewService(tr, eventsRepo.NewRepo(tr), notifSvc)
	userSvc = user.NewService(tr, userRepo.NewRepo(tr), nil, nil)
	exampleSvc = example.NewService(exampleRepo.NewRepo(tr))
	vocabularySvc = vocabService.NewService(tr, vocabRepo.NewRepo(tr), exampleSvc, dictSvc, tagSvc, nil, eventsSvc)

	var err error
	usr, err = userSvc.GetUserByNickname(ctx, "admin")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	code := m.Run()

	tp.DropDB(ctx)
	os.Exit(code)
}

func TestGetCountEvents(t *testing.T) {
	ctx := context.Background()

	t.Run("zero events", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			tempUID, err := userSvc.AddUser(ctx, entityUser.User{
				ID:       uuid.New(),
				Nickname: "user_temp",
				Email:    "user_temp@user_temp.com",
				Role:     runtime.User,
			}, "password_temp")
			if err != nil {
				return err
			}
			vocabTemp, err := vocabularySvc.UserAddVocabulary(ctx, entityVocab.Vocab{
				UserID:        usr.ID,
				Name:          "vocab_temp",
				NativeLang:    "en",
				TranslateLang: "ru",
				Access:        1,
			})
			if err != nil {
				return err
			}

			isNotif, err := notifSvc.SetVocabNotification(ctx, tempUID, vocabTemp.ID)
			if err != nil {
				return err
			}
			if !isNotif {
				return fmt.Errorf("can't set notification")
			}

			count, err := eventsSvc.GetCountEvents(ctx, tempUID)
			if err != nil {
				return err
			}

			assert.Equal(t, 0, count)

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
	t.Run("one events", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			tempUID, err := userSvc.AddUser(ctx, entityUser.User{
				ID:       uuid.New(),
				Nickname: "user_temp",
				Email:    "user_temp@user_temp.com",
				Role:     runtime.User,
			}, "password_temp")
			if err != nil {
				return err
			}
			vocabTemp, err := vocabularySvc.UserAddVocabulary(ctx, entityVocab.Vocab{
				UserID:        usr.ID,
				Name:          "vocab_temp",
				NativeLang:    "en",
				TranslateLang: "ru",
				Access:        1,
			})
			if err != nil {
				return err
			}

			isNotif, err := notifSvc.SetVocabNotification(ctx, tempUID, vocabTemp.ID)
			if err != nil {
				return err
			}
			if !isNotif {
				return fmt.Errorf("can't set notification")
			}

			_, err = vocabularySvc.AddWord(ctx, usr.ID, entityVocab.VocabWordData{
				VocabID: vocabTemp.ID,
				Native: entityVocab.DictWord{
					Text: "temp",
				},
			})
			if err != nil {
				return err
			}

			time.Sleep(time.Millisecond * 300)

			count, err := eventsSvc.GetCountEvents(ctx, tempUID)
			if err != nil {
				return err
			}

			assert.Equal(t, 1, count)

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
	t.Run("one events after subscribe", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			tempUID, err := userSvc.AddUser(ctx, entityUser.User{
				ID:       uuid.New(),
				Nickname: "user_temp",
				Email:    "user_temp@user_temp.com",
				Role:     runtime.User,
			}, "password_temp")
			if err != nil {
				return err
			}
			vocabTemp, err := vocabularySvc.UserAddVocabulary(ctx, entityVocab.Vocab{
				UserID:        usr.ID,
				Name:          "vocab_temp",
				NativeLang:    "en",
				TranslateLang: "ru",
				Access:        1,
			})
			if err != nil {
				return err
			}

			_, err = vocabularySvc.AddWord(ctx, usr.ID, entityVocab.VocabWordData{
				VocabID: vocabTemp.ID,
				Native: entityVocab.DictWord{
					Text: "temp",
				},
			})
			if err != nil {
				return err
			}

			time.Sleep(time.Millisecond * 300)

			isNotif, err := notifSvc.SetVocabNotification(ctx, tempUID, vocabTemp.ID)
			if err != nil {
				return err
			}
			if !isNotif {
				return fmt.Errorf("can't set notification")
			}

			count, err := eventsSvc.GetCountEvents(ctx, tempUID)
			if err != nil {
				return err
			}

			assert.Equal(t, 0, count)

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
}

func TestGetEvents(t *testing.T) {
	ctx := context.Background()

	t.Run("zero events", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			tempUID, err := userSvc.AddUser(ctx, entityUser.User{
				ID:       uuid.New(),
				Nickname: "user_temp",
				Email:    "user_temp@user_temp.com",
				Role:     runtime.User,
			}, "password_temp")
			if err != nil {
				return err
			}
			vocabTemp, err := vocabularySvc.UserAddVocabulary(ctx, entityVocab.Vocab{
				UserID:        usr.ID,
				Name:          "vocab_temp",
				NativeLang:    "en",
				TranslateLang: "ru",
				Access:        1,
			})
			if err != nil {
				return err
			}

			isNotif, err := notifSvc.SetVocabNotification(ctx, tempUID, vocabTemp.ID)
			if err != nil {
				return err
			}
			if !isNotif {
				return fmt.Errorf("can't set notification")
			}

			events, err := eventsSvc.GetEvents(ctx, tempUID)
			if err != nil {
				return err
			}

			assert.Equal(t, 0, len(events))

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
	t.Run("one events", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			tempUID, err := userSvc.AddUser(ctx, entityUser.User{
				ID:       uuid.New(),
				Nickname: "user_temp",
				Email:    "user_temp@user_temp.com",
				Role:     runtime.User,
			}, "password_temp")
			if err != nil {
				return err
			}
			vocabTemp, err := vocabularySvc.UserAddVocabulary(ctx, entityVocab.Vocab{
				UserID:        usr.ID,
				Name:          "vocab_temp",
				NativeLang:    "en",
				TranslateLang: "ru",
				Access:        1,
			})
			if err != nil {
				return err
			}

			isNotif, err := notifSvc.SetVocabNotification(ctx, tempUID, vocabTemp.ID)
			if err != nil {
				return err
			}
			if !isNotif {
				return fmt.Errorf("can't set notification")
			}

			_, err = vocabularySvc.AddWord(ctx, usr.ID, entityVocab.VocabWordData{
				VocabID: vocabTemp.ID,
				Native: entityVocab.DictWord{
					Text: "temp",
				},
			})
			if err != nil {
				return err
			}

			time.Sleep(time.Millisecond * 300)

			events, err := eventsSvc.GetEvents(ctx, tempUID)
			if err != nil {
				return err
			}

			assert.Equal(t, 1, len(events))

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
	t.Run("one events after subscribe", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			tempUID, err := userSvc.AddUser(ctx, entityUser.User{
				ID:       uuid.New(),
				Nickname: "user_temp",
				Email:    "user_temp@user_temp.com",
				Role:     runtime.User,
			}, "password_temp")
			if err != nil {
				return err
			}
			vocabTemp, err := vocabularySvc.UserAddVocabulary(ctx, entityVocab.Vocab{
				UserID:        usr.ID,
				Name:          "vocab_temp",
				NativeLang:    "en",
				TranslateLang: "ru",
				Access:        1,
			})
			if err != nil {
				return err
			}

			_, err = vocabularySvc.AddWord(ctx, usr.ID, entityVocab.VocabWordData{
				VocabID: vocabTemp.ID,
				Native: entityVocab.DictWord{
					Text: "temp",
				},
			})
			if err != nil {
				return err
			}

			time.Sleep(time.Millisecond * 300)

			isNotif, err := notifSvc.SetVocabNotification(ctx, tempUID, vocabTemp.ID)
			if err != nil {
				return err
			}
			if !isNotif {
				return fmt.Errorf("can't set notification")
			}

			events, err := eventsSvc.GetEvents(ctx, tempUID)
			if err != nil {
				return err
			}

			assert.Equal(t, 0, len(events))

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
}

func TestReadEvent(t *testing.T) {
	ctx := context.Background()

	t.Run("read event", func(t *testing.T) {
		err := tr.CreateTransaction(ctx, func(ctx context.Context) error {
			eid, err := eventsSvc.AddEvent(ctx, entity.Event{
				User: entity.UserData{ID: usr.ID},
				Type: entity.VocabWordCreated,
				Payload: entity.PayloadDataVocab{
					VocabID:    new(uuid.UUID),
					VocabTitle: "vocab_temp",
				},
			})
			if err != nil {
				return err
			}

			err = eventsSvc.ReadEvent(ctx, usr.ID, eid)
			if err != nil {
				return err
			}

			return transactor.ErrCancelTx
		})

		assert.ErrorIs(t, err, transactor.ErrCancelTx)
	})
}
