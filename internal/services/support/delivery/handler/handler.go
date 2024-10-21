package handler

import (
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	ginExt "github.com/av-ugolkov/lingua-evo/internal/delivery/handler/gin"
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

func Create(r *gin.Engine, supportSvc *support.Service) {
	h := &Handler{
		supportSvc: supportSvc,
	}
	h.register(r)
}

func (h *Handler) register(r *gin.Engine) {
	r.POST(handler.SupportRequest, h.sendRequest)
}

func (h *Handler) sendRequest(c *gin.Context) {
	ctx := c.Request.Context()

	var data SupportRq
	err := c.Bind(&data)
	if err != nil {
		ginExt.SendErrorWithMsg(c, http.StatusBadRequest,
			fmt.Errorf("support.delivery.Handler.sendRequest - check body: %w", err),
			"You doesn't fill one or several fields.")
		return
	}

	if !utils.IsEmailValid(data.Email) {
		ginExt.SendErrorWithMsg(c, http.StatusBadRequest,
			fmt.Errorf("support.delivery.Handler.sendRequest - email format is invalid"),
			"Email format is invalid")
		return
	}

	if utf8.RuneCountInString(data.Name) > 100 {
		data.Name = fmt.Sprintf("%s...", data.Name[:100])
	}

	if len(data.Message) == 0 {
		ginExt.SendErrorWithMsg(c, http.StatusBadRequest,
			fmt.Errorf("support.delivery.Handler.sendRequest - message is too long"),
			"Message is empty")
		return
	}

	if utf8.RuneCountInString(data.Message) > 500 {
		ginExt.SendErrorWithMsg(c, http.StatusBadRequest,
			fmt.Errorf("support.delivery.Handler.sendRequest - message is too long"),
			"Message is too long")
		return
	}

	err = h.supportSvc.SendRequest(ctx, support.SupportRequest{
		Email:   data.Email,
		Name:    data.Name,
		Type:    data.Type,
		Message: data.Message,
	})
	if err != nil {
		ginExt.SendErrorWithMsg(c, http.StatusBadRequest,
			fmt.Errorf("support.delivery.Handler.sendRequest - check body: %w", err),
			"Something went wrong. Try a bit later!")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"msg": "Thank you for your message! We will try to respond as soon as possible.",
	})
}
