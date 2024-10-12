package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	"github.com/av-ugolkov/lingua-evo/internal/services/user"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/user"
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
		Username string `json:"username"`
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
		ID          uuid.UUID    `json:"id"`
		Name        string       `json:"name"`
		Email       string       `json:"email,omitempty"`
		Role        runtime.Role `json:"role"`
		LastVisited time.Time    `json:"last_visited,omitempty"`
	}
)

type Handler struct {
	userSvc *user.Service
}

func Create(g *gin.Engine, userSvc *user.Service) {
	h := newHandler(userSvc)
	h.register(g)
}

func newHandler(userSvc *user.Service) *Handler {
	return &Handler{
		userSvc: userSvc,
	}
}

func (h *Handler) register(g *gin.Engine) {
	g.POST(handler.SignUp, h.signUp)
	g.GET(handler.UserByID, middleware.Auth(h.getUserByID))
	g.GET(handler.Users, h.getUsers)
}

func (h *Handler) signUp(c *gin.Context) {
	var data CreateUserRq
	err := c.Bind(&data)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.createAccount - check body: %v", err))
		return
	}

	if !utils.IsUsernameValid(data.Username) {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.createAccount - invalid user name"),
		)
		return
	}

	if !utils.IsPasswordValid(data.Password) {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.createAccount - invalid password"),
		)
		return
	}

	if !utils.IsEmailValid(data.Email) {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.createAccount - invalid email"),
		)
		return
	}

	uid, err := h.userSvc.SignUp(c.Request.Context(), entity.UserCreate{
		Name:     data.Username,
		Password: data.Password,
		Email:    data.Email,
		Role:     runtime.User,
		Code:     data.Code,
	})
	if err != nil {
		ginExt.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("user.delivery.Handler.createAccount - create user: %v", err),
		)
		return
	}

	createUserRs := &CreateUserRs{
		UserID: uid,
	}
	c.JSON(http.StatusCreated, createUserRs)
}

func (h *Handler) getUserByID(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Errorf("user.delivery.Handler.getUserByID - unauthorized: %v", err),
		})
		return
	}
	userData, err := h.userSvc.GetUserByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("user.delivery.Handler.getUserByID: %v", err),
		})
		return
	}

	userRs := &UserRs{
		ID:    userData.ID,
		Name:  userData.Name,
		Email: userData.Email,
		Role:  userData.Role,
	}

	c.JSON(http.StatusOK, userRs)
}

func (h *Handler) getUsers(c *gin.Context) {
	ctx := c.Request.Context()

	uid, _ := runtime.UserIDFromContext(ctx)

	page, err := ginExt.GetQueryInt(c, paramsPage)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.getUsers - get query [page]: %v", err))
		return
	}

	perPage, err := ginExt.GetQueryInt(c, paramsPerPage)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.getUsers - get query [per_page]: %v", err))
		return
	}

	search, err := ginExt.GetQuery(c, paramsSearch)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.getUsers - get query [search]: %v", err))
		return
	}

	sort, err := ginExt.GetQueryInt(c, paramsSort)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.getUsers - get query [sort]: %v", err))
		return
	}

	order, err := ginExt.GetQueryInt(c, paramsOrder)
	if err != nil {
		ginExt.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.getUsers - get query [order]: %v", err))
		return
	}

	users, countUsers, err := h.userSvc.GetUsers(ctx, uid, page, perPage, sort, order, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("user.delivery.Handler.getUsers: %v", err),
		})
		return
	}

	usersRs := make([]UserRs, 0, len(users))
	for _, u := range users {
		usersRs = append(usersRs, UserRs{
			ID:          u.ID,
			Name:        u.Name,
			Role:        u.Role,
			LastVisited: u.LastVisited,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users":       usersRs,
		"count_users": countUsers,
	})
}
