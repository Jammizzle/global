package muxhandler

import (
	"encoding/json"
	"net/http"
)

type JsonRes map[string]interface{}

func JsonResponse(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(obj)
}
