package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery"
	ginExtension "github.com/av-ugolkov/lingua-evo/internal/pkg/http/gin_extension"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/middleware"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	"github.com/av-ugolkov/lingua-evo/internal/services/user"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/user"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		ID    uuid.UUID    `json:"id"`
		Name  string       `json:"name"`
		Email string       `json:"email"`
		Role  runtime.Role `json:"role"`
	}

	UserEditPasswordRq struct {
		OldPassword string `json:"old_password"`
		Password    string `json:"password"`
		EmailCode   int    `json:"email_code"`
	}

	UserEditEmailRq struct {
		OldEmail string `json:"old_email"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	UserEditUsernameRq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
)

type Handler struct {
	userSvc *user.Service
}

func Create(r *gin.Engine, userSvc *user.Service) {
	h := newHandler(userSvc)
	h.register(r)
}

func newHandler(userSvc *user.Service) *Handler {
	return &Handler{
		userSvc: userSvc,
	}
}

func (h *Handler) register(r *gin.Engine) {
	r.POST(delivery.SignUp, h.signUp)
	r.GET(delivery.UserByID, middleware.Auth(h.getUserByID))
	r.POST(delivery.UserEditPassword, middleware.Auth(h.editPassword))
	r.POST(delivery.UserEditEmail, middleware.Auth(h.editEmail))
	r.POST(delivery.UserEditName, middleware.Auth(h.editName))
}

func (h *Handler) signUp(c *gin.Context) {
	var data CreateUserRq
	err := c.Bind(&data)
	if err != nil {
		ginExtension.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.createAccount - check body: %v", err))
		return
	}

	if !utils.IsUsernameValid(data.Username) {
		ginExtension.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.createAccount - invalid user name"),
		)
		return
	}

	if !utils.IsPasswordValid(data.Password) {
		ginExtension.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.createAccount - invalid password"),
		)
		return
	}

	if !utils.IsEmailValid(data.Email) {
		ginExtension.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.createAccount - invalid email"),
		)
		return
	}

	uid, err := h.userSvc.SignUp(c.Request.Context(), entity.UserData{
		ID:       uuid.New(),
		Name:     data.Username,
		Password: data.Password,
		Email:    data.Email,
		Role:     runtime.User,
		Code:     data.Code,
	})
	if err != nil {
		ginExtension.SendError(c, http.StatusInternalServerError,
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

func (h *Handler) editPassword(c *gin.Context) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExtension.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.editPassword: not found user id"))
		return
	}
	var data UserEditPasswordRq
	err = c.Bind(&data)
	if err != nil {
		ginExtension.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.editPassword - check body: %v", err))
		return
	}

	err = h.userSvc.EditPassword(ctx, entity.UserPasword{
		ID:          uid,
		OldPassword: data.OldPassword,
		Password:    data.Password,
		Code:        data.EmailCode,
	})
	if err != nil {
		ginExtension.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("user.delivery.Handler.editPassword: %v", err),
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) editEmail(c *gin.Context) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExtension.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.editPassword: not found user id"))
		return
	}
	var data UserEditEmailRq
	err = c.Bind(&data)
	if err != nil {
		ginExtension.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.editPassword - check body: %v", err))
		return
	}

	err = h.userSvc.EditEmail(ctx, entity.EditUserData{
		ID:       uid,
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		ginExtension.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("user.delivery.Handler.editPassword: %v", err),
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) editName(c *gin.Context) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		ginExtension.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.editPassword: not found user id"))
		return
	}
	var data UserEditUsernameRq
	err = c.Bind(&data)
	if err != nil {
		ginExtension.SendError(c, http.StatusBadRequest,
			fmt.Errorf("user.delivery.Handler.editPassword - check body: %v", err))
		return
	}

	err = h.userSvc.EditUsername(ctx, entity.EditUserData{
		ID:       uid,
		Username: data.Username,
		Password: data.Password,
	})
	if err != nil {
		ginExtension.SendError(c, http.StatusInternalServerError,
			fmt.Errorf("user.delivery.Handler.editPassword: %v", err),
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
