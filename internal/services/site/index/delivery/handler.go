package index

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"lingua-evo/internal/config"
	entityLanguage "lingua-evo/internal/services/lingua/language/entity"
	dtoWord "lingua-evo/internal/services/lingua/word/dto"
	entityWord "lingua-evo/internal/services/lingua/word/entity"
	entitySession "lingua-evo/internal/services/session/entity"
	entityUser "lingua-evo/internal/services/user/entity"

	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/http/static"
	"lingua-evo/pkg/token"
	"lingua-evo/runtime"
)

const (
	mainURL  = "/"
	indexURL = "/index"

	getAccountDataURL = "/get-account-data"

	indexPagePath = "website/index.html"
)

type (
	sessionSvc interface {
		GetSession(ctx context.Context, sid string) (*entitySession.Session, error)
	}

	userSvc interface {
		GetUserByID(ctx context.Context, uid uuid.UUID) (*entityUser.User, error)
	}

	wordSvc interface {
		GetRandomWord(ctx context.Context, w *dtoWord.RandomWordRq) (*entityWord.Word, error)
	}

	Handler struct {
		sessionSvc sessionSvc
		userSvc    userSvc
		wordSvc    wordSvc
	}
)

func Create(r *mux.Router, sessionSvc sessionSvc, userSvc userSvc, wordSvc wordSvc) {
	handler := newHandler(sessionSvc, userSvc, wordSvc)
	handler.register(r)
}

func newHandler(sessionSvc sessionSvc, userSvc userSvc, wordSvc wordSvc) *Handler {
	return &Handler{
		sessionSvc: sessionSvc,
		userSvc:    userSvc,
		wordSvc:    wordSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(mainURL, h.openPage).Methods(http.MethodGet)
	r.HandleFunc(indexURL, h.openPageIndex).Methods(http.MethodGet)

	r.HandleFunc(getAccountDataURL, h.getAccountData).Methods(http.MethodGet)
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

	randomWord, err := h.wordSvc.GetRandomWord(r.Context(), &dtoWord.RandomWordRq{LanguageCode: language})
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("site.index.delivery.Handler.get - GetRandomWord: %v", err))
		return
	}

	data := struct {
		Language *entityLanguage.Language
		Word     *entityWord.Word
	}{
		Language: &entityLanguage.Language{
			Code: language,
		},
		Word: randomWord,
	}

	handler.SetCookie(w, "language", language)
	err = t.Execute(w, data)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("site.index.delivery.Handler.get - Execute: %v", err))
		return
	}
}

func (h *Handler) getAccountData(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Access-Token")
	fingerprint := r.Header.Get("Fingerprint")

	claims, err := token.ValidateJWT(accessToken, config.GetConfig().JWT.Secret)
	if err != nil {
		handler.SendError(w, http.StatusUnauthorized, err)
		return
	}

	ctx := r.Context()

	session, err := h.sessionSvc.GetSession(ctx, claims.ID)
	if err != nil {
		handler.SendError(w, http.StatusUnauthorized, err)
		return
	}
	if session.Fingerprint != fingerprint {
		handler.SendError(w, http.StatusUnauthorized, err)
		return
	}

	u, err := h.userSvc.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		handler.SendError(w, http.StatusUnauthorized, err)
		return
	}
	fmt.Print(u)
	//static.ParseFiles(indexPagePath)

	//_, _ = w.Write(b)
}
