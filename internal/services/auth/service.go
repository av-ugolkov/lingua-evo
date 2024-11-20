package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/google"
	msgerr "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/token"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"golang.org/x/oauth2"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	RedisCreateUser = "create_user"
)

type (
	sessionRepo interface {
		SetSession(ctx context.Context, key string, s Session, ttl time.Duration) error
		GetSession(ctx context.Context, key string) (Session, error)
		GetCountSession(ctx context.Context, userID uuid.UUID) (int64, error)
		DeleteSession(ctx context.Context, key string) error
		SetAccountCode(ctx context.Context, key string, code int, ttl time.Duration) error
		GetAccountCode(ctx context.Context, key string) (int, error)
	}

	userSvc interface {
		AddUser(ctx context.Context, userCreate entityUser.User, pswHash string) (uuid.UUID, error)
		AddGoogleUser(ctx context.Context, userCreate entityUser.GoogleUser) (uuid.UUID, error)
		GetUser(ctx context.Context, login string) (*entityUser.User, error)
		GetUserByNickname(ctx context.Context, nickname string) (*entityUser.User, error)
		GetPswHash(ctx context.Context, uid uuid.UUID) (string, error)
		GetUserByEmail(ctx context.Context, email string) (*entityUser.User, error)
		UpdateVisitedAt(ctx context.Context, uid uuid.UUID) error
		GetUserByGoogleID(ctx context.Context, googleID string) (*entityUser.User, error)
	}

	emailSvc interface {
		SendAuthCode(toEmail string, code int) error
	}
)

type Service struct {
	repo    sessionRepo
	userSvc userSvc
	email   emailSvc
}

func NewService(repo sessionRepo, userSvc userSvc, email emailSvc) *Service {
	return &Service{
		repo:    repo,
		userSvc: userSvc,
		email:   email,
	}
}

func (s *Service) SignIn(ctx context.Context, user, password, fingerprint string, refreshTokenID uuid.UUID) (*Tokens, error) {
	u, err := s.userSvc.GetUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn: %w", err)
	}

	pswHash, err := s.userSvc.GetPswHash(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn: %w", err)
	}

	if err := utils.CheckPasswordHash(password, pswHash); err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn - [%w]: %v", ErrWrongPassword, err)
	}

	err = s.userSvc.UpdateVisitedAt(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	now := time.Now().UTC()
	session := Session(u.ID)

	err = s.addRefreshSession(ctx, fmt.Sprintf("%s:%s", fingerprint, refreshTokenID), session)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn: %v", err)
	}

	accessToken, err := token.NewJWTToken(u.ID, refreshTokenID, now.Add(duration))
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn: %v", err)
	}

	tokens := &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenID.String(),
	}

	return tokens, nil
}

func (s *Service) SignUp(ctx context.Context, usr User, fingerprint string) (uuid.UUID, error) {
	if err := s.validateEmail(ctx, usr.Email); err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp: %v", err)
	}

	code, err := s.repo.GetAccountCode(ctx, fmt.Sprintf("%s:%s:%s", fingerprint, usr.Email, RedisCreateUser))
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp: %v", err)
	}

	if code != usr.Code {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp: code mismatch")
	}

	if err := s.validateUsername(ctx, usr.Nickname); err != nil {
		return uuid.Nil, err
	}

	if err := validatePassword(usr.Password); err != nil {
		return uuid.Nil, err
	}

	pswHash, err := utils.HashPassword(usr.Password)
	if err != nil {
		return uuid.Nil, err
	}

	uid, err := s.userSvc.AddUser(ctx, entityUser.User{
		Nickname: usr.Nickname,
		Email:    usr.Email,
		Role:     usr.Role,
	}, pswHash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - AddUser: %v", err)
	}

	return uid, nil
}

// RefreshSessionToken - the method is called from the client
func (s *Service) RefreshSessionToken(ctx context.Context, newTokenID, oldTokenID uuid.UUID, fingerprint string) (*Tokens, error) {
	oldRefreshSession, err := s.repo.GetSession(ctx, fmt.Sprintf("%s:%s", fingerprint, oldTokenID))
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", err)
	}

	err = s.addRefreshSession(ctx, fmt.Sprintf("%s:%s", fingerprint, newTokenID), oldRefreshSession)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", err)
	}

	err = s.repo.DeleteSession(ctx, fmt.Sprintf("%s:%s", fingerprint, oldTokenID))
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", err)
	}

	err = s.userSvc.UpdateVisitedAt(ctx, uuid.UUID(oldRefreshSession))
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second

	accessToken, err := token.NewJWTToken(uuid.UUID(oldRefreshSession), newTokenID, time.Now().UTC().Add(duration))
	if err != nil {
		return nil, fmt.Errorf("auth.Service.CreateSession: %v", err)
	}

	tokens := &Tokens{
		AccessToken:  accessToken,
		RefreshToken: newTokenID.String(),
	}

	return tokens, nil
}

