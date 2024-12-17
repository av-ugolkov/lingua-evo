package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	user "github.com/av-ugolkov/lingua-evo/internal/services/user/service"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gofiber/fiber/v2"
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

func Create(r *fiber.App, userSvc *user.Service) {
	h := &Handler{
		userSvc: userSvc,
	}

	r.Get(handler.UserByID, middleware.Auth(h.getUserByID))
	r.Get(handler.Users, h.getUsers)

	h.initSettingsHandler(r)
}

func (h *Handler) getUserByID(c *fiber.Ctx) error {
	ctx := c.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized,
			fmt.Sprintf("user.delivery.Handler.getUserByID: %v", err))
	}
	usr, err := h.userSvc.GetUserByID(ctx, userID)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError,
			fmt.Sprintf("user.delivery.Handler.getUserByID: %v", err))
	}

	userRs := &UserRs{
		ID:       usr.ID,
		Nickname: usr.Nickname,
		Email:    usr.Email,
		Role:     usr.Role,
	}

	return c.Status(http.StatusOK).JSON(userRs)
}

func (h *Handler) getUsers(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, _ := runtime.UserIDFromContext(ctx)

	page := c.QueryInt(paramsPage, 1)
	perPage := c.QueryInt(paramsPerPage)
	if perPage == 0 {
		return fiber.NewError(fiber.StatusBadRequest,
			fmt.Sprintf("user.delivery.Handler.getUsers: not found query [%s]", paramsPerPage))
	}

	search := c.Query(paramsSearch)

	sort := c.QueryInt(paramsSort, -1)
	if sort == -1 {
		return fiber.NewError(fiber.StatusBadRequest,
			fmt.Sprintf("user.delivery.Handler.getUsers: not found query [%s]", paramsSort))
	}

	order := c.QueryInt(paramsOrder)
	if order == -1 {
		return fiber.NewError(fiber.StatusBadRequest,
			fmt.Sprintf("user.delivery.Handler.getUsers: not found query [%s]", paramsOrder))
	}

	users, countUsers, err := h.userSvc.GetUsers(ctx, uid, page, perPage, sort, order, search)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError,
			fmt.Sprintf("user.delivery.Handler.getUsers: %v", err))
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

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"users":       usersRs,
		"count_users": countUsers,
	})
}
