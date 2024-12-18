package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/fext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
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
	userID, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}
	usr, err := h.userSvc.GetUserByID(ctx, userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fext.E(err))
	}

	userRs := &UserRs{
		ID:       usr.ID,
		Nickname: usr.Nickname,
		Email:    usr.Email,
		Role:     usr.Role,
	}

	return c.Status(http.StatusOK).JSON(fext.D(userRs))
}

func (h *Handler) getUsers(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, _ := fext.UserIDFromContext(c)

	page := c.QueryInt(paramsPage, 1)
	perPage := c.QueryInt(paramsPerPage)
	if perPage == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("user.Handler.getUsers: not found query [%s]", paramsPerPage)))
	}

	search := c.Query(paramsSearch)

	sort := c.QueryInt(paramsSort, -1)
	if sort == -1 {
		return c.Status(fiber.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("user.Handler.getUsers: not found query [%s]", paramsSort)))
	}

	order := c.QueryInt(paramsOrder)
	if order == -1 {
		return c.Status(fiber.StatusBadRequest).JSON(fext.E(
			fmt.Errorf("user.Handler.getUsers: not found query [%s]", paramsOrder)))
	}

	users, countUsers, err := h.userSvc.GetUsers(ctx, uid, page, perPage, sort, order, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fext.E(err))
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

	return c.Status(http.StatusOK).JSON(fext.D(fiber.Map{
		"users":       usersRs,
		"count_users": countUsers,
	}))
}
