package static

import (
	"embed"
	"fmt"
	"html/template"
	"io"
)

//go:embed website
var files embed.FS

// OpenFile: open static file and return array bytes.
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
