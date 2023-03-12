package auth

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) authPut(w http.ResponseWriter, r *http.Request) {
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
