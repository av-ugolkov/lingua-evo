package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	ginext "github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	msgerr "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/user"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/gin-gonic/gin"
)

type (
	SecurityCodeRq struct {
		Data string `json:"data"`
	}

	UpdatePswRq struct {
		OldPsw string `json:"old_psw"`
		NewPsw string `json:"new_psw"`
		Code   string `json:"code"`
	}

	UpdateEmailRq struct {
		NewEmail string `json:"new_email"`
		Code     string `json:"code"`
	}

	UpdateNickname struct {
		Nickname string `json:"nickname"`
	}
)

func (h *Handler) initSettingsHandler(g *ginext.Engine) {
	g.GET(handler.AccountSettingsAccount, middleware.Auth(h.getSettingsAccount))
	g.GET(handler.AccountSettingsPersonalInfo, middleware.Auth(h.getSettingsPersonalInfo))
	g.GET(handler.AccountSettingsEmailNotif, middleware.Auth(h.getSettingsEmailNotif))
	g.POST(handler.AccountSettingsUpdatePswCode, middleware.Auth(h.updatePswSendCode))
	g.POST(handler.AccountSettingsUpdatePsw, middleware.Auth(h.updatePsw))
	g.POST(handler.AccountSettingsUpdateEmailCode, middleware.Auth(h.updateEmailSendCode))
	g.POST(handler.AccountSettingsUpdateEmail, middleware.Auth(h.updateEmail))
	g.POST(handler.AccountSettingsUpdateNickname, middleware.Auth(h.updateNickname))
}

func (h *Handler) getSettingsAccount(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("user.handler.Handler.getSettingsAccount: %v", err),
				msgerr.ErrMsgUnauthorized)
	}

	usr, err := h.userSvc.GetUserByID(ctx, uid)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("user.handler.Handler.getSettingsAccount: %v", err)
	}

	userRs := UserRs{
		Nickname: usr.Nickname,
		Email:    usr.Email,
	}

	return http.StatusOK, userRs, nil
}

func (h *Handler) getSettingsPersonalInfo(c *ginext.Context) (int, any, error) {
	return http.StatusOK, nil, nil
}

func (h *Handler) getSettingsEmailNotif(c *ginext.Context) (int, any, error) {
	return http.StatusOK, nil, nil
}

func (h *Handler) updatePswSendCode(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("user.handler.Handler.updatePswSendCode: %v", err),
				msgerr.ErrMsgUnauthorized)
	}

	var data SecurityCodeRq
	err = c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("user.handler.Handler.updatePswSendCode: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	err = h.userSvc.SendSecurityCodeForUpdatePsw(ctx, uid, data.Data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("user.handler.Handler.updatePswSendCode: %v", err)
	}

	return http.StatusOK, nil, nil
}

func (h *Handler) updatePsw(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("user.handler.Handler.updatePsw: %v", err),
				msgerr.ErrMsgUnauthorized)
	}

	var data UpdatePswRq
	err = c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("user.handler.Handler.updatePsw: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	err = h.userSvc.UpdatePsw(ctx, uid, data.OldPsw, data.NewPsw, data.Code)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("user.handler.Handler.updatePsw: %v", err)
	}

	return http.StatusOK, nil, nil
}

func (h *Handler) updateEmailSendCode(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("user.handler.Handler.updateEmailSendCode: %v", err),
				msgerr.ErrMsgUnauthorized)
	}

	ttl, err := h.userSvc.SendSecurityCodeForUpdateEmail(ctx, uid)
	switch {
	case errors.Is(err, entity.ErrDuplicateCode):
		return http.StatusConflict, gin.H{"ttl": ttl}, fmt.Errorf("user.handler.Handler.updateEmailSendCode: %w", err)
	case err != nil:
		return http.StatusInternalServerError, nil,
			fmt.Errorf("user.handler.Handler.updateEmailSendCode: %w", err)
	}

	return http.StatusOK, gin.H{"msg": "Could you check your email. We have sent you a code for updating your email"}, nil
}

func (h *Handler) updateEmail(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("user.handler.Handler.updateEmail: %v", err),
				msgerr.ErrMsgUnauthorized)
	}

	var data UpdateEmailRq
	err = c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("user.handler.Handler.updateEmail: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	err = h.userSvc.UpdateEmail(ctx, uid, data.NewEmail, data.Code)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("user.handler.Handler.updateEmail: %w", err)
	}

	return http.StatusOK, nil, nil
}

func (h *Handler) updateNickname(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("user.handler.Handler.updateNickname: %v", err),
				msgerr.ErrMsgUnauthorized)
	}

	var data UpdateNickname
	err = c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("user.handler.Handler.updateNickname: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	err = h.userSvc.UpdateNickname(ctx, uid, data.Nickname)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("user.handler.Handler.updateNickname: %w", err)
	}

	return http.StatusOK, nil, nil
}
