package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	errEmptyBody = errors.New("body is empty")
)

func CheckBody(w http.ResponseWriter, r *http.Request, body any) error {
	if r.Body == http.NoBody {
		return fmt.Errorf("tools.CheckBody: %w", errEmptyBody)
	}

	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		return fmt.Errorf("tools.CheckBody - decode: %w", err)
	}

	return nil
}
