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
	SecurityCode struct {
		Data string `json:"data"`
		Type string `json:"type"`
	}
)

func (h *Handler) initSettingsHandler(g *ginext.Engine) {
	g.POST(handler.UserSecurityCode, middleware.Auth(h.securityCode))
}

func (h *Handler) securityCode(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("user.delivery.Handler.securityCode: %v", err),
				msgerr.ErrMsgUnauthorized)
	}

	var data SecurityCode
	err = c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerr.New(fmt.Errorf("user.delivery.Handler.securityCode: %v", err),
				msgerr.ErrMsgBadRequest)
	}

	err = h.userSvc.SendSecurityCodeForUpdatePsw(ctx, uid, data.Data, data.Type)
	if err != nil {
		return http.StatusInternalServerError, nil,
			fmt.Errorf("user.delivery.Handler.securityCode: %v", err)
	}

	return http.StatusOK, nil, nil
}
