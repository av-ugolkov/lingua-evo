package dictionary

import (
	"context"
	"testing"

	"github.com/av-ugolkov/lingua-evo/internal/services/dictionary/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_AddWord(t *testing.T) {
	t.Run("AddWord", func(t *testing.T) {
		var (
			ctx  = context.Background()
			data = model.WordRq{
				Text:          "word",
				Pronunciation: "[wɜ:d]",
				LangCode:      "en",
			}
		)
		repoWordMock := new(mockRepoDictionary)
		repoWordMock.On("AddWords", ctx, mock.Anything).Return(mock.Anything, nil)
		s := &Service{repo: repoWordMock}

		got, err := s.AddWord(ctx, data)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}
