package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	jsoniter "github.com/json-iterator/go"
)

type (
	userRepo interface {
		AddUser(ctx context.Context, u *User) (uuid.UUID, error)
		EditUser(ctx context.Context, u *User) error
		GetUserByID(ctx context.Context, uid uuid.UUID) (*User, error)
		GetUserByName(ctx context.Context, name string) (*User, error)
		GetUserByEmail(ctx context.Context, email string) (*User, error)
		GetUserByToken(ctx context.Context, token uuid.UUID) (*User, error)
		RemoveUser(ctx context.Context, u *User) error
		GetUserData(ctx context.Context, uid uuid.UUID) (*Data, error)
		GetUserSubscriptions(ctx context.Context, uid uuid.UUID) ([]Subscriptions, error)
		GetUsers(ctx context.Context, page, perPage, sort, order int, search string) ([]UserData, int, error)
		UpdateLastVisited(ctx context.Context, uid uuid.UUID) error
	}

	redis interface {
		Get(ctx context.Context, key string) (string, error)
		GetAccountCode(ctx context.Context, email string) (int, error)
	}
)

type Service struct {
	repo  userRepo
	redis redis
}

func NewService(repo userRepo, redis redis) *Service {
	return &Service{
		repo:  repo,
		redis: redis,
	}
}

func (s *Service) SignUp(ctx context.Context, userCreate UserCreate) (uuid.UUID, error) {
	if err := s.validateEmail(ctx, userCreate.Email); err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - validateEmail: %v", err)
	}

	code, err := s.redis.GetAccountCode(ctx, userCreate.Email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - GetAccountCode: %v", err)
	}

	if code != userCreate.Code {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp: code mismatch")
	}

	if err := s.validateUsername(ctx, userCreate.Name); err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - validateUsername: %v", err)
	}

	if err := validatePassword(userCreate.Password); err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - validatePassword: %v", err)
	}

	hashPassword, err := utils.HashPassword(userCreate.Password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - hashPassword: %v", err)
	}

	user := &User{
		ID:           userCreate.ID,
		Name:         userCreate.Name,
		PasswordHash: hashPassword,
		Email:        userCreate.Email,
		Role:         userCreate.Role,
		CreatedAt:    time.Now().UTC(),
		LastVisitAt:  time.Now().UTC(),
	}

	uid, err := s.repo.AddUser(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("user.Service.AddUser: %w", err)
	}

	return uid, nil
}

func (s *Service) EditUser(ctx context.Context, user *User) error {
	return nil
}

func (s *Service) GetUser(ctx context.Context, login string) (*User, error) {
	user, err := s.repo.GetUserByName(ctx, login)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("user.Service.GetUser - by name: %w", err)
	} else if errors.Is(err, pgx.ErrNoRows) {
		user, err = s.repo.GetUserByEmail(ctx, login)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user.Service.GetUser - by email: %w", err)
		}
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("user.Service.GetUser - by [%s]: %w", login, ErrNotFoundUser)
	}

	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, uid uuid.UUID) (*User, error) {
	user, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("user.Service.GetUserByID: %w", err)
	}

	return user, nil
}

func (s *Service) GetUserByName(ctx context.Context, name string) (*User, error) {
	return s.repo.GetUserByName(ctx, name)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *Service) GetUserByRefreshToken(ctx context.Context, token uuid.UUID) (*User, error) {
	sessionJson, err := s.redis.Get(ctx, token.String())
	if err != nil {
		return nil, fmt.Errorf("user.Service.GetUserByRefreshToken: %w", err)
	}

	var session Session

	err = jsoniter.Unmarshal([]byte(sessionJson), &session)
	if err != nil {
		return nil, fmt.Errorf("user.Service.GetUserByRefreshToken: %w", err)
	}

	return s.repo.GetUserByToken(ctx, session.UserID)
}

func (s *Service) RemoveUser(ctx context.Context, user *User) error {
	return nil
}

func (s *Service) UserCountWord(ctx context.Context, userID uuid.UUID) (int, error) {
	data, err := s.repo.GetUserData(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("user.Service.UserCountWord - get user data: %w", err)
	}

	subscriptions, err := s.repo.GetUserSubscriptions(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("user.Service.UserCountWord - get user subscriptions: %w", err)
	}

	maxWords := data.MaxCountWords
	for _, sub := range subscriptions {
		maxWords += sub.CountWord
	}

	return maxWords, nil
}

func (s *Service) GetUsers(ctx context.Context, page, perPage, sort, order int, search string) ([]UserData, int, error) {
	users, countUsers, err := s.repo.GetUsers(ctx, page, perPage, sort, order, search)
	if err != nil {
		return nil, 0, fmt.Errorf("user.Service.GetUsers: %w", err)
	}

	return users, countUsers, nil
}

func (s *Service) UpdateLastVisited(ctx context.Context, uid uuid.UUID) error {
	err := s.repo.UpdateLastVisited(ctx, uid)
	if err != nil {
		return fmt.Errorf("user.Service.UpdateLastVisited: %w", err)
	}

	return nil
}

func (s *Service) validateEmail(ctx context.Context, email string) error {
	if !utils.IsEmailValid(email) {
		return ErrEmailNotCorrect
	}

	userData, err := s.GetUserByEmail(ctx, email)
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
	if len(username) <= UsernameLen {
		return ErrUsernameLen
	}

	userData, err := s.GetUserByName(ctx, username)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	} else if errors.Is(err, pgx.ErrNoRows) {
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
