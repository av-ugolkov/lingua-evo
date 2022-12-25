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

type handler struct {
	logger *logging.Logger
	//RTCache cache.Repository
}

func NewHandler(logger *logging.Logger) *handler {
	return &handler{
		logger: logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, authURL, h.Auth)
	router.HandlerFunc(http.MethodPut, authURL, h.Auth)
	router.HandlerFunc(http.MethodPost, signupURL, h.Signup)
}

func (h *handler) Signup(w http.ResponseWriter, r *http.Request) {
	var nu newUser
	if err := json.NewDecoder(r.Body).Decode(&nu); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	//TODO validate username and password
	//TODO create user using UserService
	jsonBytes, errCode := h.generateAccessToken()
	if errCode != 0 {
		w.WriteHeader(errCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(jsonBytes)
}

func (h *handler) Auth(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var u user
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
	case http.MethodPut:
		var refreshTokenS refresh
		if err := json.NewDecoder(r.Body).Decode(&refreshTokenS); err != nil {
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
	}

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

func (h *handler) generateAccessToken() ([]byte, int) {
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
