package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// These follow /v1
func (cfg *config) routesV1() http.Handler {
	router := chi.NewRouter()

	router.Post("/users/create", cfg.handlerUserCreate)
	router.Post("/users/login", cfg.handlerLogin)
	router.Post("/users/delete", cfg.handlerDeleteAllUsers)
	router.With(cfg.AuthMiddleware).Post("/users/logout", cfg.handlerLogout)
	router.With(cfg.AuthMiddleware).Get("/users/me", cfg.handlerGetMe)

	return router
}

// These follow /api/steam
func (cfg *config) routesAPI() http.Handler {
	router := chi.NewRouter()

	router.Get("/player-summaries", cfg.steamAPI.HandlerGetPlayerSummaries)
	router.Get("/friends", cfg.steamAPI.HandlerGetFriendList)
	router.Get("/games", cfg.steamAPI.HandlerGetOwnedGames)
	router.Get("/achievements", cfg.steamAPI.HandlerGetPlayerAchievements)
	router.Get("/friends/matchGames", cfg.steamAPI.HandlerMatchedGamesRanking)
	router.Get("/compare-achievements", cfg.steamAPI.HandlerCompareAchievements)

	return router
}
