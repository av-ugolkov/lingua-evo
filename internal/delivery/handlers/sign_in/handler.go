package sign_in

import (
	"encoding/json"
	"net/http"
	"time"

	"lingua-evo/internal/config"
	"lingua-evo/internal/delivery/handlers/sign_in/entity"
	"lingua-evo/internal/service"
	staticFiles "lingua-evo/static"

	"lingua-evo/pkg/logging"
	linguaJWT "lingua-evo/pkg/middleware/jwt"

	"github.com/cristalhq/jwt/v3"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

const (
	signInURL = "/signin"

	signInPage = "web/sign_in/signin.html"
)

type Handler struct {
	logger *logging.Logger
	lingua *service.Lingua
	//RTCache cache.Repository
}

func Create(log *logging.Logger, ling *service.Lingua, r *httprouter.Router) {
	handler := newHandler(log, ling)
	handler.register(r)
}

func newHandler(logger *logging.Logger, lingua *service.Lingua) *Handler {
	return &Handler{
		logger: logger,
		lingua: lingua,
	}
}

func (h *Handler) register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, signInURL, h.get)
	router.HandlerFunc(http.MethodPut, signInURL, h.put)
	router.HandlerFunc(http.MethodPost, signInURL, h.post)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	file, err := staticFiles.OpenFile(signInPage)
	if err != nil {
		h.logger.Errorf("sign_in.get.OpenFile: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(file))
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	var u entity.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//TODO client to UserService and get user by username and password
	//for now stub check
	//if u.Username != "me" || u.Password != "pass" {
	//	w.WriteHeader(http.StatusNotFound)
	//	return
	//}

	jsonBytes, errCode := h.generateAccessToken()
	if errCode != 0 {
		w.WriteHeader(errCode)
		return
	}
	request, err := json.Marshal(map[string]string{
		"token": string(jsonBytes),
		"url":   "/account",
	})
	if err != nil {
		w.WriteHeader(errCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(request)
}

func (h *Handler) put(w http.ResponseWriter, r *http.Request) {
	var refreshToken entity.Refresh
	if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
		h.logger.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	/*userIdBytes, err := h.RTCache.Get([]byte(refreshTokenS.RefreshToken))
	h.Logger.Infof("refresh token user_id: %s", userIdBytes)
	if err != nil {
		h.Logger.Fatal(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	h.RTCache.Del([]byte(refreshTokenS.RefreshToken))*/
	//TODO client to UserService and get user by username

	jsonBytes, errCode := h.generateAccessToken()
	if errCode != 0 {
		w.WriteHeader(errCode)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	h.logger.Info(string(jsonBytes))
	w.Write(jsonBytes)
}

func (h *Handler) generateAccessToken() ([]byte, int) {
	key := []byte(config.GetConfig().JWT.Secret)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return nil, 418
	}
	builder := jwt.NewBuilder(signer)

	//TODO insert real user data in claims
	claims := linguaJWT.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "uuid_here",
			Audience:  []string{"users"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
		Email: "email@will.be.here",
	}
	token, err := builder.Build(claims)
	if err != nil {
		h.logger.Error(err)
		return nil, http.StatusUnauthorized
	}

	h.logger.Info("create refresh token")
	refreshTokenUuid := uuid.New()
	/*err = h.RTCache.Set([]byte(refreshTokenUuid.String()), []byte(claims.ID), 0)
	if err != nil {
		h.Logger.Error(err)
		return nil, http.StatusInternalServerError
	}*/
	jsonBytes, err := json.Marshal(map[string]string{
		"token":         token.String(),
		"refresh_token": refreshTokenUuid.String(),
	})
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return jsonBytes, 0
}
