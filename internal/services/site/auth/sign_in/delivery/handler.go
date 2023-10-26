package delivery

import (
	"fmt"
	"log/slog"
	"net/http"

	"lingua-evo/internal/services/user/service"

	staticFiles "lingua-evo"

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

	user, err := h.userSvc.GetUserByName(r.Context(), username)
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
}
