package sign_in

import (
	"net/http"
	"os"
)

const (
	signInPage = "./web/static/sign_in/signin.html"
)

func (h *Handler) getSignIn(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile(signInPage)
	if err != nil {
		h.logger.Errorf("auth.getSignIn.ReadFile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(file))
}
