package transport

import (
	"encoding/json"
	db "gophermart/internal/database"
	"gophermart/internal/services"
	"net/http"
)

const (
	ApplicationJSON = "application/json"
)

type handlersData struct {
	storage db.StoragerDB
}

func New(storage db.StoragerDB) *handlersData {
	return &handlersData{storage: storage}
}

func (h *handlersData) Registration(w http.ResponseWriter, r *http.Request) {

	var user services.UserAuthInfo
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.storage.GetUser(user)

	setResponseHeaders(w, ApplicationJSON, http.StatusOK)
}

func setResponseHeaders(w http.ResponseWriter, contentType string, statusCode int) {

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

}
