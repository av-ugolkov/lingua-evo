package sign_up

import (
	"net/http"
	"os"
)

const (
	signupPage = "./web/static/sign_up/signup.html"
)

func (h *Handler) getSignup(w http.ResponseWriter, _ *http.Request) {
	file, err := os.ReadFile(signupPage)
	if err != nil {
		h.logger.Errorf("auth.getSignup.ReadFile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(file))
}
