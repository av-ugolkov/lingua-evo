package sign_in

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"lingua-evo/pkg/tools"
)

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) postSignIn(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//TODO client to UserService and get user by username and password
	//for now stub check
	//if u.Username != "me" || u.Password != "pass" {
	//	w.WriteHeader(http.StatusNotFound)
	//	return
	//}

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
