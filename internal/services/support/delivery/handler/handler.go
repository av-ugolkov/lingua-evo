package handler

import (
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	"github.com/av-ugolkov/lingua-evo/internal/services/support"

	"github.com/gin-gonic/gin"
)

type (
	SupportRq struct {
		Email   string `json:"email"`
		Name    string `json:"name,omitempty"`
		Type    string `json:"type,omitempty"`
		Message string `json:"message"`
	}
)

type Handler struct {
	supportSvc *support.Service
}

func Create(r *ginext.Engine, supportSvc *support.Service) {
	h := &Handler{
		supportSvc: supportSvc,
	}

	r.POST(handler.SupportRequest, h.sendRequest)
}

func (h *Handler) sendRequest(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	var data SupportRq
	err := c.Bind(&data)
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerror.New(fmt.Errorf("support.delivery.Handler.sendRequest: %v", err),
				msgerror.ErrMsgBadRequest)
	}

	if !utils.IsEmailValid(data.Email) {
		return http.StatusBadRequest, nil,
			msgerror.New(fmt.Errorf("support.delivery.Handler.sendRequest - email format is invalid"),
				msgerror.ErrMsgBadEmail)
	}

	if utf8.RuneCountInString(data.Name) > 100 {
		data.Name = fmt.Sprintf("%s...", data.Name[:100])
	}

	if len(data.Message) == 0 {
		return http.StatusBadRequest, nil,
			msgerror.New(fmt.Errorf("support.delivery.Handler.sendRequest - message is too long"),
				"Message is empty")
	}

	if utf8.RuneCountInString(data.Message) > 500 {
		return http.StatusBadRequest, nil,
			msgerror.New(fmt.Errorf("support.delivery.Handler.sendRequest - message is too long"),
				"Message is too long")
	}

	err = h.supportSvc.SendRequest(ctx, support.SupportRequest{
		Email:   data.Email,
		Name:    data.Name,
		Type:    data.Type,
		Message: data.Message,
	})
	if err != nil {
		return http.StatusBadRequest, nil,
			msgerror.New(fmt.Errorf("support.delivery.Handler.sendRequest - check body: %w", err),
				"Something went wrong. Try a bit later!")
	}

	return http.StatusCreated, gin.H{
		"msg": "Thank you for your message! We will try to respond as soon as possible.",
	}, nil
}
