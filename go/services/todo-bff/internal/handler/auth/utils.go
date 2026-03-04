package auth

import (
	"fmt"
	"net/http"
)

func renderError(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	msg := fmt.Sprintf(`{"message":"%s"}`, err.Error())
	// nolint: errcheck
	w.Write([]byte(msg))
}
