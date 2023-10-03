package sign_up

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"lingua-evo/internal/config"
	"lingua-evo/internal/services/sign_up/entity"
	entityUserSvc "lingua-evo/internal/services/user/entity"
	userSvc "lingua-evo/internal/services/user/service"
	staticFiles "lingua-evo/static"

	linguaJWT "lingua-evo/pkg/middleware/jwt"
	"lingua-evo/pkg/tools"

	"github.com/cristalhq/jwt/v3"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const (
	signUpURL = "/signup"

	signupPage = "web/sign_up/signup.html"
)

type Handler struct {
	usersSvc *userSvc.UserSvc
}

func Create(r *mux.Router, srcUser *userSvc.UserSvc) {
	handler := newHandler(srcUser)
	handler.register(r)
}

func newHandler(usersSvc *userSvc.UserSvc) *Handler {
	return &Handler{
		usersSvc: usersSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(signUpURL, h.get).Methods(http.MethodGet)
	r.HandleFunc(signUpURL, h.post).Methods(http.MethodPost)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	file, err := staticFiles.OpenFile(signupPage)
	if err != nil {
		slog.Error("sign_up.get.OpenFile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(file))
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	//TODO получение данных из формы плохое решение
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if err := h.validateEmail(r.Context(), email); err != nil {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	if err := h.validateUsername(r.Context(), username); err != nil {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	if err := validatePassword(password); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	hashPassword, err := hashPassword(password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	uid, err := h.usersSvc.AddUser(r.Context(), &entityUserSvc.User{
		Username:     username,
		PasswordHash: hashPassword,
		Email:        email})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
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

	uid, err := h.usersSvc.FindEmail(ctx, email)
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

	uid, err := h.usersSvc.FindUser(ctx, username)
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
