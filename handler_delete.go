package main

import (
	"net/http"

	"github.com/Khazz0r/steam-lens/internal/api"
)

// handler to be used in a dev environment to delete all users in the database for testing purposes
func (cfg *config) handlerDeleteAllUsers(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("Reset is only allowed in dev environment"))
		if err != nil {
			api.RespondWithError(w, http.StatusInternalServerError, "Unable to write message to webpage", err)
			return
		}
		return
	}

	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Unable to delete all users", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Users table reset to initial state"))
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Unable to write message to webpage", err)
		return
	}
}
