package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "JSON encoding error: %v"}`, err), http.StatusInternalServerError)
	}
}

func WriteJSONRedis(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func WriteError(w http.ResponseWriter, status int, err string) {
	w.WriteHeader(status)
	log.Print(err)
	if err := json.NewEncoder(w).Encode(map[string]string{"err": err}); err != nil {
		http.Error(w, `{"err": err encode error}`, status)
	}
}
