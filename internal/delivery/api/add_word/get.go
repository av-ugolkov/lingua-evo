package add_word

import (
	"html/template"
	"net/http"

	"lingua-evo/internal/config"
	"lingua-evo/internal/service"
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	data := struct {
		Root      string
		Languages []*service.Language
	}{
		Root:      config.GetConfig().Front.Root,
		Languages: languages,
	}

	err = t.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
