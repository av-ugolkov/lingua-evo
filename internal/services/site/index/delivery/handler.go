package index

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"lingua-evo/internal/config"
	entityLanguage "lingua-evo/internal/services/lingua/language/entity"
	dtoWord "lingua-evo/internal/services/lingua/word/dto"
	entityWord "lingua-evo/internal/services/lingua/word/entity"
	entitySession "lingua-evo/internal/services/session/entity"
	entityUser "lingua-evo/internal/services/user/entity"

	"lingua-evo/pkg/files"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/http/handler/common"
	"lingua-evo/pkg/token"
)

const (
	mainURL  = "/"
	indexURL = "/index"

	getAccountPanelURL = "/get-account-panel"

	indexPagePath    = "website/index.html"
	accountPanelPath = "website/components/account_panel.html"
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

	r.HandleFunc(getAccountPanelURL, h.getAccountPanel).Methods(http.MethodGet)
}

func (h *Handler) openPageIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, mainURL, http.StatusPermanentRedirect)
}

func (h *Handler) openPage(w http.ResponseWriter, r *http.Request) {
	handler := handler.NewHandler(w, r)

	t, err := files.ParseFiles(indexPagePath)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("site.index.delivery.Handler.get - parseFiles: %v", err))
		return
	}
	language := handler.GetCookieLanguageOrDefault()
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("site.index.delivery.Handler.get - Cookie: %v", err))
		return
	}
	randomWord, err := h.wordSvc.GetRandomWord(r.Context(), &dtoWord.RandomWordRq{LanguageCode: language})
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("site.index.delivery.Handler.get - GetRandomWord: %v", err))
		return
	}

	data := struct {
		Language *entityLanguage.Language
		Word     *entityWord.Word
	}{
		Language: &entityLanguage.Language{
			Code: common.Language,
		},
		Word: randomWord,
	}
	handler.SetCookieLanguage(common.Language)

	err = t.Execute(w, data)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("site.index.delivery.Handler.get - Execute: %v", err))
		return
	}
}

func (h *Handler) getAccountPanel(w http.ResponseWriter, r *http.Request) {
	handler := handler.NewHandler(w, r)

	accessToken, err := handler.GetHeaderAccessToken()
	if err != nil {
		handler.SendError(http.StatusBadRequest, fmt.Errorf("site.index.delivery.Handler.getAccountPanel - GetHeaderAccessToken: %v", err))
		return
	}
	fingerprint, err := handler.GetHeaderFingerprint()
	if err != nil {
		handler.SendError(http.StatusBadRequest, fmt.Errorf("site.index.delivery.Handler.getAccountPanel - GetHeaderFingerprint: %v", err))
	}
	claims, err := token.ValidateJWT(accessToken, config.GetConfig().JWT.Secret)
	if err != nil {
		handler.SendError(http.StatusBadRequest, fmt.Errorf("site.index.delivery.Handler.getAccountPanel - ValidateJWT: %v", err))
		return
	}

	ctx := r.Context()

	session, err := h.sessionSvc.GetSession(ctx, claims.ID)
	if err != nil {
		handler.SendError(http.StatusUnauthorized, err)
		return
	}
	if session.Fingerprint != fingerprint {
		handler.SendError(http.StatusUnauthorized, err)
		return
	}

	u, err := h.userSvc.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		handler.SendError(http.StatusUnauthorized, err)
		return
	}

	t, err := files.ParseFiles(accountPanelPath)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("site.index.delivery.Handler.getAccountPanel - parseFiles: %v", err))
		return
	}

	user := struct {
		Name string
	}{
		Name: u.Name,
	}

	err = t.Execute(w, user)
	if err != nil {
		handler.SendError(http.StatusInternalServerError, fmt.Errorf("site.index.delivery.Handler.getAccountPanel - Execute: %v", err))
		return
	}
}
