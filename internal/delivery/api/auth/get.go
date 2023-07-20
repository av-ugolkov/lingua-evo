package auth

import (
	"net/http"
	"os"
)

const (
	authPagePath   = "./../web/static/auth/auth.html" //it is ok
	signupPagePath = "./../web/static/signup/signup.html"
)

func (h *Handler) getAuth(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile(authPagePath)
	if err != nil {
		h.logger.Errorf("auth.getAuth.ReadFile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(file))
}

func (h *Handler) getSignup(w http.ResponseWriter, _ *http.Request) {
	file, err := os.ReadFile(signupPagePath)
	if err != nil {
		h.logger.Errorf("auth.getSignup.ReadFile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(file))
}
