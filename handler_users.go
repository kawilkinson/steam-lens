package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Khazz0r/steam-lens/internal/api"
	"github.com/Khazz0r/steam-lens/internal/auth"
	"github.com/Khazz0r/steam-lens/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	SteamID   string    `json:"steam_id"`
}

func (cfg *config) handlerUserCreate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		SteamID  string `json:"steam_id"`
		Password string `json:"password"`
	}
	type response struct {
		User `json:"user"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	err = cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Username:       params.Username,
		HashedPassword: hashedPassword,
		SteamID:        params.SteamID,
	})
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "User already exists", err)
		return
	}

	user, err := cfg.db.GetUserByUsername(req.Context(), params.Username)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Couldn't get user", err)
		return
	}

	api.RespondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Username:  user.Username,
			SteamID:   user.SteamID,
		},
	})
}

func (cfg *config) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		User         `json:"user"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	user, err := cfg.db.GetUserByUsername(req.Context(), params.Username)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "The username provided does not exist", err)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		api.RespondWithError(w, http.StatusUnauthorized, "Cannot login, wrong password provided", err)
		return
	}

	accessToken, err := auth.MakeJWTToken(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Could not make JWT token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Could not make refresh token", err)
		return
	}

	err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add((15 * 24) * time.Hour),
	})
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Could not set refresh token into database", err)
		return
	}

	// Set HttpOnly cookies for both tokens
	http.SetCookie(w, &http.Cookie{
		Name:     "JWT_token", // JWT access token
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // false for now since I'm working in just a dev environment
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Hour),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token", // refresh token
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // false for now since I'm working in just a dev environment
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(15 * 24 * time.Hour),
	})

	api.RespondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Username:  user.Username,
			SteamID:   user.SteamID,
		},
		Token:        accessToken,  // remove this when ran in a prod environment
		RefreshToken: refreshToken, // remove too
	})
}
