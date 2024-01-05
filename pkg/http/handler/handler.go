package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"lingua-evo/pkg/http/handler/common"
	"lingua-evo/runtime"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	cookiePath     = "/"
	cookiePathAuth = "/auth"

	headerFingerprint = "Fingerprint"
)

var (
	errEmptyBody = errors.New("body is empty")
)

var (
	errHeaderNotFound = errors.New("header not found")
)

type Handler struct {
	responseWriter http.ResponseWriter
	request        *http.Request
}

func NewHandler(rw http.ResponseWriter, r *http.Request) *Handler {
	return &Handler{
		responseWriter: rw,
		request:        r,
	}
}

func (h *Handler) CheckBody(body any) error {
	if h.request.Body == http.NoBody {
		return fmt.Errorf("http.handler.CheckBody: %w", errEmptyBody)
	}

	err := json.NewDecoder(h.request.Body).Decode(body)
	if err != nil {
		return fmt.Errorf("http.handler.CheckBody - decode: %w", err)
	}

	return nil
}

func (h *Handler) SendData(data []byte) {
	_, err := h.responseWriter.Write(data)
	if err != nil {
		slog.Error(fmt.Errorf("http.handler.SendError: %v", err).Error())
	}
}

func (h *Handler) SendError(httpStatus int, err error) {
	h.responseWriter.WriteHeader(httpStatus)
	_, err = h.responseWriter.Write([]byte(err.Error()))
	if err != nil {
		slog.Error(fmt.Errorf("http.handler.SendError: %v", err).Error())
	}
}

func (h *Handler) setCookie(name, value string) {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  cookiePath,
	}
	http.SetCookie(h.responseWriter, cookie)
}

func (h *Handler) SetCookieLanguage(languageID string) {
	h.setCookie(common.Language, languageID)
}

func (h *Handler) SetCookieRefreshToken(token uuid.UUID, maxAge time.Duration) {
	cookie := &http.Cookie{
		Name:     common.RefreshToken,
		Value:    token.String(),
		MaxAge:   int(maxAge.Seconds()),
		HttpOnly: true,
		Secure:   true,
		Path:     cookiePathAuth,
	}
	http.SetCookie(h.responseWriter, cookie)
}

func (h *Handler) Cookie(name string) (*http.Cookie, error) {
	cookie, err := h.request.Cookie(name)
	switch {
	case errors.Is(err, http.ErrNoCookie):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("http.handler.GetCookie: %w", err)
	default:
		return cookie, nil
	}
}

func (h *Handler) GetCookieLanguageOrDefault() string {
	cookie, err := h.request.Cookie(common.Language)
	switch {
	case errors.Is(err, http.ErrNoCookie):
		return runtime.GetLanguage("en")
	case err != nil:
		slog.Error(fmt.Errorf("http.handler.GetCookieLanguageOrDefault: %w", err).Error())
		return runtime.GetLanguage("en")
	default:
		return cookie.Value
	}
}

func (h *Handler) WriteHeader(httpStatus int) {
	h.responseWriter.WriteHeader(httpStatus)
}

func (h *Handler) SetHeader(ney, value string) {
	h.request.Header.Set(ney, value)
}

func (h *Handler) getHeader(name string) (string, error) {
	value := h.request.Header.Get(name)
	if value == "" {
		return "", fmt.Errorf("http.handler.getHeader: %w", errHeaderNotFound)
	}
	return value, nil
}

func (h *Handler) GetHeaderAuthorization() (string, error) {
	return h.getHeader("Authorization")
}

func (h *Handler) GetHeaderAccessToken() (string, error) {
	return h.getHeader("Access-Token")
}

func (h *Handler) GetHeaderFingerprint() (string, error) {
	return h.getHeader(headerFingerprint)
}
