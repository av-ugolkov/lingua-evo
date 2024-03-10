package vocabulary

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	entityExample "github.com/av-ugolkov/lingua-evo/internal/services/example"
	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	entityWord "github.com/av-ugolkov/lingua-evo/internal/services/word"
)

type (
	repoVocabulary interface {
		GetWord(ctx context.Context, dictID, wordID uuid.UUID) (Vocabulary, error)
		AddWord(ctx context.Context, vocabulary Vocabulary) error
		DeleteWord(ctx context.Context, vocabulary Vocabulary) error
		GetRandomVocabulary(ctx context.Context, dictID uuid.UUID, limit int) ([]Vocabulary, error)
		GetVocabulary(ctx context.Context, dictID uuid.UUID) ([]Vocabulary, error)
		UpdateWord(ctx context.Context, vocabulary Vocabulary) error
	}

	exampleSvc interface {
		AddExample(ctx context.Context, text, langCode string) (uuid.UUID, error)
		UpdateExample(ctx context.Context, text, langCode string) (uuid.UUID, error)
		GetExamples(ctx context.Context, exampleIDs []uuid.UUID) ([]entityExample.Example, error)
	}

	tagSvc interface {
		AddTag(ctx context.Context, text string) (uuid.UUID, error)
		UpdateTag(ctx context.Context, text string) (uuid.UUID, error)
		GetTags(ctx context.Context, tagIDs []uuid.UUID) ([]entityTag.Tag, error)
	}

	wordSvc interface {
		AddWord(ctx context.Context, id uuid.UUID, text, langCode, pronunciation string) (uuid.UUID, error)
		UpdateWord(ctx context.Context, text, langCode, pronunciation string) (uuid.UUID, error)
		GetWords(ctx context.Context, wordIDs []uuid.UUID) ([]entityWord.Word, error)
	}
)

type Service struct {
	repo       repoVocabulary
	wordSvc    wordSvc
	exampleSvc exampleSvc
	tagSvc     tagSvc
}

func NewService(
	repo repoVocabulary,
	wordSvc wordSvc,
	exexampleSvc exampleSvc,
	tagSvc tagSvc,
) *Service {
	return &Service{
		repo:       repo,
		wordSvc:    wordSvc,
		exampleSvc: exexampleSvc,
		tagSvc:     tagSvc,
	}
}

func (s *Service) AddWord(
	ctx context.Context,
	dictID uuid.UUID,
	nativeWord Word,
	tanslateWords Words,
	examples []string,
	tags []string) (*Vocabulary, error) {
	nativeWordID, err := s.wordSvc.AddWord(ctx, uuid.New(), nativeWord.Text, nativeWord.LangCode, nativeWord.Pronunciation)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.AddWord - add native word in dictionary: %w", err)
	}

	translateWordIDs := make([]uuid.UUID, 0, len(tanslateWords))
	for _, translateWord := range tanslateWords {
		translateID, err := s.wordSvc.AddWord(ctx, uuid.New(), translateWord.Text, translateWord.LangCode, translateWord.Pronunciation)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.AddWord - add translate word in dictionary: %w", err)
		}
		translateWordIDs = append(translateWordIDs, translateID)
	}

	exampleIDs := make([]uuid.UUID, 0, len(examples))
	for _, example := range examples {
		exampleID, err := s.exampleSvc.AddExample(ctx, example, nativeWord.LangCode)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.AddWord - add example: %w", err)
		}
		exampleIDs = append(exampleIDs, exampleID)
	}

	tagIDs := make([]uuid.UUID, 0, len(tags))
	for _, tag := range tags {
		tagID, err := s.tagSvc.AddTag(ctx, tag)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.AddWord - add tag: %w", err)
		}
		tagIDs = append(tagIDs, tagID)
	}

	vocabulary := Vocabulary{
		DictionaryId:   dictID,
		NativeWord:     nativeWordID,
		TranslateWords: translateWordIDs,
		Examples:       exampleIDs,
		Tags:           tagIDs,
	}

	err = s.repo.AddWord(ctx, vocabulary)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.AddWord - add vocabulary: %w", err)
	}

	return &vocabulary, nil
}

