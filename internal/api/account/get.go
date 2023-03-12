package account

import (
	"lingua-evo/internal/api/entity"
	"net/http"
	"os"
)

const (
	accountPagePath = entity.RootPath + "/account/account.html"
)

func (h *Handler) account(w http.ResponseWriter, _ *http.Request) {
	file, err := os.ReadFile(accountPagePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = w.Write(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
