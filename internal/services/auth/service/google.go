package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/google"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	jwtToken "github.com/av-ugolkov/lingua-evo/internal/pkg/token"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/auth"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type GoogleToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) AuthByGoogle(ctx context.Context, code, fingerprint string) (_ *entity.Tokens, err error) {
	token, err := google.GetTokenByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.AuthByGoogle: %w", err)
	}

	defer func() {
		if err != nil {
			err = google.RevokeGoogleToken(token.AccessToken)
			if err != nil {
				err = fmt.Errorf("auth.Service.AuthByGoogle: %w", err)
			}
		}
	}()

	profile, err := google.GetProfile(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.AuthByGoogle: %w", err)
	}

	var session *entity.Session
	usr, err := s.userSvc.GetUserByGoogleID(ctx, profile.ID)
	if err != nil {
		slog.Warn(fmt.Sprintf("auth.Service.AuthByGoogle: %v", err.Error()))
		if token.RefreshToken == runtime.EmptyString {
			return nil, fmt.Errorf("auth.Service.AuthByGoogle: not fount user and refresh token is empty")
		}
		uid, err := s.userSvc.AddGoogleUser(ctx, entityUser.GoogleUser{
			User: entityUser.User{
				Nickname: runtime.GenerateNickname(),
				Email:    profile.Email,
				Role:     runtime.User,
			},
			GoogleID: profile.ID,
		})
		if err != nil {
			return nil, msgerr.New(fmt.Errorf("auth.Service.AuthByGoogle: %w", err),
				entity.ErrMsgUserExists)
		}

		session = &entity.Session{
			UserID:       uid,
			RefreshToken: token.RefreshToken,
			TypeToken:    entity.Google,
			ExpiresAt:    token.Expiry,
		}
	} else {
		session, err = s.repo.GetSession(ctx, fmt.Sprintf("%s:%s:%s", usr.ID, fingerprint, RedisRefreshToken))
		switch {
		case errors.Is(err, redis.Nil):
			session = &entity.Session{
				UserID:       usr.ID,
				RefreshToken: token.RefreshToken,
				TypeToken:    entity.Google,
				ExpiresAt:    token.Expiry,
			}
		case err != nil:
			return nil, fmt.Errorf("auth.Service.AuthByGoogle: %w", err)
		}
		if token.RefreshToken != runtime.EmptyString {
			session.RefreshToken = token.RefreshToken
		}
		if session.RefreshToken == runtime.EmptyString {
			return nil, fmt.Errorf("auth.Service.AuthByGoogle: not fount refresh token")
		}
		session.ExpiresAt = token.Expiry
	}

	err = s.addRefreshSession(ctx, fmt.Sprintf("%s:%s:%s", session.UserID, fingerprint, RedisRefreshToken), session)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.AuthByGoogle: %w", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	now := time.Now().UTC()

	accessToken, err := jwtToken.NewJWTToken(usr.ID, token.RefreshToken, now.Add(duration))
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn: %v", err)
	}

	err = s.userSvc.UpdateVisitedAt(ctx, usr.ID)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn: %v", err)
	}

	tokens := &entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: session.RefreshToken,
	}

	return tokens, nil
}

func (s *Service) refreshGoogleToken(ctx context.Context, uid uuid.UUID, refreshToken string) (*entity.Tokens, error) {
	_, err := google.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.refreshGoogleToken: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	accessToken, err := jwtToken.NewJWTToken(uid, refreshToken, time.Now().UTC().Add(duration))
	if err != nil {
		return nil, fmt.Errorf("auth.Service.refreshGoogleToken: %v", err)
	}

	tokens := &entity.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return tokens, nil
}
