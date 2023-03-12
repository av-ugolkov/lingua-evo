package add_word

import (
	"lingua-evo/internal/api/entity"
	"net/http"
	"os"
)

const (
	addWordPagePath = entity.RootPath + "/add_word/add_word.html"
)

func (h *Handler) addWord(w http.ResponseWriter, _ *http.Request) {
	file, err := os.ReadFile(addWordPagePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = w.Write(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
