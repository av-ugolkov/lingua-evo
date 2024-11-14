package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/user"
	user "github.com/av-ugolkov/lingua-evo/internal/services/user/service"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	paramsPage    string = "page"
	paramsPerPage string = "per_page"
	paramsSearch  string = "search"
	paramsSort    string = "sort"
	paramsOrder   string = "order"
)

type (
	CreateUserRq struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		Code     int    `json:"code"`
	}

	GetValueRq struct {
		Value string `json:"value"`
	}

	CreateUserRs struct {
		UserID uuid.UUID `json:"user_id"`
	}

	UserRs struct {
		ID        uuid.UUID    `json:"id"`
		Nickname  string       `json:"nickname"`
		Email     string       `json:"email,omitempty"`
		Role      runtime.Role `json:"role"`
		VisitedAt time.Time    `json:"visited_at,omitempty"`
	}
)

type Handler struct {
	userSvc *user.Service
}

func Create(g *ginext.Engine, userSvc *user.Service) {
	h := newHandler(userSvc)

	g.POST(handler.SignUp, h.signUp)
	g.GET(handler.UserByID, middleware.Auth(h.getUserByID))
	g.GET(handler.Users, h.getUsers)

	h.initSettingsHandler(g)
}

func newHandler(userSvc *user.Service) *Handler {
	return &Handler{
		userSvc: userSvc,
	}
}

func (h *Handler) signUp(c *ginext.Context) (int, any, error) {
	var data CreateUserRq
	err := c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(
				fmt.Errorf("user.delivery.Handler.signUp: %v", err),
				msgerr.ErrMsgBadRequest)

	}

	if !utils.IsPasswordValid(data.Password) {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("user.delivery.Handler.signUp: invalid password"),
				"Invalid password")
	}

	if !utils.IsEmailValid(data.Email) {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("user.delivery.Handler.signUp: invalid email"),
				msgerr.ErrMsgBadEmail)
	}

	uid, err := h.userSvc.SignUp(c.Request.Context(), entity.UserCreate{
		Nickname: strings.Split(data.Email, "@")[0],
		Password: data.Password,
		Email:    data.Email,
		Role:     runtime.User,
		Code:     data.Code,
	})
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("user.delivery.Handler.signUp: %v", err)
	}

	createUserRs := &CreateUserRs{
		UserID: uid,
	}

	return http.StatusCreated, createUserRs, nil
}

func (h *Handler) getUserByID(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil, fmt.Errorf("user.delivery.Handler.getUserByID: %v", err)
	}
	usr, err := h.userSvc.GetUserByID(ctx, userID)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("user.delivery.Handler.getUserByID: %v", err)
	}

	userRs := &UserRs{
		ID:       usr.ID,
		Nickname: usr.Nickname,
		Email:    usr.Email,
		Role:     usr.Role,
	}

	return http.StatusOK, userRs, nil
}

func (h *Handler) getUsers(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, _ := runtime.UserIDFromContext(ctx)

	page, err := c.GetQueryInt(paramsPage)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("user.delivery.Handler.getUsers: %v", err)
	}

	perPage, err := c.GetQueryInt(paramsPerPage)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("user.delivery.Handler.getUsers: %v", err)
	}

	search, err := c.GetQuery(paramsSearch)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("user.delivery.Handler.getUsers: %v", err)
	}

	sort, err := c.GetQueryInt(paramsSort)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("user.delivery.Handler.getUsers: %v", err)
	}

	order, err := c.GetQueryInt(paramsOrder)
	if err != nil {
		return http.StatusBadRequest, nil,
			fmt.Errorf("user.delivery.Handler.getUsers: %v", err)
	}

	users, countUsers, err := h.userSvc.GetUsers(ctx, uid, page, perPage, sort, order, search)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("user.delivery.Handler.getUsers: %v", err)
	}

	usersRs := make([]UserRs, 0, len(users))
	for _, u := range users {
		usersRs = append(usersRs, UserRs{
			ID:        u.ID,
			Nickname:  u.Nickname,
			Role:      u.Role,
			VisitedAt: u.VisitedAt,
		})
	}

	return http.StatusOK, gin.H{
		"users":       usersRs,
		"count_users": countUsers,
	}, nil
}
