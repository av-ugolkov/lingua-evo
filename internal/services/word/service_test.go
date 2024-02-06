package word

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_AddWord(t *testing.T) {
	t.Run("AddWord", func(t *testing.T) {
		var (
			ctx  = context.Background()
			word = &Word{
				ID:            uuid.New(),
				Text:          "word",
				Pronunciation: "[wɜ:d]",
				LanguageCode:  "en",
			}
		)
		repoWordMock := new(mockRepoWord)
		repoWordMock.On("GetWordByText", ctx, word).Return(word.ID, nil)
		repoWordMock.On("AddWord", ctx, word).Return(word.ID, nil)
		s := &Service{repo: repoWordMock}

		got, err := s.AddWord(ctx, word.Text, word.LanguageCode, word.Pronunciation)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, word.ID, got)
	})

}
