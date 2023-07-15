package auth

import (
	"html/template"
	"net/http"

	"lingua-evo/internal/config"
	"lingua-evo/pkg/tools/view"
)

const (
	authPagePath   = "view/auth/auth.html"
	signupPagePath = "view/signup/signup.html"
)

func (h *Handler) getAuth(w http.ResponseWriter, _ *http.Request) {
	t, err := template.ParseFiles(view.GetPathFile(authPagePath))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data := struct {
		RootPath string
	}{
		RootPath: config.GetConfig().Front.Root,
	}

	err = t.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (h *Handler) getSignup(w http.ResponseWriter, _ *http.Request) {
	t, err := template.ParseFiles(view.GetPathFile(signupPagePath))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data := struct {
		RootPath string
	}{
		RootPath: config.GetConfig().Front.Root,
	}

	err = t.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
