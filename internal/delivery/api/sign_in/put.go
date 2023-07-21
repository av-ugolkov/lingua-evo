package sign_in

import (
	"encoding/json"
	"net/http"
)

type Refresh struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) putSignIn(w http.ResponseWriter, r *http.Request) {
	var refreshToken Refresh
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
