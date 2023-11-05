package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
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

func SendData(w http.ResponseWriter, data []byte) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func SendError(w http.ResponseWriter, httpStatus int, err error) {
	w.WriteHeader(httpStatus)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		slog.Error(fmt.Errorf("internal.tools.handler.SendError: %v", err).Error())
	}
}

func GetFingerprint(r *http.Request) string {
	return strings.Join([]string{
		r.Header.Get("user-agent"),
		r.Header.Get("sec-ch-ua"),
		r.Header.Get("accept-language"),
		r.Header.Get("upgrade-insecure-req"),
	}, ":")
}