func (s *Service) SignOut(ctx context.Context, refreshToken uuid.UUID, fingerprint string) error {
	err := s.repo.DeleteSession(ctx, fmt.Sprintf("%s:%s", fingerprint, refreshToken))
	if err != nil {
		return fmt.Errorf("auth.Service.SignOut: %v", err)
	}

	return nil
}

func (s *Service) CreateCode(ctx context.Context, email string, fingerprint string) error {
	err := s.validateEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("auth.Service.CreateCode: %v", err)
	}

	creatingCode := utils.GenerateCode()

	err = s.repo.SetAccountCode(ctx, fmt.Sprintf("%s:%s:%s", fingerprint, email, RedisCreateUser), creatingCode, time.Duration(5)*time.Minute)
	if err != nil {
		return fmt.Errorf("auth.Service.CreateCode: %w", err)
	}

	err = s.email.SendAuthCode(email, creatingCode)
	if err != nil {
		return fmt.Errorf("auth.Service.CreateCode: %v", err)
	}

	return nil
}

func (s *Service) GoogleAuthUrl() string {
	return google.GetAuthUrl()
}

func (s *Service) AuthByGoogle(ctx context.Context, code, fingerprint string) (*oauth2.Token, error) {
	token, err := google.GetTokenByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.AuthByGoogle: %w", err)
	}

	/*
		client := googleOauthConfig.Client(context.Background(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var userInfo map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
			http.Error(w, "Failed to parse user info: "+err.Error(), http.StatusInternalServerError)
			return
		}
	*/

	profile, err := google.GetProfile(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.AuthByGoogle: %w", err)
	}
	var session Session
	usr, err := s.userSvc.GetUserByGoogleID(ctx, profile.ID)
	if err != nil {
		slog.Warn(fmt.Sprintf("auth.Service.AuthByGoogle: %v", err.Error()))

		uid, err := s.userSvc.AddGoogleUser(ctx, entityUser.GoogleUser{
			User: entityUser.User{
				Nickname: strings.Split(profile.Email, "@")[0],
				Email:    profile.Email,
				Role:     runtime.User,
			},
			GoogleID: profile.ID,
		})
		if err != nil {
			return nil,
				msgerr.New(fmt.Errorf("auth.Service.AuthByGoogle: %w", err),
					"The user exists with the same email")
		}

		session = Session(uid)
	} else {
		session = Session(usr.ID)
	}

	err = s.addRefreshSession(ctx, fmt.Sprintf("%s:%s", fingerprint, token.AccessToken), session)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.AuthByGoogle: %w", err)
	}

	return token, nil
}

func (s *Service) addRefreshSession(ctx context.Context, key string, refreshSession Session) error {
	ttl := time.Duration(config.GetConfig().JWT.ExpireRefresh) * time.Second
	err := s.repo.SetSession(ctx, key, refreshSession, ttl)
	if err != nil {
		return fmt.Errorf("auth.Service.addRefreshSession: %w", err)
	}
	return nil
}

func (s *Service) validateEmail(ctx context.Context, email string) error {
	if !utils.IsEmailValid(email) {
		return ErrEmailNotCorrect
	}

	userData, err := s.userSvc.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	} else if errors.Is(err, pgx.ErrNoRows) {
		return nil
	} else if userData != nil && userData.ID == uuid.Nil && err == nil {
		return ErrItIsAdmin
	} else if userData != nil && userData.ID != uuid.Nil {
		return ErrEmailBusy
	}

	return nil
}

func (s *Service) validateUsername(ctx context.Context, username string) error {
	if len(username) <= MinNicknameLen {
		return ErrNicknameLen
	}

	userData, err := s.userSvc.GetUserByNickname(ctx, username)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	} else if errors.Is(err, pgx.ErrNoRows) {
		return nil
	} else if userData.ID == uuid.Nil && err == nil {
		return ErrItIsAdmin
	} else if userData.ID != uuid.Nil {
		return ErrNicknameBusy
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < MinPasswordLen {
		return ErrPasswordLen
	}

	if !utils.IsPasswordValid(password) {
		return ErrPasswordDifficult
	}

	return nil
}
