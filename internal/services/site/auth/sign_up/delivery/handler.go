package sign_up

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"lingua-evo/internal/config"
	entityUserSvc "lingua-evo/internal/services/lingua/user/entity"
	userSvc "lingua-evo/internal/services/lingua/user/service"
	"lingua-evo/internal/services/site/auth/sign_up/entity"
	httpTools "lingua-evo/internal/tools"

	staticFiles "lingua-evo"
	linguaJWT "lingua-evo/pkg/middleware/jwt"
	"lingua-evo/pkg/tools"

	"github.com/cristalhq/jwt/v3"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const (
	signUpURL = "/signup"

	signupPage = "website/sign_up/signup.html"
)

type Handler struct {
	userSvc *userSvc.UserSvc
}

func Create(r *mux.Router, userSvc *userSvc.UserSvc) {
	handler := newHandler(userSvc)
	handler.register(r)
}

func newHandler(userSvc *userSvc.UserSvc) *Handler {
	return &Handler{
		userSvc: userSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(signUpURL, h.get).Methods(http.MethodGet)
	r.HandleFunc(signUpURL, h.post).Methods(http.MethodPost)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	file, err := staticFiles.OpenFile(signupPage)
	if err != nil {
		slog.Error(fmt.Errorf("sign_up.get.OpenFile: %v", err).Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(file))
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	var data entity.UserRq
	err := httpTools.CheckBody(w, r, &data)
	if err != nil {
		httpTools.SendError(w, http.StatusBadRequest, fmt.Errorf("auth.sign_up.Handler.post - check body: %v", err))
		return
	}
	if err := h.validateEmail(r.Context(), data.Email); err != nil {
		httpTools.SendError(w, http.StatusConflict, fmt.Errorf("auth.sign_up.Handler.post - validateEmail: %v", err))
		return
	}

	if err := h.validateUsername(r.Context(), data.Username); err != nil {
		httpTools.SendError(w, http.StatusConflict, fmt.Errorf("auth.sign_up.Handler.post - validateUsername: %v", err))
		return
	}

	if err := validatePassword(data.Password); err != nil {
		httpTools.SendError(w, http.StatusConflict, fmt.Errorf("auth.sign_up.Handler.post - validatePassword: %v", err))
		return
	}

	hashPassword, err := hashPassword(data.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	uid, err := h.userSvc.CreateUser(r.Context(), &entityUserSvc.User{
		Username:     data.Username,
		PasswordHash: hashPassword,
		Email:        data.Email,
	})
	if err != nil {
		httpTools.SendError(w, http.StatusInternalServerError, fmt.Errorf("auth.sign_up.Handler.post - create user: %v", err))
		return
	}

	slog.Info(uid.String())

	jsonBytes, errCode := h.generateAccessToken()
	if errCode != 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte{byte(errCode)})
		return
	}

	request, err := json.Marshal(map[string]string{
		"token": string(jsonBytes),
		"url":   "/account",
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	cookie := http.Cookie{
		Name:    "session_token",
		Value:   string(jsonBytes),
		Expires: time.Now().Add(120 * time.Second),
	}

	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(request)
}

func (h *Handler) generateAccessToken() ([]byte, int) {
	key := []byte(config.GetConfig().JWT.Secret)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return nil, 418
	}
	builder := jwt.NewBuilder(signer)

	//TODO insert real user data in claims
	claims := linguaJWT.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "uuid_here",
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
		Email: "email@will.be.here",
	}
	token, err := builder.Build(claims)
	if err != nil {
		slog.Error(err.Error())
		return nil, http.StatusUnauthorized
	}

	slog.Info("create refresh token")
	refreshTokenUuid := uuid.New()
	/*err = h.RTCache.Set([]byte(refreshTokenUuid.String()), []byte(claims.ID), 0)
	if err != nil {
		h.Logger.Error(err)
		return nil, http.StatusInternalServerError
	}*/
	jsonBytes, err := json.Marshal(map[string]string{
		"token":         token.String(),
		"refresh_token": refreshTokenUuid.String(),
	})
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return jsonBytes, 0
}

func (h *Handler) validateEmail(ctx context.Context, email string) error {
	if !tools.IsEmailValid(email) {
		return entity.ErrEmailNotCorrect
	}

	uid, err := h.userSvc.GetIDByEmail(ctx, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if uid == uuid.Nil && err == nil {
		return entity.ErrItIsAdmin
	} else if uid != uuid.Nil {
		return entity.ErrEmailBusy
	}

	return nil
}

func (h *Handler) validateUsername(ctx context.Context, username string) error {
	if len(username) <= entity.UsernameLen {
		return entity.ErrUsernameLen
	}

	if strings.Contains(strings.ToLower(username), entity.UsernameAdmin) {
		return entity.ErrUsernameAdmin
	}

	uid, err := h.userSvc.GetIDByName(ctx, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if uid == uuid.Nil && err == nil {
		return entity.ErrItIsAdmin
	} else if uid != uuid.Nil {
		return entity.ErrUsernameBusy
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < entity.PasswordLen {
		return entity.ErrPasswordLen
	}

	if !tools.IsPasswordValid(password) {
		return entity.ErrPasswordDifficult
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
