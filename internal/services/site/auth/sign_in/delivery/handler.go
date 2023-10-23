package delivery

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"lingua-evo/internal/config"
	"lingua-evo/internal/services/lingua/user/service"
	"lingua-evo/internal/services/site/auth/sign_in/entity"

	staticFiles "lingua-evo"
	linguaJWT "lingua-evo/pkg/middleware/jwt"

	"github.com/cristalhq/jwt/v3"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	signInURL = "/signin"

	signInPage = "website/sign_in/signin.html"
)

type (
	Handler struct {
		userSvc *service.UserSvc
	}
)

func Create(r *mux.Router, userSvc *service.UserSvc) {
	handler := newHandler(userSvc)
	handler.register(r)
}

func newHandler(userSvc *service.UserSvc) *Handler {
	return &Handler{
		userSvc: userSvc,
	}
}

func (h *Handler) register(r *mux.Router) {
	r.HandleFunc(signInURL, h.get).Methods(http.MethodGet)
	r.HandleFunc(signInURL, h.put).Methods(http.MethodPut)
	r.HandleFunc(signInURL, h.post).Methods(http.MethodPost)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	file, err := staticFiles.OpenFile(signInPage)
	if err != nil {
		slog.Error(fmt.Errorf("sign_in.get.OpenFile: %v", err).Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(file))
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	//email := r.FormValue("email")
	//password := r.FormValue("password")

	user, err := h.userSvc.GetIDByName(r.Context(), username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	fmt.Printf("find user: %v", user)
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
	_, _ = w.Write(request)
}

func (h *Handler) put(w http.ResponseWriter, r *http.Request) {
	var refreshToken entity.Refresh
	if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
		slog.Error(err.Error())
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
	slog.Info(string(jsonBytes))
	_, _ = w.Write(jsonBytes)
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
		slog.Error(err.Error())
		return nil, http.StatusUnauthorized
	}

	slog.Info("create refresh token")
	refreshTokenUuid := uuid.New()
	/*err = h.RTCache.Set([]byte(refreshTokenUuid.String()), []byte(claims.ID), 0)
	if err != nil {
		slog.Error(err)
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
