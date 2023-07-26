package static

import (
	"embed"
	"fmt"
	"html/template"
	"io"
)

//go:embed web
var files embed.FS

// Files returns a filesystem with static files.
func OpenFile(fileName string) ([]byte, error) {
	file, err := files.Open(fileName)
	if err != nil {
		return []byte{}, fmt.Errorf("templates.OpenFile.Open: %w", err)
	}

	body, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, fmt.Errorf("templates.OpenFile.ReadAll: %w", err)
	}

	return body, nil
}

func ParseFiles(fileNames ...string) (*template.Template, error) {
	return template.ParseFS(files, fileNames...)
}