func (s *Service) UpdateWord(ctx context.Context,
	dictID uuid.UUID,
	oldWordID uuid.UUID,
	nativeWord Word,
	tanslateWords Words,
	examples []string,
	tags []string) (*VocabularyWord, error) {
	nativeWordID, err := s.wordSvc.UpdateWord(ctx, nativeWord.Text, nativeWord.LangCode, nativeWord.Pronunciation)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.UpdateWord - add native word in dictionary: %w", err)
	}

	translateWordIDs := make([]uuid.UUID, 0, len(tanslateWords))
	for _, translateWord := range tanslateWords {
		translateWordID, err := s.wordSvc.UpdateWord(ctx, translateWord.Text, translateWord.LangCode, translateWord.Pronunciation)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.UpdateWord - add translate word in dictionary: %w", err)
		}
		translateWordIDs = append(translateWordIDs, translateWordID)
	}

	exampleIDs := make([]uuid.UUID, 0, len(examples))
	for _, example := range examples {
		exampleID, err := s.exampleSvc.UpdateExample(ctx, example, nativeWord.LangCode)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.UpdateWord - add example: %w", err)
		}
		exampleIDs = append(exampleIDs, exampleID)
	}

	tagIDs := make([]uuid.UUID, 0, len(tags))
	for _, tag := range tags {
		tagID, err := s.tagSvc.UpdateTag(ctx, tag)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.UpdateWord - add tag: %w", err)
		}
		tagIDs = append(tagIDs, tagID)
	}

	vocabulary := Vocabulary{
		DictionaryId:   dictID,
		NativeWord:     nativeWordID,
		TranslateWords: translateWordIDs,
		Examples:       exampleIDs,
		Tags:           tagIDs,
	}

	if oldWordID != nativeWordID {
		err = s.repo.DeleteWord(ctx, Vocabulary{DictionaryId: dictID, NativeWord: oldWordID})
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.UpdateWord - delete old word: %w", err)
		}
		err = s.repo.AddWord(ctx, vocabulary)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.UpdateWord - add new word: %w", err)
		}
		return &VocabularyWord{
			NativeWord:     nativeWord,
			TranslateWords: tanslateWords.GetValues(),
			Examples:       examples,
			Tags:           tags,
		}, nil
	}
	err = s.repo.UpdateWord(ctx, vocabulary)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.UpdateWord - update vocabulary: %w", err)
	}

	return &VocabularyWord{
		NativeWord:     nativeWord,
		TranslateWords: tanslateWords.GetValues(),
		Examples:       examples,
		Tags:           tags,
	}, nil
}

func (s *Service) DeleteWord(ctx context.Context, dictID, nativeWordID uuid.UUID) error {
	vocabulary := Vocabulary{
		DictionaryId: dictID,
		NativeWord:   nativeWordID,
	}

	err := s.repo.DeleteWord(ctx, vocabulary)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.DeleteWord - delete word: %w", err)
	}
	return nil
}

