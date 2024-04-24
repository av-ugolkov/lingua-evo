package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	"github.com/google/uuid"
	"io"
	"log/slog"

	pg "github.com/av-ugolkov/lingua-evo/internal/db/postgres"
	"github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	dictRepo "github.com/av-ugolkov/lingua-evo/internal/services/dictionary/delivery/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/example"
	repoExample "github.com/av-ugolkov/lingua-evo/internal/services/example/delivery/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/language"
	repoLang "github.com/av-ugolkov/lingua-evo/internal/services/language/delivery/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/tag"
	repoTag "github.com/av-ugolkov/lingua-evo/internal/services/tag/delivery/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/user"
	repoUser "github.com/av-ugolkov/lingua-evo/internal/services/user/delivery/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	repoVocabulary "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/delivery/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/word"
	repoWord "github.com/av-ugolkov/lingua-evo/internal/services/word/delivery/repository"
)

//go:embed words_en.json
var words embed.FS

func main() {
	connStr := "postgresql://lingua:ib6vACdec2Fmht4lnX153@localhost:6432/pg-lingua-evo"

	db, err := pg.NewDB(connStr)
	if err != nil {
		slog.Error(fmt.Errorf("can't create pg pool: %v", err).Error())
		return
	}

	err = fillWord(db)
	if err != nil {
		slog.Error(fmt.Errorf("can't fill db: %v", err).Error())
	}
}

func fillWord(db *sql.DB) error {
	file, err := words.Open("words_en.json")
	if err != nil {
		return fmt.Errorf("fildDB.fillWord - Open: %w", err)
	}

	jsonFile, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("fildDB.fillWord - ReadAll: %w", err)
	}

	var data struct {
		Dictionary []struct {
			Word          string   `json:"word"`
			Pronunciation string   `json:"pronunciation,omitempty"`
			Examples      []string `json:"example,omitempty"`
			Translates    []string `json:"translate,omitempty"`
			Description   string   `json:"description,omitempty"`
		} `json:"dictionary"`
	}

	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return fmt.Errorf("fildDB.fillWord - unmarshal: %w", err)
	}

	userSvc := user.NewService(repoUser.NewRepo(db), nil)
	userData, err := userSvc.GetUser(context.Background(), "makedonskiy")
	if err != nil {
		return fmt.Errorf("fildDB.fillWord - GetUser: %w", err)
	}

	fmt.Println(userData)
	tr := transactor.NewTransactor(db)

	langSvc := language.NewService(repoLang.NewRepo(db))
	exampleSvc := example.NewService(repoExample.NewRepo(db))
	tagSvc := tag.NewService(repoTag.NewRepo(db))
	repoVocab := repoVocabulary.NewRepo(db)
	dictSvc := dictionary.NewService(dictRepo.NewRepo(db), langSvc)
	vocabSvc := vocabulary.NewService(tr, repoVocab, langSvc, tagSvc)
	wordSvc := word.NewService(tr, repoWord.NewRepo(db), vocabSvc, dictSvc, exampleSvc)

	vocab, err := vocabSvc.AddVocabulary(context.Background(), vocabulary.Vocabulary{
		ID:            uuid.New(),
		UserID:        userData.ID,
		Name:          "default",
		NativeLang:    "en",
		TranslateLang: "ru",
	})
	if err != nil {
		return fmt.Errorf("fildDB.fillWord - AddDictionary: %w", err)
	}

	for _, d := range data.Dictionary {
		translates := make([]dictionary.DictWord, 0, len(d.Translates))
		for _, text := range d.Translates {
			translates = append(translates, dictionary.DictWord{
				ID:            uuid.New(),
				Text:          text,
				Pronunciation: "",
				LangCode:      vocab.TranslateLang,
			})
		}

		examples := make([]example.Example, 0, len(d.Examples))
		for _, text := range d.Examples {
			examples = append(examples, example.Example{
				ID:   uuid.New(),
				Text: text,
			})
		}

		vocabWord, err := wordSvc.AddWord(context.Background(), word.VocabWordData{
			ID:      uuid.New(),
			VocabID: vocab.ID,
			Native: dictionary.DictWord{
				ID:            uuid.New(),
				Text:          d.Word,
				Pronunciation: d.Pronunciation,
				LangCode:      vocab.NativeLang,
			},
			Translates: translates,
			Examples:   examples,
		})
		if err != nil {
			slog.Error(fmt.Errorf("fail insert word [%s]: %v", d.Word, err).Error())
			continue
		}

		slog.Info(fmt.Sprintf("inserted word [%s]", vocabWord.ID))
	}
	return nil
}
