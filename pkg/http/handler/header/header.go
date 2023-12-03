package header

import (
	"errors"
	"fmt"
	"lingua-evo/runtime"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	language     = "language"
	refreshToken = "refresh_token"

	cookiePath     = "/"
	cookiePathAuth = "/auth"

	headerFingerprint = "Fingerprint"
)

var (
	errHeaderNotFound = errors.New("header not found")
)

type Header struct {
	responseWriter http.ResponseWriter
	request        *http.Request
}

func NewHeader(w http.ResponseWriter, r *http.Request) *Header {
	return &Header{
		responseWriter: w,
		request:        r,
	}
}

func (h *Header) setCookie(name, value string) {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  cookiePath,
	}
	http.SetCookie(h.responseWriter, cookie)
}

func (h *Header) SetCookieLanguage(languageID string) {
	h.setCookie(language, languageID)
}

func (h *Header) SetCookieRefreshToken(token uuid.UUID, maxAge time.Duration) {
	cookie := &http.Cookie{
		Name:     refreshToken,
		Value:    token.String(),
		MaxAge:   int(maxAge.Seconds()),
		HttpOnly: true,
		Secure:   true,
		Path:     cookiePathAuth,
	}
	http.SetCookie(h.responseWriter, cookie)
}

func (h *Header) GetCookieRefreshToken() (*http.Cookie, error) {
	cookie, err := h.request.Cookie(refreshToken)
	switch {
	case errors.Is(err, http.ErrNoCookie):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("tools.GetCookie: %w", err)
	default:
		return cookie, nil
	}
}

func (h *Header) GetCookieLanguageOrDefault() string {
	cookie, err := h.request.Cookie(language)
	switch {
	case errors.Is(err, http.ErrNoCookie):
		return runtime.GetLanguage("en")
	case err != nil:
		slog.Error(fmt.Errorf("tools.GetCookie: %w", err).Error())
		return runtime.GetLanguage("en")
	default:
		return cookie.Value
	}
}

func (h *Header) SetHeader(ney, value string) {
	h.request.Header.Set(ney, value)
}

func (h *Header) getHeader(name string) (string, error) {
	value := h.request.Header.Get(name)
	if value == "" {
		return "", fmt.Errorf("http.handler.header.getHeader: %w", errHeaderNotFound)
	}
	return value, nil
}

func (h *Header) GetHeaderAuthorization() (string, error) {
	return h.getHeader("Authorization")
}

func (h *Header) GetHeaderFingerprint() (string, error) {
	return h.getHeader(headerFingerprint)
}
