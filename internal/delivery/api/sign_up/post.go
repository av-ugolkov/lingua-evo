package sign_up

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"lingua-evo/internal/delivery/repository"

	"lingua-evo/pkg/tools"
)

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) postSignUp(w http.ResponseWriter, r *http.Request) {
	var user User
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
	if err := h.validatePassword(user.Password); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	hashPassword, err := h.hashPassword(user.Password)
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

func (h *Handler) validatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password must be more 5 characters")
	}

	if !tools.IsPasswordValid(password) {
		return fmt.Errorf("password must be more difficult")
	}

	return nil
}

func (h *Handler) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
