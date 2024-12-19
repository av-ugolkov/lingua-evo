package handler

import (
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/fext"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msgerr"
	"github.com/av-ugolkov/lingua-evo/internal/services/games/service"

	"github.com/gofiber/fiber/v2"
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

func Create(r *fiber.App, gameSvc *service.Service) {
	h := &Handler{
		gameSvc: gameSvc,
	}

	r.Get(handler.ReviseGame, h.getReviseGame)
}

func (h *Handler) getReviseGame(c *fiber.Ctx) error {
	ctx := c.Context()

	uid, err := fext.UserIDFromContext(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fext.E(err, msgerr.ErrMsgUnauthorized))
	}

	var data ReviseGameRq
	err = c.BodyParser(&data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	err = h.gameSvc.GameRevise(ctx, uid, reviseGameFromRsToEntity(data))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fext.E(err, msgerr.ErrMsgInternal))
	}

	return c.SendStatus(http.StatusOK)
}
