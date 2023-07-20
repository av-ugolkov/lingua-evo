package add_word

import (
	"html/template"
	"net/http"

	"lingua-evo/internal/service"
)

const (
	addWordPagePath = "./../view/add_word/add_word.html"
)

func (h *Handler) getAddWord(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(addWordPagePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	languages, err := h.lingua.GetLanguages(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	data := struct {
		Languages []*service.Language
	}{
		Languages: languages,
	}

	err = t.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
