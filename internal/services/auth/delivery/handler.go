package delivery

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"lingua-evo/internal/config"
	"lingua-evo/internal/services/auth/dto"
	"lingua-evo/internal/services/auth/service"
	"lingua-evo/pkg/http/handler"
	"lingua-evo/pkg/middleware"
	"lingua-evo/pkg/token"
	"lingua-evo/runtime"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	login   = "/auth/login"
	refresh = "/auth/refresh"
	logout  = "/auth/logout"
)

type (
	Handler struct {
		authSvc *service.AuthSvc
	}
)

func Create(r *mux.Router, authSvc *service.AuthSvc) {
	handler := newHandler(authSvc)
	handler.register(r)
}

func newHandler(authSvc *service.AuthSvc) *Handler {
	return &Handler{
		authSvc: authSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(login, h.login).Methods(http.MethodPost)
	r.HandleFunc(refresh, middleware.Auth(h.refresh)).Methods(http.MethodPost)
	r.HandleFunc(logout, middleware.Auth(h.logout)).Methods(http.MethodPost)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data dto.CreateSessionRq
	err := decodeBasicAuth(r.Header.Get("Authorization"), &data)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.createSession - check body: %v", err))
		return
	}
	data.Fingerprint = r.Header.Get("Fingerprint")

	ctx := r.Context()
	tokens, err := h.authSvc.Login(ctx, &data)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.createSession - create session: %v", err))
		return
	}

	b, err := json.Marshal(&dto.CreateSessionRs{
		AccessToken: tokens.AccessToken,
	})
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.createSession - marshal: %v", err))
		return
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken.String(),
		MaxAge:   int(duration.Seconds()),
		HttpOnly: true,
		Secure:   true,
		Path:     "/auth",
	})

	_, _ = w.Write(b)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	refreshToken, err := handler.GetCookie(r, "refresh_token")
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - get cookie: %v", err))
		return
	}

	tokenID, err := uuid.Parse(refreshToken.Value)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - get cookie: %v", err))
		return
	}

	accessToken, err := handler.GetHeader(r, "access_token")
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - get cookie: %v", err))
		return
	}

	claims, err := token.ValidateJWT(accessToken, config.GetConfig().JWT.Secret)
	if err != nil {
		handler.SendError(w, http.StatusUnauthorized, err)
		return
	}

	ctx := r.Context()

	tokens, err := h.authSvc.RefreshSessionToken(ctx, claims.UserID, tokenID)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - RefreshSessionToken: %v", err))
		return
	}
	b, err := json.Marshal(&dto.CreateSessionRs{
		AccessToken: tokens.AccessToken,
	})
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - marshal: %v", err))
		return
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken.String(),
		MaxAge:   int(duration.Seconds()),
		HttpOnly: true,
		Secure:   true,
		Path:     "/auth",
	})

	_, _ = w.Write(b)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	ctx := r.Context()
	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		handler.SendError(w, http.StatusUnauthorized, fmt.Errorf("auth.delivery.Handler.logout - unauthorized: %v", err))
		return
	}

	err = h.authSvc.Logout(ctx, uid)
	if err != nil {
		handler.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.logout - logout: %v", err))
		return
	}

	_, _ = w.Write([]byte("done"))
}

func decodeBasicAuth(auth string, data *dto.CreateSessionRq) error {
	base, err := base64.StdEncoding.DecodeString(strings.Split(auth, " ")[1])
	if err != nil {
		return fmt.Errorf("auth.delivery.decodeBasicAuth - decode base64: %v", err)
	}
	authData := strings.Split(string(base), ":")
	if len(authData) != 2 {
		return fmt.Errorf("auth.delivery.decodeBasicAuth - invalid auth data")
	}

	data.User = authData[0]
	data.Password = authData[1]

	return nil
}
