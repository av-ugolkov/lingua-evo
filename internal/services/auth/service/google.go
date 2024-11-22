package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/google"
	msgerr "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/auth"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"golang.org/x/oauth2"
)

type GoogleToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (g *GoogleToken) ReadToken(token string) ([]byte, error) {
	return []byte(token), nil
}

func (g *GoogleToken) WriteToken(token string) error {
	return nil
}

func (s *Service) AuthByGoogle(ctx context.Context, code, fingerprint string) (*oauth2.Token, error) {
	token, err := google.GetTokenByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.AuthByGoogle: %w", err)
	}

	profile, err := google.GetProfile(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.AuthByGoogle: %w", err)
	}
	var session *entity.Session
	usr, err := s.userSvc.GetUserByGoogleID(ctx, profile.ID)
	if err != nil {
		slog.Warn(fmt.Sprintf("auth.Service.AuthByGoogle: %v", err.Error()))

		uid, err := s.userSvc.AddGoogleUser(ctx, entityUser.GoogleUser{
			User: entityUser.User{
				Nickname: runtime.GenerateNickname(),
				Email:    profile.Email,
				Role:     runtime.User,
			},
			GoogleID: profile.ID,
		})
		if err != nil {
			return nil,
				msgerr.New(fmt.Errorf("auth.Service.AuthByGoogle: %w", err),
					"The user exists with the same email or nickname")
		}

		session = &entity.Session{
			UserID:       uid,
			RefreshToken: token.RefreshToken,
			TypeToken:    entity.Google,
			ExpiresAt:    token.Expiry,
		}
	} else {
		session = &entity.Session{
			UserID:       usr.ID,
			RefreshToken: token.RefreshToken,
			TypeToken:    entity.Google,
			ExpiresAt:    token.Expiry,
		}
	}

	err = s.addRefreshSession(ctx, fmt.Sprintf("%s:%s:%s", session.UserID, fingerprint, RedisRefreshToken), session)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.AuthByGoogle: %w", err)
	}

	return token, nil
}
