package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

var (
	errEmptyBody = errors.New("body is empty")
)

func CheckBody(w http.ResponseWriter, r *http.Request, body any) error {
	if r.Body == http.NoBody {
		return fmt.Errorf("tools.CheckBody: %w", errEmptyBody)
	}

	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		return fmt.Errorf("tools.CheckBody - decode: %w", err)
	}

	return nil
}

func SendError(w http.ResponseWriter, httpStatus int, err error) {
	w.WriteHeader(httpStatus)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		slog.Error(fmt.Errorf("internal.tools.handler.SendError: %v", err).Error())
	}
}
