package add_word

import (
	"html/template"
	"lingua-evo/internal/api/entity"
	"net/http"
)

const (
	addWordPagePath = entity.RootPath + "/add_word/add_word.html"
)

func (h *Handler) getAddWord(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(addWordPagePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	languages, err := h.lingua.GetLanguages(r.Context())
	if err != nil {
		return
	}
	err = t.Execute(w, languages)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
