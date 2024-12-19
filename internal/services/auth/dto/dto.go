package dto

import (
	entity "github.com/av-ugolkov/lingua-evo/internal/services/auth"

	"github.com/google/uuid"
)

type CreateSessionRq struct {
	User        string `json:"user"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type CreateSessionRs struct {
	AccessToken string `json:"access_token"`
}

func CreateSessionToDTO(session *entity.Tokens) *CreateSessionRs {
	return &CreateSessionRs{
		AccessToken: session.AccessToken,
	}
}

type CreateUserRq struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Code     int    `json:"code" binding:"required"`
}

type CreateUserRs struct {
	UserID uuid.UUID `json:"user_id"`
}

type CreateCodeRq struct {
	Email string `json:"email"`
}

type GoogleAuthCode struct {
	Code     string   `json:"code"`
	State    string   `json:"state"`
	Scope    []string `json:"scope"`
	Authuser int      `json:"authuser"`
	Prompt   string   `json:"prompt"`
}
