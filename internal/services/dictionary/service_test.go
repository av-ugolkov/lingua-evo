package dictionary

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/av-ugolkov/lingua-evo/internal/services/dictionary/model"

	"github.com/stretchr/testify/assert"
)

func TestService_AddWord(t *testing.T) {
	t.Run("AddWord", func(t *testing.T) {
		var (
			ctx    = context.Background()
			wordID = uuid.MustParse("efd91d3f-fd7e-4da3-9860-4e3a0012c887")
			data   = model.WordRq{
				Text:          "word",
				Pronunciation: "[wɜ:d]",
				LangCode:      "en",
			}
		)
		repoWordMock := new(mockRepoDictionary)
		repoWordMock.On("AddWords", ctx, mock.Anything).Return([]uuid.UUID{wordID}, nil)
		langSvcMock := new(mockLangSvc)
		langSvcMock.On("CheckLanguage", ctx, data.LangCode).Return(nil)

		s := &Service{repo: repoWordMock, langSvc: langSvcMock}

		got, err := s.AddWord(ctx, data)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}
