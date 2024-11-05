package service

import (
	"context"
	"fmt"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/events"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/db/postgres"
	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	"github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	dictRepository "github.com/av-ugolkov/lingua-evo/internal/services/dictionary/repository"
	eventsRepository "github.com/av-ugolkov/lingua-evo/internal/services/events/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/example"
	exampleRepo "github.com/av-ugolkov/lingua-evo/internal/services/example/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/language"
	langRepository "github.com/av-ugolkov/lingua-evo/internal/services/language/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/notifications"
	repoNotif "github.com/av-ugolkov/lingua-evo/internal/services/notifications/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/tag"
	tagRepo "github.com/av-ugolkov/lingua-evo/internal/services/tag/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/user"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"
	userRepository "github.com/av-ugolkov/lingua-evo/internal/services/user/repository"
	entityVocab "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	vocabRepository "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/repository"
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

	langSvc = language.NewService(langRepository.NewRepo(tr))
	dictSvc = dictionary.NewService(dictRepository.NewRepo(tr), langSvc)
	notifSvc = notifications.NewService(repoNotif.NewRepo(tr))
	tagSvc = tag.NewService(tagRepo.NewRepo(tr))
	eventsSvc = NewService(tr, eventsRepository.NewRepo(tr), notifSvc)
	userSvc = user.NewService(userRepository.NewRepo(tr), nil, tr)
	exampleSvc = example.NewService(exampleRepo.NewRepo(tr))
	vocabularySvc = vocabService.NewService(tr, vocabRepository.NewRepo(tr), userSvc, exampleSvc, dictSvc, tagSvc, nil, eventsSvc)

	var err error
	usr, err = userSvc.GetUserByName(ctx, "admin")
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
			tempUID, err := userSvc.AddUser(ctx, entityUser.UserCreate{
				ID:       uuid.New(),
				Name:     "user_temp",
				Password: "password_temp",
				Email:    "user_temp@user_temp.com",
				Role:     runtime.User,
				Code:     0,
			})
			if err != nil {
				return err
			}
			vocabTemp, err := vocabularySvc.UserAddVocabulary(ctx, entityVocab.Vocab{
				UserID:        usr.ID,
				Name:          "vocab_temp",
				NativeLang:    "en",
				TranslateLang: "ru",
				Access:        1,
				Tags:          []tag.Tag{},
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
			tempUID, err := userSvc.AddUser(ctx, entityUser.UserCreate{
				ID:       uuid.New(),
				Name:     "user_temp",
				Password: "password_temp",
				Email:    "user_temp@user_temp.com",
				Role:     runtime.User,
				Code:     0,
			})
			if err != nil {
				return err
			}
			vocabTemp, err := vocabularySvc.UserAddVocabulary(ctx, entityVocab.Vocab{
				UserID:        usr.ID,
				Name:          "vocab_temp",
				NativeLang:    "en",
				TranslateLang: "ru",
				Access:        1,
				Tags:          []tag.Tag{},
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
				Native: dictionary.DictWord{
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
			tempUID, err := userSvc.AddUser(ctx, entityUser.UserCreate{
				ID:       uuid.New(),
				Name:     "user_temp",
				Password: "password_temp",
				Email:    "user_temp@user_temp.com",
				Role:     runtime.User,
				Code:     0,
			})
			if err != nil {
				return err
			}
			vocabTemp, err := vocabularySvc.UserAddVocabulary(ctx, entityVocab.Vocab{
				UserID:        usr.ID,
				Name:          "vocab_temp",
				NativeLang:    "en",
				TranslateLang: "ru",
				Access:        1,
				Tags:          []tag.Tag{},
			})
			if err != nil {
				return err
			}

			_, err = vocabularySvc.AddWord(ctx, usr.ID, entityVocab.VocabWordData{
				VocabID: vocabTemp.ID,
				Native: dictionary.DictWord{
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
			tempUID, err := userSvc.AddUser(ctx, entityUser.UserCreate{
				ID:       uuid.New(),
				Name:     "user_temp",
				Password: "password_temp",
				Email:    "user_temp@user_temp.com",
				Role:     runtime.User,
				Code:     0,
			})
			if err != nil {
				return err
			}
			vocabTemp, err := vocabularySvc.UserAddVocabulary(ctx, entityVocab.Vocab{
				UserID:        usr.ID,
				Name:          "vocab_temp",
				NativeLang:    "en",
				TranslateLang: "ru",
				Access:        1,
				Tags:          []tag.Tag{},
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
			tempUID, err := userSvc.AddUser(ctx, entityUser.UserCreate{
				ID:       uuid.New(),
				Name:     "user_temp",
				Password: "password_temp",
				Email:    "user_temp@user_temp.com",
				Role:     runtime.User,
				Code:     0,
			})
			if err != nil {
				return err
			}
			vocabTemp, err := vocabularySvc.UserAddVocabulary(ctx, entityVocab.Vocab{
				UserID:        usr.ID,
				Name:          "vocab_temp",
				NativeLang:    "en",
				TranslateLang: "ru",
				Access:        1,
				Tags:          []tag.Tag{},
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
				Native: dictionary.DictWord{
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
			tempUID, err := userSvc.AddUser(ctx, entityUser.UserCreate{
				ID:       uuid.New(),
				Name:     "user_temp",
				Password: "password_temp",
				Email:    "user_temp@user_temp.com",
				Role:     runtime.User,
				Code:     0,
			})
			if err != nil {
				return err
			}
			vocabTemp, err := vocabularySvc.UserAddVocabulary(ctx, entityVocab.Vocab{
				UserID:        usr.ID,
				Name:          "vocab_temp",
				NativeLang:    "en",
				TranslateLang: "ru",
				Access:        1,
				Tags:          []tag.Tag{},
			})
			if err != nil {
				return err
			}

			_, err = vocabularySvc.AddWord(ctx, usr.ID, entityVocab.VocabWordData{
				VocabID: vocabTemp.ID,
				Native: dictionary.DictWord{
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
