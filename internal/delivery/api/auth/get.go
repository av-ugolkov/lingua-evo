package auth

import (
	"html/template"
	"net/http"

	"lingua-evo/internal/config"
	"lingua-evo/pkg/tools/view"
)

const (
	authPagePath = "view/auth/auth.html"
)

func (h *Handler) getAuth(w http.ResponseWriter, _ *http.Request) {
	t, err := template.ParseFiles(view.GetPathFile(authPagePath))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rootPath := struct {
		Root string
	}{
		Root: config.GetConfig().Front.Root,
	}

	err = t.Execute(w, rootPath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
