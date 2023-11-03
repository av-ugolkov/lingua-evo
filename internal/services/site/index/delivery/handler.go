package index

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	entityLanguage "lingua-evo/internal/services/lingua/language/entity"
	dtoWord "lingua-evo/internal/services/lingua/word/dto"
	entityWord "lingua-evo/internal/services/lingua/word/entity"
	entityUser "lingua-evo/internal/services/user/entity"

	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/http/static"
	"lingua-evo/runtime"
)

const (
	mainURL  = "/"
	indexURL = "/index"

	indexPagePath = "website/index.html"
)

type (
	userSvc interface {
		GetUserByRefreshToken(ctx context.Context, token uuid.UUID) (*entityUser.User, error)
	}

	wordSvc interface {
		GetRandomWord(ctx context.Context, w *dtoWord.RandomWordRq) (*entityWord.Word, error)
	}

	userInfo struct {
		IsLogin bool
		Name    string
	}

	Handler struct {
		userSvc userSvc
		wordSvc wordSvc
	}
)

func Create(r *mux.Router, userSvc userSvc, wordSvc wordSvc) {
	handler := newHandler(userSvc, wordSvc)
	handler.register(r)
}

func newHandler(userSvc userSvc, wordSvc wordSvc) *Handler {
	return &Handler{
		userSvc: userSvc,
		wordSvc: wordSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(mainURL, h.openPage).Methods(http.MethodGet)
	r.HandleFunc(indexURL, h.openPageIndex).Methods(http.MethodGet)
}

func (h *Handler) openPageIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, mainURL, http.StatusPermanentRedirect)
}

func (h *Handler) openPage(w http.ResponseWriter, r *http.Request) {
	t, err := static.ParseFiles(indexPagePath)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("site.index.delivery.Handler.get - parseFiles: %v", err))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	language := runtime.GetLanguage("en")
	languageCookie, err := handler.GetCookie(r, "language")
	if err != nil {
		slog.Warn(err.Error())
	}
	if languageCookie != nil {
		language = languageCookie.Value
	}

	user := &userInfo{
		IsLogin: false,
	}
	cookie, err := r.Cookie("refresh_token")
	if errors.Is(err, http.ErrNoCookie) {
		slog.Warn(fmt.Sprintf("not found cookie: %s", "token"))
	} else if err != nil {
		slog.Error(fmt.Errorf("site.index.delivery.Handler.get - GetCookie: %v", err).Error())
	}
	if cookie != nil {
		u, err := h.userSvc.GetUserByRefreshToken(r.Context(), uuid.MustParse(cookie.Value))
		if err != nil {
			slog.Error(fmt.Errorf("site.index.delivery.Handler.get - not found user by token [%s]: %v", cookie.Value, err).Error())
		} else {
			user.IsLogin = true
			user.Name = u.Username
		}
	}

	randomWord, err := h.wordSvc.GetRandomWord(r.Context(), &dtoWord.RandomWordRq{LanguageCode: language})
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("site.index.delivery.Handler.get - GetRandomWord: %v", err))
		return
	}

	data := struct {
		Language *entityLanguage.Language
		Word     *entityWord.Word
		User     *userInfo
	}{
		Language: &entityLanguage.Language{
			Code: language,
		},
		Word: randomWord,
		User: user,
	}

	handler.SetCookie(w, "language", language)
	err = t.Execute(w, data)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("site.index.delivery.Handler.get - Execute: %v", err))
		return
	}
}
