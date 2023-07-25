package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io"
)

//go:embed index.html
//go:embed sign_in/signin.html
//go:embed sign_up/signup.html
//go:embed dictionary/add_word/add_word.html
//go:embed account/account.html
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
