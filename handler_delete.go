package main

import "net/http"

// handler to be used in a dev environment to delete all users in the database for testing purposes
func (cfg *apiConfig) handlerDeleteAllUsers(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("Reset is only allowed in dev environment"))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to write message to webpage", err)
		}
		return
	}

	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {

	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Users table reset to initial state"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to write message to webpage", err)
	}
}
