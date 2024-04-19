package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/token"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	"github.com/av-ugolkov/lingua-evo/internal/services/auth/model"
	entityUser "github.com/av-ugolkov/lingua-evo/internal/services/user"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
)

var (
	errNotEqualFingerprints = errors.New("new fingerprint is not equal old fingerprint")
)

type (
	sessionRepo interface {
		SetSession(ctx context.Context, tokenID uuid.UUID, s *Session, expiration time.Duration) error
		GetSession(ctx context.Context, refreshTokenID uuid.UUID) (*Session, error)
		GetCountSession(ctx context.Context, userID uuid.UUID) (int64, error)
		DeleteSession(ctx context.Context, session uuid.UUID) error
		SetAccountCode(ctx context.Context, email string, code int, expiration time.Duration) error
		GetAccountCode(ctx context.Context, email string) (int, error)
	}

	userSvc interface {
		AddUser(ctx context.Context, user *entityUser.User) (uuid.UUID, error)
		GetUser(ctx context.Context, login string) (*entityUser.User, error)
		GetUserByID(ctx context.Context, uid uuid.UUID) (*entityUser.User, error)
		GetUserByEmail(ctx context.Context, email string) (*entityUser.User, error)
		GetUserByName(ctx context.Context, name string) (*entityUser.User, error)
	}

	Service struct {
		repo    sessionRepo
		userSvc userSvc
	}
)

func NewService(repo sessionRepo, userSvc userSvc) *Service {
	return &Service{
		repo:    repo,
		userSvc: userSvc,
	}
}

func (s *Service) SignUp(ctx context.Context, data model.CreateUserRq, role runtime.Role) (uuid.UUID, error) {
	if err := s.validateEmail(ctx, data.Email); err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - validateEmail: %v", err)
	}

	code, err := s.repo.GetAccountCode(ctx, data.Email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - GetAccountCode: %v", err)
	}

	if code != data.Code {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp: code mismatch")
	}

	if err := s.validateUsername(ctx, data.Username); err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - validateUsername: %v", err)
	}

	if err := validatePassword(data.Password); err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - validatePassword: %v", err)
	}

	hashPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - hashPassword: %v", err)
	}

	user := &entityUser.User{
		ID:           uuid.New(),
		Name:         data.Username,
		PasswordHash: hashPassword,
		Email:        data.Email,
		Role:         role,
		CreatedAt:    time.Now().UTC(),
		LastVisitAt:  time.Now().UTC(),
	}

	uid, err := s.userSvc.AddUser(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp: %w", err)
	}

	return uid, nil
}

func (s *Service) SignIn(ctx context.Context, user, password, fingerprint string) (*Tokens, error) {
	u, err := s.userSvc.GetUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn - getUser: %v", err)
	}
	if err := utils.CheckPasswordHash(password, u.PasswordHash); err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn - incorrect password: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	now := time.Now().UTC()
	session := &Session{
		UserID:      u.ID,
		Fingerprint: fingerprint,
		CreatedAt:   now,
	}

	refreshTokenID := uuid.New()
	err = s.addRefreshSession(ctx, refreshTokenID, session)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn - addRefreshSession: %v", err)
	}

	claims := &Claims{
		ID:        refreshTokenID,
		UserID:    u.ID,
		ExpiresAt: now.Add(duration),
	}

	accessToken, err := token.NewJWTToken(u.ID, claims.ID, claims.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.SignIn - jwt.NewToken: %v", err)
	}

	tokens := &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenID,
	}

	return tokens, nil
}

