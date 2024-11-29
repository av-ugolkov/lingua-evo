package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
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
	GetValueRq struct {
		Value string `json:"value"`
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
	h := &Handler{
		userSvc: userSvc,
	}

	g.GET(handler.UserByID, middleware.Auth(h.getUserByID))
	g.GET(handler.Users, h.getUsers)

	h.initSettingsHandler(g)
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
