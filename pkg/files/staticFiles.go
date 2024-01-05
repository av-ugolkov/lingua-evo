package files

import (
	"fmt"
	"html/template"
	"os"
)

var additionalPath string

func InitStaticFiles(path string) {
	additionalPath = path
}

// OpenFile: open static file and return array bytes.
func OpenFile(fileName string) ([]byte, error) {
	content, err := os.ReadFile(fmt.Sprintf("%s%s", additionalPath, fileName))
	if err != nil {
		return []byte{}, fmt.Errorf("pkg.files.OpenFile.ReadFile: %w", err)
	}

	return content, nil
}

func ParseFiles(fileNames ...string) (*template.Template, error) {
	for i := 0; i < len(fileNames); i++ {
		fileNames[i] = fmt.Sprintf("%s%s", additionalPath, fileNames[i])
	}
	return template.ParseFiles(fileNames...)
}