// RefreshSessionToken - the method is called from the client
func (s *Service) RefreshSessionToken(ctx context.Context, refreshTokenID uuid.UUID, fingerprint string) (*Tokens, error) {
	oldRefreshSession, err := s.repo.GetSession(ctx, refreshTokenID)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken - get session: %v", err)
	}

	if oldRefreshSession.Fingerprint != fingerprint {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken: %w", errNotEqualFingerprints)
	}

	newSession := &Session{
		UserID:      oldRefreshSession.UserID,
		Fingerprint: oldRefreshSession.Fingerprint,
		CreatedAt:   time.Now().UTC(),
	}

	tokenID := uuid.New()
	err = s.addRefreshSession(ctx, tokenID, newSession)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken - addRefreshSession: %v", err)
	}

	err = s.repo.DeleteSession(ctx, refreshTokenID)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.RefreshSessionToken - delete session: %v", err)
	}

	additionalTime := config.GetConfig().JWT.ExpireAccess
	duration := time.Duration(additionalTime) * time.Second
	claims := &Claims{
		ID:        tokenID,
		UserID:    oldRefreshSession.UserID,
		ExpiresAt: time.Now().UTC().Add(duration),
	}
	u, err := s.userSvc.GetUserByID(ctx, oldRefreshSession.UserID)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.CreateSession - get user by ID: %v", err)
	}

	accessToken, err := token.NewJWTToken(u.ID, claims.ID, claims.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("auth.Service.CreateSession - create access token: %v", err)
	}

	tokens := &Tokens{
		AccessToken:  accessToken,
		RefreshToken: tokenID,
	}

	return tokens, nil
}

func (s *Service) SignOut(ctx context.Context, refreshToken uuid.UUID, fingerprint string) error {
	oldRefreshSession, err := s.repo.GetSession(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("auth.Service.Logout - GetSession: %v", err)
	}

	if oldRefreshSession.Fingerprint != fingerprint {
		return fmt.Errorf("auth.Service.Logout: %w", errNotEqualFingerprints)
	}

	err = s.repo.DeleteSession(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("auth.Service.Logout - DeleteSession: %v", err)
	}

	return nil
}

func (s *Service) CreateCode(ctx context.Context, email string) error {
	if err := s.validateEmail(ctx, email); err != nil {
		return fmt.Errorf("user.delivery.Handler.createAccount - validateEmail: %v", err)
	}

	from := "amak07@yandex.ru"     //"makedonskiy07@gmail.com"
	password := "sbdkuhybqeesvqvg" //"zseo awvn mmjd sklz"

	toEmailAddress := email
	to := []string{toEmailAddress}

	host := "smtp.yandex.ru" //"smtp.gmail.com"
	port := "587"            //"587"
	address := host + ":" + port

	creatingCode := rand.Intn(999999-100000) + 100000

	subject := "Subject: Create account on Lingua Evo\r\n\r\n"
	body := fmt.Sprintf("Аccount creation code: %d", creatingCode)
	message := []byte(subject + body)

	authEmail := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, authEmail, from, to, message)
	if err != nil {
		return fmt.Errorf("auth.Service.CreateCode - send mail: %v", err)
	}

	err = s.repo.SetAccountCode(ctx, email, creatingCode, time.Duration(10)*time.Minute)
	if err != nil {
		return fmt.Errorf("auth.Service.CreateCode: %w", err)
	}

	return nil
}

func (s *Service) addRefreshSession(ctx context.Context, tokenID uuid.UUID, refreshSession *Session) error {
	expiration := time.Duration(config.GetConfig().JWT.ExpireRefresh) * time.Second
	err := s.repo.SetSession(ctx, tokenID, refreshSession, expiration)
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
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if userData.ID == uuid.Nil && err == nil {
		return ErrItIsAdmin
	} else if userData.ID != uuid.Nil {
		return ErrEmailBusy
	}

	return nil
}

func (s *Service) validateUsername(ctx context.Context, username string) error {
	if len(username) <= UsernameLen {
		return ErrUsernameLen
	}

	userData, err := s.userSvc.GetUserByName(ctx, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	} else if errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if userData.ID == uuid.Nil && err == nil {
		return ErrItIsAdmin
	} else if userData.ID != uuid.Nil {
		return ErrUsernameBusy
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
