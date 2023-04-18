package auth

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewUser struct {
	User  User   `json:"user"`
	Email string `json:"email*"`
}

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	var nu NewUser
	if err := json.NewDecoder(r.Body).Decode(&nu); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//TODO validate username and password
	//TODO create user using UserService
	jsonBytes, errCode := h.generateAccessToken()
	if errCode != 0 {
		w.WriteHeader(errCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(jsonBytes)
}

func (h *Handler) postAuth(w http.ResponseWriter, r *http.Request) {
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
