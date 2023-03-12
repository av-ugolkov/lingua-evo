package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"lingua-evo/internal/config"
	"lingua-evo/pkg/logging"
	linguaJWT "lingua-evo/pkg/middleware/jwt"

	"github.com/cristalhq/jwt/v3"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

const (
	authURL   = "/auth"
	signupURL = "/api/signup"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type newUser struct {
	user
	Email string `json:"email*"`
}

type refresh struct {
	RefreshToken string `json:"refresh_token"`
}

type Handler struct {
	logger *logging.Logger
	//RTCache cache.Repository
}

func NewHandler(logger *logging.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, authURL, h.authGet)
	router.HandlerFunc(http.MethodPut, authURL, h.authPut)
	router.HandlerFunc(http.MethodPost, authURL, h.authPost)
	router.HandlerFunc(http.MethodPost, signupURL, h.signup)
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
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
