package muxhandler

import (
	"encoding/json"
	"net/http"
)

type jsonRes map[string]interface{}

func jsonResponse(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(obj)
}
