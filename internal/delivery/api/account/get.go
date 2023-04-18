package account

import (
	"net/http"
	"os"

	"lingua-evo/pkg/tools/view"
)

const (
	accountPagePath = "view/account/account.html"
)

func (h *Handler) account(w http.ResponseWriter, _ *http.Request) {
	file, err := os.ReadFile(view.GetPathFile(accountPagePath))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = w.Write(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
