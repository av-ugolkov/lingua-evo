package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"lingua-evo/pkg/http/handler/common"
	"lingua-evo/runtime"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	cookiePath     = "/"
	cookiePathAuth = "/auth"

	headerFingerprint = "Fingerprint"
)

var (
	errEmptyBody      = errors.New("body is empty")
	errHeaderNotFound = errors.New("header not found")
)

var (
	emptyJson = []byte("{}")
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

func (h *Handler) SendData(httpStatus int, data []byte) {
	h.responseWriter.WriteHeader(httpStatus)
	_, err := h.responseWriter.Write(data)
	if err != nil {
		slog.Error(fmt.Errorf("http.handler.SendError: %v", err).Error())
	}
}

func (h *Handler) SendEmptyData(httpStatus int) {
	h.responseWriter.WriteHeader(httpStatus)
	_, err := h.responseWriter.Write(emptyJson)
	if err != nil {
		slog.Error(fmt.Errorf("http.handler.SendError: %v", err).Error())
	}
}

func (h *Handler) SendError(httpStatus int, err error) {
	h.responseWriter.WriteHeader(httpStatus)
	slog.Error(fmt.Errorf("http.handler.SendError: %v", err).Error())
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

func (h *Handler) DeleteCookie(name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(h.responseWriter, cookie)
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
	h.responseWriter.Header().Set(ney, value)
}

func (h *Handler) getHeader(name string) (string, error) {
	value := h.request.Header.Get(name)
	if value == "" {
		return "", fmt.Errorf("http.handler.getHeader: %w", errHeaderNotFound)
	}
	return value, nil
}

func (h *Handler) GetHeaderAuthorization(typeAuth common.TypeAuth) (string, error) {
	token, err := h.getHeader("Authorization")
	if err != nil {
		return "", fmt.Errorf("http.handler.GetHeaderAuthorization: %w", err)
	}

	if !strings.HasPrefix(token, string(typeAuth)) {
		return "", fmt.Errorf("http.handler.GetHeaderAuthorization - invalid type auth [%s]: %s", typeAuth, token)
	}

	return token[len(string(typeAuth))+1:], nil
}

func (h *Handler) GetHeaderFingerprint() (string, error) {
	return h.getHeader(headerFingerprint)
}
