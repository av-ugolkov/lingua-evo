package dictionary

import (
	"context"
	entityLanguage "github.com/av-ugolkov/lingua-evo/internal/services/language"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_AddWord(t *testing.T) {
	t.Run("AddWord", func(t *testing.T) {
		var (
			ctx  = context.Background()
			data = []Word{
				{
					ID:            uuid.New(),
					Text:          "word",
					Pronunciation: "[wɜ:d]",
					LangCode:      "en",
				},
			}
		)
		repoWordMock := new(mockRepoDictionary)
		repoWordMock.On("AddWords", ctx, mock.Anything).Return([]uuid.UUID{data[0].ID}, nil)
		langSvcMock := new(mockLangSvc)
		langSvcMock.On("GetAvailableLanguages", ctx).Return([]*entityLanguage.Language{{Code: data[0].LangCode}}, nil)

		s := &Service{repo: repoWordMock, langSvc: langSvcMock}

		got, err := s.AddWords(ctx, data)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}
