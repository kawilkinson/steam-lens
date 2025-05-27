package main

import (
	"context"
	"net/http"

	"github.com/Khazz0r/steam-lens/internal/api"
	"github.com/Khazz0r/steam-lens/internal/auth"
	"github.com/google/uuid"
)

type contextKey string

const userIDContextKey = contextKey("userID")

func (cfg *config) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("JWT_token")
		if err != nil {
			api.RespondWithError(w, http.StatusUnauthorized, "Not authorized to do this, missing the proper refresh_token", err)
			return
		}

		userID, err := auth.ValidateJWT(cookie.Value, cfg.jwtSecret)
		if err != nil || userID == uuid.Nil {
			api.RespondWithError(w, http.StatusUnauthorized, "Not authorized to do this, invalid token provided", err)
			return
		}

		context := context.WithValue(req.Context(), userIDContextKey, userID)
		next.ServeHTTP(w, req.WithContext(context))
	})
}