func (s *Service) GetRandomWords(ctx context.Context, dictID uuid.UUID, limit int) ([]VocabularyWord, error) {
	vocabularies, err := s.repo.GetRandomVocabulary(ctx, dictID, limit)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.GetWords - get words: %w", err)
	}

	vocabularyWords := make([]VocabularyWord, 0, len(vocabularies))

	var eg errgroup.Group
	eg.SetLimit(10)
	for _, vocab := range vocabularies {
		vocab := vocab
		eg.Go(func() error {
			words, err := s.wordSvc.GetWords(ctx, []uuid.UUID{vocab.NativeWord})
			if err != nil {
				return fmt.Errorf("get words: %w", err)
			}
			if len(words) == 0 {
				return fmt.Errorf("not found word by id [%v]", vocab.NativeWord)
			}

			translateWords, err := s.wordSvc.GetWords(ctx, vocab.TranslateWords)
			if err != nil {
				return fmt.Errorf("get translate words: %w", err)
			}
			translates := make([]string, 0, len(translateWords))
			for _, word := range translateWords {
				translates = append(translates, word.Text)
			}

			examples, err := s.exampleSvc.GetExamples(ctx, vocab.Examples)
			if err != nil {
				return fmt.Errorf("get examples: %w", err)
			}
			examplesStr := make([]string, 0, len(examples))
			for _, example := range examples {
				examplesStr = append(examplesStr, example.Text)
			}

			tags, err := s.tagSvc.GetTags(ctx, vocab.Tags)
			if err != nil {
				return fmt.Errorf("get tags: %w", err)
			}
			tagsStr := make([]string, 0, len(tags))
			for _, tag := range tags {
				tagsStr = append(tagsStr, tag.Text)
			}

			vocabularyWords = append(vocabularyWords, VocabularyWord{
				NativeWord: Word{
					Text:          words[0].Text,
					Pronunciation: words[0].Pronunciation,
				},
				TranslateWords: translates,
				Examples:       examplesStr,
				Tags:           tagsStr,
			})
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("vocabulary.Service.GetWords - get words: %w", err)
	}

	return vocabularyWords, nil
}

func (s *Service) GetWord(ctx context.Context, dictID uuid.UUID, wordID uuid.UUID) (*VocabularyWord, error) {
	word, err := s.repo.GetWord(ctx, dictID, wordID)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.GetWord: %w", err)
	}
	var vocab VocabularyWord
	var eg errgroup.Group
	eg.Go(func() error {
		words, err := s.wordSvc.GetWords(ctx, []uuid.UUID{wordID})
		if err != nil {
			return fmt.Errorf("get words: %w", err)
		}
		if len(words) == 0 {
			return fmt.Errorf("not found word by id [%v]", wordID)
		}

		vocab.NativeWord = Word{
			Text:          words[0].Text,
			Pronunciation: words[0].Pronunciation,
		}

		return nil
	})
	eg.Go(func() error {
		translateWords, err := s.wordSvc.GetWords(ctx, word.TranslateWords)
		if err != nil {
			return fmt.Errorf("get translate words: %w", err)
		}
		vocab.TranslateWords = make([]string, 0, len(translateWords))
		for _, word := range translateWords {
			vocab.TranslateWords = append(vocab.TranslateWords, word.Text)
		}
		return nil
	})
	eg.Go(func() error {
		examples, err := s.exampleSvc.GetExamples(ctx, word.Examples)
		if err != nil {
			return fmt.Errorf("get examples: %w", err)
		}
		vocab.Examples = make([]string, 0, len(examples))
		for _, example := range examples {
			vocab.Examples = append(vocab.Examples, example.Text)
		}
		return nil
	})
	eg.Go(func() error {
		tags, err := s.tagSvc.GetTags(ctx, word.Tags)
		if err != nil {
			return fmt.Errorf("get tags: %w", err)
		}
		vocab.Tags = make([]string, 0, len(tags))
		for _, tag := range tags {
			vocab.Tags = append(vocab.Tags, tag.Text)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("vocabulary.Service.GetWords - get data: %w", err)
	}

	return &vocab, nil
}

func (s *Service) GetWords(ctx context.Context, dictID uuid.UUID) ([]VocabularyWord, error) {
	vocabularies, err := s.repo.GetVocabulary(ctx, dictID)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.GetWords - get words: %w", err)
	}

	vocabularyWords := make([]VocabularyWord, 0, len(vocabularies))

	var eg errgroup.Group
	eg.SetLimit(10)
	for _, vocab := range vocabularies {
		eg.Go(func() error {
			words, err := s.wordSvc.GetWords(ctx, []uuid.UUID{vocab.NativeWord})
			if err != nil {
				return fmt.Errorf("get words: %w", err)
			}
			if len(words) == 0 {
				return fmt.Errorf("not found word by id [%v]", vocab.NativeWord)
			}

			translateWords, err := s.wordSvc.GetWords(ctx, vocab.TranslateWords)
			if err != nil {
				return fmt.Errorf("get translate words: %w", err)
			}
			translates := make([]string, 0, len(translateWords))
			for _, word := range translateWords {
				translates = append(translates, word.Text)
			}

			examples, err := s.exampleSvc.GetExamples(ctx, vocab.Examples)
			if err != nil {
				return fmt.Errorf("get examples: %w", err)
			}
			examplesStr := make([]string, 0, len(examples))
			for _, example := range examples {
				examplesStr = append(examplesStr, example.Text)
			}

			tags, err := s.tagSvc.GetTags(ctx, vocab.Tags)
			if err != nil {
				return fmt.Errorf("get tags: %w", err)
			}
			tagsStr := make([]string, 0, len(tags))
			for _, tag := range tags {
				tagsStr = append(tagsStr, tag.Text)
			}

			vocabularyWords = append(vocabularyWords, VocabularyWord{
				Id: words[0].ID,
				NativeWord: Word{
					Text:          words[0].Text,
					Pronunciation: words[0].Pronunciation,
				},
				TranslateWords: translates,
				Examples:       examplesStr,
				Tags:           tagsStr,
			})
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("vocabulary.Service.GetWords - get words: %w", err)
	}

	return vocabularyWords, nil
}
