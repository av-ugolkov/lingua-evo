package auth

import (
	"lingua-evo/internal/api/entity"
	"net/http"
	"os"
)

const (
	authPagePath = entity.RootPath + "/auth/auth.html"
)

func (h *Handler) authGet(w http.ResponseWriter, _ *http.Request) {
	file, err := os.ReadFile(authPagePath)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = w.Write(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
