package fext

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Ferr interface {
	Msg() string
}

type contextUserIDKey struct{}

var ctxUserIDKey = &contextUserIDKey{}

var (
	errUserIDNotFound = errors.New("user id not found")
)

func SetUserIDToContext(c *fiber.Ctx, uid uuid.UUID) {
	c.Locals(ctxUserIDKey, uid)
}

func UserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals(ctxUserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errUserIDNotFound
	}
	return userID, nil
}

func D(data any) fiber.Map {
	return fiber.Map{
		"data": data,
	}
}

func E(err error, msg ...string) fiber.Map {
	slog.Error(err.Error())
	if len(msg) > 0 {
		return fiber.Map{
			"msg": msg[0],
		}
	}

	var ferr Ferr
	switch {
	case errors.As(err, &ferr):
		return fiber.Map{
			"msg": ferr.Msg(),
		}
	default:
		return fiber.Map{
			"msg": err.Error(),
		}

	}
}

func DE(data any, err error) fiber.Map {
	jsonData := fiber.Map{
		"data": data,
	}
	if err != nil {
		var ferr Ferr
		switch {
		case errors.As(err, &ferr):
			jsonData["msg"] = ferr.Msg()
		default:
			jsonData["msg"] = err.Error()
		}
	}
	return jsonData
}
