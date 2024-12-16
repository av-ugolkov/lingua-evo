package handler

import (
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	ginext "github.com/av-ugolkov/lingua-evo/internal/pkg/gin-ext"
	msgerr "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/internal/services/games/service"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
)

type (
	ReviseGameRq struct {
		VocabID   uuid.UUID `json:"vocab_id"`
		CountWord int       `json:"count_word"`
		TypeGame  string    `json:"type_game"`
	}
)

type Handler struct {
	gameSvc *service.Service
}

func Create(r *ginext.Engine, gameSvc *service.Service) {
	h := &Handler{
		gameSvc: gameSvc,
	}

	r.GET(handler.ReviseGame, h.getReviseGame)
}

func (h *Handler) getReviseGame(c *ginext.Context) (int, any, error) {
	ctx := c.Request.Context()

	uid, err := runtime.UserIDFromContext(ctx)
	if err != nil {
		return http.StatusUnauthorized, nil,
			msgerr.New(fmt.Errorf("games.delivery.Handler.getReviseGame: %v", err),
				msgerr.ErrMsgUnauthorized)
	}

	var data ReviseGameRq
	err = c.Bind(&data)
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("games.delivery.Handler.getReviseGame: %v", err),
				msgerr.ErrMsgInternal)
	}

	err = h.gameSvc.GameRevise(ctx, uid, reviseGameFromRsToEntity(data))
	if err != nil {
		return http.StatusInternalServerError, nil,
			msgerr.New(fmt.Errorf("games.delivery.Handler.getReviseGame: %v", err),
				msgerr.ErrMsgInternal)
	}

	return http.StatusOK, nil, nil
}
