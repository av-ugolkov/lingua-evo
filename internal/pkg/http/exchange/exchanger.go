package exchange

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/av-ugolkov/lingua-evo/runtime"
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
	emptyJson = []byte(runtime.EmptyJson)
)

type Exchanger struct {
	responseWriter http.ResponseWriter
	request        *http.Request
}

func NewExchanger(rw http.ResponseWriter, r *http.Request) *Exchanger {
	return &Exchanger{
		responseWriter: rw,
		request:        r,
	}
}

func (e *Exchanger) Context() context.Context {
	return e.request.Context()
}

func (e *Exchanger) CheckBody(body any) error {
	defer func() {
		_ = e.request.Body.Close()
	}()

	if e.request.Body == http.NoBody {
		return fmt.Errorf("http.exchange.Exchanger.CheckBody: %w", errEmptyBody)
	}

	err := json.NewDecoder(e.request.Body).Decode(body)
	if err != nil {
		return fmt.Errorf("http.exchange.Exchanger.CheckBody - decode: %w", err)
	}

	return nil
}

func (e *Exchanger) SendData(httpStatus int, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("http.exchange.Exchanger.SendData - marshal: %v", err)
		slog.Error(err.Error())
		e.SendError(http.StatusInternalServerError, err)
		return
	}

	e.responseWriter.WriteHeader(httpStatus)
	_, err = e.responseWriter.Write(b)
	if err != nil {
		err = fmt.Errorf("http.exchange.Exchanger.SendData: %v", err)
		slog.Error(err.Error())
		e.SendError(http.StatusInternalServerError, err)
	}
}

func (e *Exchanger) SendEmptyData(httpStatus int) {
	e.responseWriter.WriteHeader(httpStatus)
	_, err := e.responseWriter.Write(emptyJson)
	if err != nil {
		slog.Error(fmt.Errorf("http.exchange.Exchanger.SendError: %v", err).Error())
	}
}

func (e *Exchanger) SendError(httpStatus int, err error) {
	e.responseWriter.WriteHeader(httpStatus)
	slog.Error(fmt.Errorf("http.exchange.Exchanger.SendError: %v", err).Error())
	_, err = e.responseWriter.Write([]byte(err.Error()))
	if err != nil {
		slog.Error(fmt.Errorf("http.exchange.Exchanger.SendError: %v", err).Error())
	}
}

func (e *Exchanger) setCookie(name, value string) {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
		Path:  cookiePath,
	}
	http.SetCookie(e.responseWriter, cookie)
}

func (e *Exchanger) SetContentType(contentType ContentType) {
	e.SetHeader("Content-Type", string(contentType))
}

func (e *Exchanger) SetCookieLanguage(languageID string) {
	e.setCookie(Language, languageID)
}

func (e *Exchanger) SetCookieRefreshToken(token uuid.UUID, maxAge time.Duration) {
	cookie := &http.Cookie{
		Name:     RefreshToken,
		Value:    token.String(),
		MaxAge:   int(maxAge.Seconds()),
		HttpOnly: true,
		Secure:   true,
		Path:     cookiePathAuth,
	}
	http.SetCookie(e.responseWriter, cookie)
}

func (e *Exchanger) Cookie(name string) (*http.Cookie, error) {
	cookie, err := e.request.Cookie(name)
	switch {
	case errors.Is(err, http.ErrNoCookie):
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("http.exchange.Exchanger.GetCookie: %w", err)
	default:
		return cookie, nil
	}
}

func (e *Exchanger) DeleteCookie(name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    runtime.EmptyString,
		MaxAge:   0,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(e.responseWriter, cookie)
}

func (e *Exchanger) GetCookieLanguageOrDefault() string {
	cookie, err := e.request.Cookie(Language)
	switch {
	case errors.Is(err, http.ErrNoCookie):
		return runtime.GetLanguage("en")
	case err != nil:
		slog.Error(fmt.Errorf("http.exchange.Exchanger.GetCookieLanguageOrDefault: %w", err).Error())
		return runtime.GetLanguage("en")
	default:
		return cookie.Value
	}
}

func (e *Exchanger) SetHeader(ney, value string) {
	e.responseWriter.Header().Set(ney, value)
}

func (e *Exchanger) getHeader(name string) (string, error) {
	value := e.request.Header.Get(name)
	if value == runtime.EmptyString {
		return runtime.EmptyString, fmt.Errorf("http.exchange.Exchanger.getHeader: %w", errHeaderNotFound)
	}
	return value, nil
}

func (e *Exchanger) GetHeaderAuthorization(typeAuth TypeAuth) (string, error) {
	token, err := e.getHeader("Authorization")
	if err != nil {
		return runtime.EmptyString, fmt.Errorf("http.exchange.Exchanger.GetHeaderAuthorization: %w", err)
	}

	if !strings.HasPrefix(token, string(typeAuth)) {
		return runtime.EmptyString, fmt.Errorf("http.exchange.Exchanger.GetHeaderAuthorization - invalid type auth [%s]: %s", typeAuth, token)
	}

	return token[len(string(typeAuth))+1:], nil
}

func (e *Exchanger) GetHeaderFingerprint() (string, error) {
	return e.getHeader(headerFingerprint)
}

func (e *Exchanger) QueryParamString(key string) (string, error) {
	if !e.request.URL.Query().Has(key) {
		return runtime.EmptyString, fmt.Errorf("http.exchange.Exchanger.QueryParamString: not found query for key [%s]", key)
	}

	s := e.request.URL.Query().Get(key)
	return s, nil
}

func (e *Exchanger) QueryParamInt(key string) (int, error) {
	if !e.hasQueryParam(key) {
		return 0, fmt.Errorf("http.exchange.Exchanger.QueryParamInt: not found query for key [%s]", key)
	}

	s := e.request.URL.Query().Get(key)
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("http.exchange.Exchanger.QueryParamInt - cann't parse query param [%s]: %w", key, err)
	}
	return v, nil
}

func (e *Exchanger) QueryParamBool(key string) (bool, error) {
	if !e.hasQueryParam(key) {
		return false, fmt.Errorf("http.exchange.Exchanger.QueryParamBool: not found query for key [%s]", key)
	}

	s := e.request.URL.Query().Get(key)
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, fmt.Errorf("http.exchange.Exchanger.QueryParamBool - cann't parse query param [%s]: %w", key, err)
	}
	return v, nil
}

func (e *Exchanger) QueryParamUUID(key string) (uuid.UUID, error) {
	if !e.request.URL.Query().Has(key) {
		return uuid.Nil, fmt.Errorf("http.exchange.Exchanger.QueryParamString: not found query for key [%s]", key)
	}

	s := e.request.URL.Query().Get(key)

	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("http.exchange.Exchanger.QueryParamUUID - cann't parse query param [%s]: %w", key, err)
	}

	return u, nil
}

func (e *Exchanger) hasQueryParam(key string) bool {
	return e.request.URL.Query().Has(key)
}
