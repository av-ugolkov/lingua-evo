package index

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	staticFiles "lingua-evo"
	entityLanguage "lingua-evo/internal/services/lingua/language/entity"
	dtoWord "lingua-evo/internal/services/lingua/word/dto"
	entityWord "lingua-evo/internal/services/lingua/word/entity"
	wordSvc "lingua-evo/internal/services/lingua/word/service"

	"lingua-evo/pkg/http/handler"
)

const (
	mainURL  = "/"
	indexURL = "/index"

	indexPagePath = "website/index.html"
)

type Handler struct {
	wordSvc *wordSvc.WordSvc
}

func Create(r *mux.Router, wordSvc *wordSvc.WordSvc) {
	handler := newHandler(wordSvc)
	handler.register(r)
}

func newHandler(wordSvc *wordSvc.WordSvc) *Handler {
	return &Handler{
		wordSvc: wordSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(mainURL, h.get).Methods(http.MethodGet)
	r.HandleFunc(indexURL, h.getIndex).Methods(http.MethodGet)
}

func (h *Handler) getIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, mainURL, http.StatusPermanentRedirect)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	t, err := staticFiles.ParseFiles(indexPagePath)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("site.index.delivery.Handler.get - parseFiles: %v", err))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	randomWord, err := h.wordSvc.GetRandomWord(r.Context(), &dtoWord.RandomWordRq{LanguageCode: "en"})
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("site.index.delivery.Handler.get - GetRandomWord: %v", err))
		return
	}
	data := struct {
		Language *entityLanguage.Language
		Word     *entityWord.Word
	}{
		//TODO нужно то ли с браузера, то ли еще откуда-то брать язык
		Language: &entityLanguage.Language{
			Code: "en",
		},
		Word: randomWord,
	}

	err = t.Execute(w, data)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("site.index.delivery.Handler.get - Execute: %v", err))
		return
	}
}
