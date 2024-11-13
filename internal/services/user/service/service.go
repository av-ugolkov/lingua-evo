package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	jsoniter "github.com/json-iterator/go"
)

type (
	userRepo interface {
		AddUser(ctx context.Context, u *entity.User, pswHash string) (uuid.UUID, error)
		GetUserByID(ctx context.Context, uid uuid.UUID) (*entity.User, error)
		GetUserByNickname(ctx context.Context, name string) (*entity.User, error)
		GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
		RemoveUser(ctx context.Context, uid uuid.UUID) error
		GetUserData(ctx context.Context, uid uuid.UUID) (*entity.UserData, error)
		GetUserMaxCountWords(ctx context.Context, uid uuid.UUID) (int, error)
		AddUserData(ctx context.Context, data entity.UserData) error
		AddUserNewsletters(ctx context.Context, data entity.UserNewsletters) error
		GetUserSubscriptions(ctx context.Context, uid uuid.UUID) ([]entity.Subscriptions, error)
		GetUsers(ctx context.Context, page, perPage, sort, order int, search string) ([]entity.User, int, error)
		UpdateVisitedAt(ctx context.Context, uid uuid.UUID) error

		repoSettings
	}

	redis interface {
		Get(ctx context.Context, key string) (string, error)
		GetAccountCode(ctx context.Context, email string) (int, error)
		SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error)
	}
)

//go:generate mockery --inpackage --outpkg service --testonly --name "userRepo|redis|emailSvc"

const (
	DEFAULT_COUNT_WORDS = 300
)

type Service struct {
	tr       *transactor.Transactor
	repo     userRepo
	redis    redis
	emailSvc emailSvc
}

func NewService(tr *transactor.Transactor, repo userRepo, redis redis, emailSvc emailSvc) *Service {
	return &Service{
		tr:       tr,
		repo:     repo,
		redis:    redis,
		emailSvc: emailSvc,
	}
}

func (s *Service) SignUp(ctx context.Context, userCreate entity.UserCreate) (uuid.UUID, error) {
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

	uid, err := s.AddUser(ctx, userCreate)
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth.Service.SignUp - AddUser: %v", err)
	}

	return uid, nil
}

func (s *Service) AddUser(ctx context.Context, usr entity.UserCreate) (_ uuid.UUID, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("auth.Service.AddUser: %v", err)
		}
	}()

	if err := s.validateUsername(ctx, usr.Nickname); err != nil {
		return uuid.Nil, err
	}

	if err := validatePassword(usr.Password); err != nil {
		return uuid.Nil, err
	}

	hashPassword, err := utils.HashPassword(usr.Password)
	if err != nil {
		return uuid.Nil, err
	}

	user := &entity.User{
		Nickname:  usr.Nickname,
		Email:     usr.Email,
		Role:      usr.Role,
		CreatedAt: time.Now().UTC(),
		VisitedAt: time.Now().UTC(),
	}

	uid := uuid.Nil
	err = s.tr.CreateTransaction(ctx, func(ctx context.Context) error {
		uid, err = s.repo.AddUser(ctx, user, hashPassword)
		if err != nil {
			return err
		}

		err = s.repo.AddUserData(ctx, entity.UserData{UID: uid})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func (s *Service) GetUser(ctx context.Context, login string) (*entity.User, error) {
	user, err := s.repo.GetUserByNickname(ctx, login)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("user.Service.GetUser - by name: %w", err)
	} else if errors.Is(err, pgx.ErrNoRows) {
		user, err = s.repo.GetUserByEmail(ctx, login)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user.Service.GetUser - by email: %w", err)
		}
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("user.Service.GetUser - by [%s]: %w", login, entity.ErrNotFoundUser)
	}

	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, uid uuid.UUID) (*entity.User, error) {
	user, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("user.Service.GetUserByID: %w", err)
	}

	return user, nil
}

func (s *Service) GetUserByNickname(ctx context.Context, name string) (*entity.User, error) {
	usr, err := s.repo.GetUserByNickname(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("user.Service.GetUserByNickname: %w", err)
	}

	return usr, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *Service) GetUserByRefreshToken(ctx context.Context, token uuid.UUID) (*entity.User, error) {
	sessionJson, err := s.redis.Get(ctx, token.String())
	if err != nil {
		return nil, fmt.Errorf("user.Service.GetUserByRefreshToken: %w", err)
	}

	var session entity.Session

	err = jsoniter.Unmarshal([]byte(sessionJson), &session)
	if err != nil {
		return nil, fmt.Errorf("user.Service.GetUserByRefreshToken: %w", err)
	}

	return s.repo.GetUserByID(ctx, session.UserID)
}

func (s *Service) RemoveUser(ctx context.Context, user *entity.User) error {
	return nil
}

func (s *Service) UserCountWord(ctx context.Context, uid uuid.UUID) (int, error) {
	maxCountWords, err := s.repo.GetUserMaxCountWords(ctx, uid)
	if err != nil {
		return 0, fmt.Errorf("user.Service.UserCountWord - get user data: %w", err)
	}

	subscriptions, err := s.repo.GetUserSubscriptions(ctx, uid)
	if err != nil {
		return 0, fmt.Errorf("user.Service.UserCountWord - get user subscriptions: %w", err)
	}

	for _, sub := range subscriptions {
		maxCountWords += sub.CountWord
	}

	return maxCountWords, nil
}

func (s *Service) GetUsers(ctx context.Context, uid uuid.UUID, page, perPage, sort, order int, search string) ([]entity.User, int, error) {
	users, countUsers, err := s.repo.GetUsers(ctx, page, perPage, sort, order, search)
	if err != nil {
		return nil, 0, fmt.Errorf("user.Service.GetUsers: %w", err)
	}

	return users, countUsers, nil
}

func (s *Service) UpdateVisitedAt(ctx context.Context, uid uuid.UUID) error {
	err := s.repo.UpdateVisitedAt(ctx, uid)
	if err != nil {
		return fmt.Errorf("user.Service.UpdateVisitedAt: %w", err)
	}

	return nil
}

func (s *Service) validateEmail(ctx context.Context, email string) error {
	if !utils.IsEmailValid(email) {
		return entity.ErrEmailNotCorrect
	}

	userData, err := s.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	} else if errors.Is(err, pgx.ErrNoRows) {
		return nil
	} else if userData != nil && userData.ID == uuid.Nil && err == nil {
		return entity.ErrItIsAdmin
	} else if userData != nil && userData.ID != uuid.Nil {
		return entity.ErrEmailBusy
	}

	return nil
}

func (s *Service) validateUsername(ctx context.Context, username string) error {
	if len(username) <= entity.MinUsernameLen {
		return entity.ErrUsernameLen
	}

	userData, err := s.GetUserByNickname(ctx, username)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	} else if errors.Is(err, pgx.ErrNoRows) {
		return nil
	} else if userData.ID == uuid.Nil && err == nil {
		return entity.ErrItIsAdmin
	} else if userData.ID != uuid.Nil {
		return entity.ErrUsernameBusy
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < entity.MinPasswordLen {
		return entity.ErrPasswordLen
	}

	if !utils.IsPasswordValid(password) {
		return entity.ErrPasswordDifficult
	}

	return nil
}
