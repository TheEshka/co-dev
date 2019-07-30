package common

import (
	"encoding/json"
	errors2 "github.com/misgorod/co-dev/errors"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func RespondError(w http.ResponseWriter, err error) {
	m, c := errors2.Resolve(err)
	RespondJSON(w, c, map[string]string{"error": m})
}
