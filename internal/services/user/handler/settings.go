package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler/middleware"
	ginext "github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	msgerr "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/runtime"
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
)

func (h *Handler) initSettingsHandler(g *ginext.Engine) {
	g.POST(handler.UserUpdatePswCode, middleware.Auth(h.updatePswSendCode))
	g.POST(handler.UserUpdatePsw, middleware.Auth(h.updatePsw))
	g.POST(handler.UserUpdateEmailCode, middleware.Auth(h.updateEmailSendCode))
	g.POST(handler.UserUpdateEmail, middleware.Auth(h.updateEmail))
}

func (h *Handler) updatePswSendCode(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("user.delivery.Handler.updatePswSendCode: %v", err),
				msgerr.ErrMsgUnauthorized)
	}

	var data SecurityCodeRq
	err = c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("user.delivery.Handler.updatePswSendCode: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	err = h.userSvc.SendSecurityCodeForUpdatePsw(ctx, uid, data.Data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("user.delivery.Handler.updatePswSendCode: %v", err)
	}

	return http.StatusOK, nil, nil
}

func (h *Handler) updatePsw(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("user.delivery.Handler.updatePsw: %v", err),
				msgerr.ErrMsgUnauthorized)
	}

	var data UpdatePswRq
	err = c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("user.delivery.Handler.updatePsw: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	err = h.userSvc.UpdatePsw(ctx, uid, data.OldPsw, data.NewPsw, data.Code)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("user.delivery.Handler.updatePsw: %v", err)
	}

	return http.StatusOK, nil, nil
}

func (h *Handler) updateEmailSendCode(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("user.delivery.Handler.updateEmailSendCode: %v", err),
				msgerr.ErrMsgUnauthorized)
	}

	var data SecurityCodeRq
	err = c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("user.delivery.Handler.updateEmailSendCode: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	err = h.userSvc.SendSecurityCodeForUpdateEmail(ctx, uid, data.Data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("user.delivery.Handler.updateEmailSendCode: %v", err)
	}

	return http.StatusOK, nil, nil
}

func (h *Handler) updateEmail(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("user.delivery.Handler.updatePsw: %v", err),
				msgerr.ErrMsgUnauthorized)
	}

	var data UpdatePswRq
	err = c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("user.delivery.Handler.updatePsw: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	err = h.userSvc.UpdatePsw(ctx, uid, data.OldPsw, data.NewPsw, data.Code)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("user.delivery.Handler.updatePsw: %v", err)
	}

	return http.StatusOK, nil, nil
}
