package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (cfg *config) routesV1() http.Handler {
	router := chi.NewRouter()

	router.Post("/users/create", cfg.handlerUserCreate)
	router.Post("/users/login", cfg.handlerLogin)
	router.Post("/users/delete", cfg.handlerDeleteAllUsers)

	return router
}
