package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"time"

	pg "github.com/av-ugolkov/lingua-evo/internal/db/postgres"

	"github.com/google/uuid"
	"golang.org/x/text/language"
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
			Word          string `json:"word"`
			Pronunciation string `json:"pronunciation"`
		} `json:"dictionary"`
	}

	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return fmt.Errorf("fildDB.fillWord - unmarshal: %w", err)
	}

	insertQuery := fmt.Sprintf(`INSERT INTO "word_%s" (id, text, pronunciation, lang_code, created_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (text) DO UPDATE SET pronunciation=$3, created_at=$5;`, language.English.String())
	for _, d := range data.Dictionary {
		_, err := db.Exec(insertQuery, uuid.New(), d.Word, d.Pronunciation, language.English, time.Now().UTC())
		if err != nil {
			slog.Error(fmt.Errorf("fail insert word [%s]: %v", d.Word, err).Error())
		}
	}
	return nil
}
