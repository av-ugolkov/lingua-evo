package dictionary

import (
	"LinguaEvo/storage"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Dictionary struct {
	basePath string
}

const defaultPerm = 0774

func New(basePath string) *Dictionary {
	return &Dictionary{basePath: basePath}
}

func (d Dictionary) AddWord(w *storage.Word) error {
	fPath := filepath.Join(d.basePath, w.UserName)
	log.Printf("fPath: %s", fPath)
	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(w)
	if err != nil {
		return fmt.Errorf("words.AddWord.fileName: %w", err)
	}

	fPath = filepath.Join(fPath, fName)
	file, err := os.Create(fPath)
	if err != nil {
		return fmt.Errorf("words.AddWord.Create: %w", err)
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(w); err != nil {
		return fmt.Errorf("words.AddWord.NewEncoder: %w", err)
	}

	return nil
}

func (d Dictionary) EditWord(w *storage.Word) error {
	log.Println(w)
	return nil
}

func (d Dictionary) RemoveWord(w *storage.Word) error {
	fileName, err := fileName(w)
	if err != nil {
		return fmt.Errorf("words.RemoveWord.fileName: %w", err)
	}

	path := filepath.Join(d.basePath, w.UserName, fileName)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("words.RemoveWord.Remove [%s]: %w", path, err)
	}

	return nil
}

func (d Dictionary) PickRandomWord(w *storage.Word) (*storage.Word, error) {
	path := filepath.Join(d.basePath, w.UserName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, storage.ErrNoSavePages
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("words.PickRandomWord.ReadDir: %w", err)
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavePages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return d.decodePage(filepath.Join(path, file.Name()))
}

func (d Dictionary) SharedWord(w *storage.Word) (*storage.Word, error) {
	log.Println(w)
	return nil, nil
}

func (d Dictionary) IsExists(w *storage.Word) (bool, error) {
	fileName, err := fileName(w)
	if err != nil {
		return false, fmt.Errorf("words.IsExists.fileName: %w", err)
	}

	path := filepath.Join(d.basePath, w.UserName, fileName)

	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("words.IsExists.Stat [%s]: %w", path, err)
	}

	return true, nil
}

func (d Dictionary) decodePage(filePath string) (*storage.Word, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("words.decodePage.Open: %w", err)
	}
	defer func() { _ = f.Close() }()

	var word storage.Word
	if err := gob.NewDecoder(f).Decode(&word); err != nil {
		return nil, fmt.Errorf("words.decodePage.Decode: %w", err)
	}

	return &word, nil
}

func fileName(p *storage.Word) (string, error) {
	return p.Hash()
}
