package sign_up

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"lingua-evo/internal/config"
	"lingua-evo/internal/delivery/handlers/sign_up/entity"
	"lingua-evo/internal/delivery/repository"
	"lingua-evo/internal/service"

	"lingua-evo/pkg/logging"
	linguaJWT "lingua-evo/pkg/middleware/jwt"
	"lingua-evo/pkg/tools"

	"github.com/cristalhq/jwt/v3"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

const (
	signUpURL = "/signup"

	signupPage = "./web/static/sign_up/signup.html"
)

type Handler struct {
	logger *logging.Logger
	lingua *service.Lingua
	//RTCache cache.Repository
}

func Create(log *logging.Logger, ling *service.Lingua, r *httprouter.Router) {
	handler := newHandler(log, ling)
	handler.register(r)
}

func newHandler(logger *logging.Logger, lingua *service.Lingua) *Handler {
	return &Handler{
		logger: logger,
		lingua: lingua,
	}
}

func (h *Handler) register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, signUpURL, h.get)
	router.HandlerFunc(http.MethodPost, signUpURL, h.post)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile(signupPage)
	if err != nil {
		h.logger.Errorf("auth.getSignup.ReadFile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(file))
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if !tools.IsEmailValid(user.Email) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email is not correct"))
		return
	}

	if err := h.validateUsername(r.Context(), user.Username); err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return
	}
	if err := validatePassword(user.Password); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	hashPassword, err := hashPassword(user.Password)
	if err != nil {
		return
	}

	uid, err := h.lingua.AddUser(r.Context(), &repository.User{
		Username:     user.Username,
		PasswordHash: hashPassword,
		Email:        user.Email})
	if err != nil {
		return
	}

	h.logger.Println(uid)

	jsonBytes, errCode := h.generateAccessToken()
	if errCode != 0 {
		w.WriteHeader(errCode)
		return
	}

	request, err := json.Marshal(map[string]string{
		"token": string(jsonBytes),
		"url":   "/account",
	})
	if err != nil {
		w.WriteHeader(errCode)
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
	w.Write(request)
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
		h.logger.Error(err)
		return nil, http.StatusUnauthorized
	}

	h.logger.Info("create refresh token")
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

func (h *Handler) validateUsername(ctx context.Context, username string) error {
	if len(username) < 4 {
		return fmt.Errorf("username must be more 3 characters")
	}

	if strings.Contains(username, "admin") {
		return fmt.Errorf("username can not contains admin")
	}

	uid, err := h.lingua.FindUser(ctx, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if uid == uuid.Nil && err == nil {
		return fmt.Errorf("it is admin")
	} else if uid != uuid.Nil {
		return fmt.Errorf("this username is busy")
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password must be more 5 characters")
	}

	if !tools.IsPasswordValid(password) {
		return fmt.Errorf("password must be more difficult")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
