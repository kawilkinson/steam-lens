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

func (cfg *config) routesAPI() http.Handler {
	router := chi.NewRouter()

	router.Get("/player-summaries", cfg.steamAPI.HandlerGetPlayerSummaries)
	router.Get("/friends", cfg.steamAPI.HandlerGetFriendList)
	router.Get("/games", cfg.steamAPI.HandlerGetOwnedGames)
	router.Get("/achievements", cfg.steamAPI.HandlerGetPlayerAchievements)
	router.Get("/friends/matchGames", cfg.steamAPI.HandlerCompareOwnedGames)

	return router
}
