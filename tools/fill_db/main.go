package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	pg "github.com/av-ugolkov/lingua-evo/internal/db/postgres"
	"github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	dictRepo "github.com/av-ugolkov/lingua-evo/internal/services/dictionary/delivery/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/example"
	repoExample "github.com/av-ugolkov/lingua-evo/internal/services/example/delivery/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/tag"
	repoTag "github.com/av-ugolkov/lingua-evo/internal/services/tag/delivery/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/user"
	repoUser "github.com/av-ugolkov/lingua-evo/internal/services/user/delivery/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	repoVocabulary "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/delivery/repository"
	"github.com/av-ugolkov/lingua-evo/internal/services/word"
	repoWord "github.com/av-ugolkov/lingua-evo/internal/services/word/delivery/repository"
	"github.com/google/uuid"
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
	user, err := userSvc.GetUser(context.Background(), "makedonskiy")
	if err != nil {
		return fmt.Errorf("fildDB.fillWord - GetUser: %w", err)
	}

	wordSvc := word.NewService(repoWord.NewRepo(db))
	exampleSvc := example.NewService(repoExample.NewRepo(db))
	tagSvc := tag.NewService(repoTag.NewRepo(db))
	repoVocab := repoVocabulary.NewRepo(db)
	dictSvc := dictionary.NewService(dictRepo.NewRepo(db), repoVocab)
	vocabSvc := vocabulary.NewService(repoVocab, wordSvc, exampleSvc, tagSvc)

	dictID, err := dictSvc.AddDictionary(context.Background(), user.ID, uuid.New(), "default")
	if err != nil {
		return fmt.Errorf("fildDB.fillWord - AddDictionary: %w", err)
	}

	const (
		nativeTable    = "en"
		translateTable = "ru"
	)

	for _, d := range data.Dictionary {
		translateWords := make([]vocabulary.Word, 0, len(d.Translates))
		for _, word := range d.Translates {
			translateWords = append(translateWords, vocabulary.Word{Text: word, Pronunciation: "", LangCode: translateTable})
		}
		vocabulary, err := vocabSvc.AddWord(context.Background(), dictID, vocabulary.Word{Text: d.Word, Pronunciation: d.Pronunciation, LangCode: nativeTable},
			translateWords, d.Examples, nil)
		if err != nil {
			slog.Error(fmt.Errorf("fail insert word [%s]: %v", d.Word, err).Error())
			continue
		}

		slog.Info(fmt.Sprintf("inserted word [%s]", vocabulary.NativeWord))
	}
	return nil
}
