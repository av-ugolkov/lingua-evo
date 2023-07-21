package sign_up

import (
	"encoding/json"
	"net/http"
	"time"

	"lingua-evo/internal/config"
	"lingua-evo/internal/service"
	"lingua-evo/pkg/logging"
	linguaJWT "lingua-evo/pkg/middleware/jwt"

	"github.com/cristalhq/jwt/v3"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

const (
	signUpURL = "/signup"
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
	router.HandlerFunc(http.MethodGet, signUpURL, h.getSignUp)
	router.HandlerFunc(http.MethodPost, signUpURL, h.postSignUp)
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
