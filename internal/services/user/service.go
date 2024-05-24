package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"

	"github.com/google/uuid"
)

type (
	userRepo interface {
		AddUser(ctx context.Context, u *User) (uuid.UUID, error)
		EditPassword(ctx context.Context, u *User) error
		EditEmail(ctx context.Context, u *User) error
		EditUsername(ctx context.Context, u *User) error
		GetUserByID(ctx context.Context, uid uuid.UUID) (*User, error)
		GetUserByName(ctx context.Context, name string) (*User, error)
		GetUserByEmail(ctx context.Context, email string) (*User, error)
		GetUserByToken(ctx context.Context, token uuid.UUID) (*User, error)
		RemoveUser(ctx context.Context, u *User) error
	}

	sessionRepo interface {
		Get(ctx context.Context, key string) (string, error)
		GetAccountCode(ctx context.Context, email string) (int, error)
	}

	Service struct {
		repo        userRepo
		sessionRepo sessionRepo
	}
)

func NewService(repo userRepo, sessionRepo sessionRepo) *Service {
	return &Service{
		repo:        repo,
		sessionRepo: sessionRepo,
	}
}

func (s *Service) SignUp(ctx context.Context, userData UserData) (uuid.UUID, error) {
	if err := s.validateEmail(ctx, userData.Email); err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - validateEmail: %v", err)
	}

	code, err := s.sessionRepo.GetAccountCode(ctx, userData.Email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - GetAccountCode: %v", err)
	}

	if code != userData.Code {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp: code mismatch")
	}

	if err := s.validateUsername(ctx, userData.Name); err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - validateUsername: %v", err)
	}

	if err := validatePassword(userData.Password); err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - validatePassword: %v", err)
	}

	hashPassword, err := utils.HashPassword(userData.Password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - hashPassword: %v", err)
	}

	user := &User{
		ID:           userData.ID,
		Name:         userData.Name,
		PasswordHash: hashPassword,
		Email:        userData.Email,
		Role:         userData.Role,
		CreatedAt:    time.Now().UTC(),
		LastVisitAt:  time.Now().UTC(),
	}

	uid, err := s.repo.AddUser(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("user.Service.AddUser: %w", err)
	}

	return uid, nil
}

func (s *Service) EditPassword(ctx context.Context, user UserPasword) error {
	if err := validatePassword(user.Password); err != nil {
		return fmt.Errorf("auth.Service.EditPassword - validatePassword: %v", err)
	}

	u, err := s.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("auth.Service.EditPassword - GetUserByID: %v", err)
	}
	code, err := s.sessionRepo.GetAccountCode(ctx, u.Email)
	if err != nil {
		return fmt.Errorf("auth.Service.EditPassword - GetAccountCode: %v", err)
	}

	if code != user.Code {
		return fmt.Errorf("auth.Service.EditPassword: code mismatch")
	}

	u.PasswordHash, err = utils.HashPassword(u.PasswordHash)
	if err != nil {
		return fmt.Errorf("auth.Service.EditPassword - HashPassword: %v", err)
	}

	err = s.repo.EditPassword(ctx, u)
	if err != nil {
		return fmt.Errorf("auth.Service.EditPassword - EditPassword: %v", err)
	}
	return nil
}

func (s *Service) EditEmail(ctx context.Context, editUser EditUserData) error {
	if err := s.validateEmail(ctx, editUser.Email); err != nil {
		return fmt.Errorf("auth.Service.EditEmail - validatePassword: %v", err)
	}

	u, err := s.repo.GetUserByID(ctx, editUser.ID)
	if err != nil {
		return fmt.Errorf("auth.Service.EditEmail - GetUserByID: %v", err)
	}

	if err := utils.CheckPasswordHash(editUser.Password, u.PasswordHash); err != nil {
		return fmt.Errorf("auth.Service.EditEmail - incorrect password: %v", err)
	}

	if u.Email == editUser.Email {
		return fmt.Errorf("auth.Service.EditEmail: new email is the same as old email")
	}
	u.Email = editUser.Email

	err = s.repo.EditEmail(ctx, u)
	if err != nil {
		return fmt.Errorf("auth.Service.EditPassword - EditPassword: %v", err)
	}
	return nil
}

func (s *Service) EditUsername(ctx context.Context, editUser EditUserData) error {
	u, err := s.repo.GetUserByID(ctx, editUser.ID)
	if err != nil {
		return fmt.Errorf("auth.Service.EditEmail - GetUserByID: %v", err)
	}

	if err := utils.CheckPasswordHash(editUser.Password, u.PasswordHash); err != nil {
		return fmt.Errorf("auth.Service.EditEmail - incorrect password: %v", err)
	}

	if u.Name == editUser.Username {
		return fmt.Errorf("auth.Service.EditEmail: new username is the same as old username")
	}
	u.Name = editUser.Username

	err = s.repo.EditEmail(ctx, u)
	if err != nil {
		return fmt.Errorf("auth.Service.EditPassword - EditPassword: %v", err)
	}
	return nil
}

func (s *Service) GetUser(ctx context.Context, login string) (*User, error) {
	user, err := s.repo.GetUserByName(ctx, login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user.Service.GetUser - by name: %w", err)
	} else if errors.Is(err, sql.ErrNoRows) {
		user, err = s.repo.GetUserByEmail(ctx, login)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user.Service.GetUser - by email: %w", err)
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user.Service.GetUser - by [%s]: %w", login, ErrNotFoundUser)
	}

	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, uid uuid.UUID) (*User, error) {
	return s.repo.GetUserByID(ctx, uid)
}

func (s *Service) GetUserByName(ctx context.Context, name string) (*User, error) {
	return s.repo.GetUserByName(ctx, name)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *Service) GetUserByRefreshToken(ctx context.Context, token uuid.UUID) (*User, error) {
	sessionJson, err := s.sessionRepo.Get(ctx, token.String())
	if err != nil {
		return nil, fmt.Errorf("user.Service.GetUserByRefreshToken: %w", err)
	}

	var session Session

	err = json.Unmarshal([]byte(sessionJson), &session)
	if err != nil {
		return nil, fmt.Errorf("user.Service.GetUserByRefreshToken: %w", err)
	}

	return s.repo.GetUserByToken(ctx, session.UserID)
}

func (s *Service) RemoveUser(ctx context.Context, user *User) error {
	return nil
}

func (s *Service) validateEmail(ctx context.Context, email string) error {
	if !utils.IsEmailValid(email) {
		return ErrEmailNotCorrect
	}

	userData, err := s.GetUserByEmail(ctx, email)
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

	userData, err := s.GetUserByName(ctx, username)
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
