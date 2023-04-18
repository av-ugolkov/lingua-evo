package add_word

import (
	"html/template"
	"net/http"

	"lingua-evo/pkg/tools/view"
)

const (
	addWordPagePath = "view/add_word/add_word.html"
)

func (h *Handler) getAddWord(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(view.GetPathFile(addWordPagePath))
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
