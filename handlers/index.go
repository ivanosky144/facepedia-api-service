package handlers

import (
	"encoding/json"
	"net/http"
)

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error formatting JSON response", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonResponse)
}
