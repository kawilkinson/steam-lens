package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
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

	// get platform to determine if dev or prod, if dev make devPlatform false for secure cookies
	devPlatform := os.Getenv("PLATFORM") != "dev"

	// Set HttpOnly cookies for both tokens
	http.SetCookie(w, &http.Cookie{
		Name:     "JWT_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   devPlatform,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Hour),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   devPlatform,
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

func (cfg *config) handlerLogout(w http.ResponseWriter, req *http.Request) {
	// get platform to determine if dev or prod, if dev make devPlatform false for secure cookies
	devPlatform := os.Getenv("PLATFORM") != "dev"

	refreshCookie, err := req.Cookie("refresh_token")
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, "Unable to retrieve refresh_token from request", err)
		return
	}
	if refreshCookie.Value == "" {
		api.RespondWithError(w, http.StatusBadRequest, "refresh_token provided is not valid, please try logging out again", err)
		return
	}

	err = cfg.db.DeleteRefreshToken(req.Context(), refreshCookie.Value)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Unable to delete refresh token from database", err)
		return
	}

	// Invalidate both JWT token and refresh token cookies for logging out
	http.SetCookie(w, &http.Cookie{
		Name:     "JWT_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   devPlatform,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(-1 * time.Hour),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   devPlatform,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(-1 * time.Hour),
	})

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "Successfully logged out"}`))
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Unable to write out logout message", err)
		return
	}
}

func (cfg *config) handlerGetMe(w http.ResponseWriter, req *http.Request) {
	type response struct {
		User `json:"user"`
	}

	userID, exists := req.Context().Value(userIDContextKey).(uuid.UUID)
	if !exists || userID == uuid.Nil {
		api.RespondWithError(w, http.StatusUnauthorized, "Not authorized to perform this action", nil)
		return
	}

	user, err := cfg.db.GetUserByID(req.Context(), userID)
	if err != nil {
		api.RespondWithError(w, http.StatusNotFound, "Could not find user by that ID", err)
		return
	}

	api.RespondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:       user.ID,
			Username: user.Username,
			SteamID:  user.SteamID,
		},
	})
}

type updateUserParams struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	SteamID  *string `json:"steam_id"`
}

// Helper function to safely dereference a string pointer that could be empty
func deref(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

func (cfg *config) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	userID, exists := req.Context().Value(userIDContextKey).(uuid.UUID)
	if !exists || userID == uuid.Nil {
		api.RespondWithError(w, http.StatusUnauthorized, "You are not authorized to update this account", nil)
		return
	}

	decoder := json.NewDecoder(req.Body)

	params := updateUserParams{}
	err := decoder.Decode(&params)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Unable to decode request body", err)
		return
	}

	var usernamePtr *string
	var hashedPasswordPtr *string
	var steamIDPtr *string

	if params.Username != nil && strings.TrimSpace(*params.Username) != "" {
		usernamePtr = params.Username
	}

	if params.Password != nil && strings.TrimSpace(*params.Password) != "" {
		hash, err := auth.HashPassword(*params.Password)
		if err != nil {
			api.RespondWithError(w, http.StatusBadRequest, "Unable to hash password", err)
			return
		}
		hashedPasswordPtr = &hash
	}

	if params.SteamID != nil && strings.TrimSpace(*params.SteamID) != "" {
		steamIDPtr = params.SteamID
	}

	err = cfg.db.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:        userID,
		Column1:   deref(usernamePtr),
		Column2:   deref(hashedPasswordPtr),
		Column3:   deref(steamIDPtr),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Failed to update user with provided fields, duplicate username found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message":"updated"}`))
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, "Failed to write reply", err)
		return
	}
}
