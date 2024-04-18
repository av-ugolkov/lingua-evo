package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/av-ugolkov/lingua-evo/internal/services/auth/model"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"net/http"
	"strings"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/delivery"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/http/exchange"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/services/auth"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	authSvc *auth.Service
}

func Create(r *mux.Router, authSvc *auth.Service) {
	h := newHandler(authSvc)
	h.register(r)
}

func newHandler(authSvc *auth.Service) *Handler {
	return &Handler{
		authSvc: authSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(delivery.SignUp, h.signUp).Methods(http.MethodPost)
	r.HandleFunc(delivery.SignIn, h.signIn).Methods(http.MethodPost)
	r.HandleFunc(delivery.Refresh, h.refresh).Methods(http.MethodGet)
	r.HandleFunc(delivery.SignOut, middleware.Auth(h.signOut)).Methods(http.MethodGet)
	r.HandleFunc(delivery.SendCode, h.sendCode).Methods(http.MethodPost)
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	ex := exchange.NewExchanger(w, r)
	var data model.CreateUserRq
	err := ex.CheckBody(&data)
	if err != nil {
		ex.SendError(http.StatusBadRequest, fmt.Errorf("user.delivery.Handler.createAccount - check body: %v", err))
		return
	}

	uid, err := h.authSvc.SignUp(ex.Context(), data, runtime.User)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("user.delivery.Handler.createAccount - create user: %v", err))
		return
	}

	createUserRs := &model.CreateUserRs{
		UserID: uid,
	}

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SendData(http.StatusCreated, createUserRs)
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	ex := exchange.NewExchanger(w, r)
	authorization, err := ex.GetHeaderAuthorization(exchange.AuthTypeBasic)
	if err != nil {
		ex.SendError(http.StatusBadRequest, fmt.Errorf("auth.delivery.Handler.signin: %v", err))
		return
	}
	var data model.CreateSessionRq
	err = decodeBasicAuth(authorization, &data)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.signin - check body: %v", err))
		return
	}
	data.Fingerprint, err = ex.GetHeaderFingerprint()
	if err != nil {
		ex.SendError(http.StatusBadRequest, fmt.Errorf("auth.delivery.Handler.signin - GetHeaderFingerprint: %v", err))
		return
	}
	ctx := r.Context()
	tokens, err := h.authSvc.SignIn(ctx, data.User, data.Password, data.Fingerprint)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.signin - create session: %v", err))
		return
	}

	sessionRs := &model.CreateSessionRs{
		AccessToken: tokens.AccessToken,
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second

	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SetCookieRefreshToken(tokens.RefreshToken, duration)
	ex.SendData(http.StatusOK, sessionRs)
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	ex := exchange.NewExchanger(w, r)
	refreshToken, err := ex.Cookie(exchange.RefreshToken)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - get cookie: %v", err))
		return
	}
	if refreshToken == nil {
		ex.SendError(http.StatusUnauthorized, fmt.Errorf("auth.delivery.Handler.refresh - refresh token is nil"))
		return
	}

	refreshID, err := uuid.Parse(refreshToken.Value)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - get cookie: %v", err))
		return
	}

	fingerprint, err := ex.GetHeaderFingerprint()
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - get fingerprint: %v", err))
		return
	}

	ctx := r.Context()
	tokens, err := h.authSvc.RefreshSessionToken(ctx, refreshID, fingerprint)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.refresh - refresh token: %v", err))
		return
	}

	sessionRs := &model.CreateSessionRs{
		AccessToken: tokens.AccessToken,
	}

	additionalTime := config.GetConfig().JWT.ExpireRefresh
	duration := time.Duration(additionalTime) * time.Second
	ex.SetContentType(exchange.ContentTypeJSON)
	ex.SetCookieRefreshToken(tokens.RefreshToken, duration)
	ex.SendData(http.StatusOK, sessionRs)
}

func (h *Handler) signOut(ctx context.Context, ex *exchange.Exchanger) {
	refreshToken, err := ex.Cookie(exchange.RefreshToken)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.logout - get cookie: %v", err))
		return
	}
	if refreshToken == nil {
		ex.SendError(http.StatusUnauthorized, fmt.Errorf("auth.delivery.Handler.logout: %v", err))
		return
	}

	refreshID, err := uuid.Parse(refreshToken.Value)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.logout - parse: %v", err))
		return
	}

	fingerprint, err := ex.GetHeaderFingerprint()
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.logout - get fingerprint: %v", err))
		return
	}

	err = h.authSvc.SignOut(ctx, refreshID, fingerprint)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.logout - logout: %v", err))
		return
	}

	ex.DeleteCookie(exchange.RefreshToken)
	ex.SendEmptyData(http.StatusOK)
}

func (h *Handler) sendCode(w http.ResponseWriter, r *http.Request) {
	ex := exchange.NewExchanger(w, r)

	var data model.CreateCodeRq
	err := ex.CheckBody(&data)
	if err != nil {
		ex.SendError(http.StatusBadRequest, fmt.Errorf("auth.delivery.handler.Handler.sendCode - check body: %v", err))
		return
	}

	err = h.authSvc.CreateCode(ex.Context(), data.Email)
	if err != nil {
		ex.SendError(http.StatusInternalServerError, fmt.Errorf("auth.delivery.Handler.sendCode - create code: %v", err))
		return
	}

	ex.SendEmptyData(http.StatusOK)
}

func decodeBasicAuth(basicToken string, data *model.CreateSessionRq) error {
	base, err := base64.StdEncoding.DecodeString(basicToken)
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
