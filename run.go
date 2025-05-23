package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/Khazz0r/steam-lens/internal/api"
	"github.com/Khazz0r/steam-lens/internal/database"
	_ "github.com/lib/pq"
)

type config struct {
	db        *database.Queries
	platform  string
	jwtSecret string
	steamAPI  *api.ApiConfig
}

//go:embed static/*
var staticFiles embed.FS

func run() error {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("warning: assuming default configuration. .env unreadable: %v\n", err)
	}

	platform := getEnvOrFail("PLATFORM")
	dbURL := getEnvOrFail("DATABASE_URL")
	port := getEnvOrFail("PORT")
	jwtSecret := getEnvOrFail("JWTSECRET")
	steamAPIKey := getEnvOrFail("STEAM_API_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}

	cfg := &config{
		db:        database.New(db),
		platform:  platform,
		jwtSecret: jwtSecret,
		steamAPI: &api.ApiConfig{
			SteamApiKey: steamAPIKey,
			PlayerCache: api.Cache[api.Player]{
				Cache:     map[string]api.CachedData[api.Player]{},
				RenewTime: 24 * time.Hour,
			},
			FriendListCache: api.Cache[api.FriendList]{
				Cache:     map[string]api.CachedData[api.FriendList]{},
				RenewTime: 10 * time.Minute,
			},
			OwnedGamesCache: api.Cache[api.OwnedGames]{
				Cache:     map[string]api.CachedData[api.OwnedGames]{},
				RenewTime: 60 * time.Minute,
			},
		},
	}

	friendListCleaner := api.Cleaner[api.FriendList]{
		Name:     "FriendListCache",
		Cache:    &cfg.steamAPI.FriendListCache,
		Interval: 1 * time.Hour,
	}
	friendListCleaner.CacheCleanerStart()

	playerCleaner := api.Cleaner[api.Player]{
		Name:     "PlayerCache",
		Cache:    &cfg.steamAPI.PlayerCache,
		Interval: 1 * time.Hour,
	}
	playerCleaner.CacheCleanerStart()

	ownedGamesCleaner := api.Cleaner[api.OwnedGames]{
		Name:     "OwnedGamesCache",
		Cache:    &cfg.steamAPI.OwnedGamesCache,
		Interval: 1 * time.Hour,
	}
	ownedGamesCleaner.CacheCleanerStart()

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	router.Get("/", serveIndex)
	router.Mount("/v1", cfg.routesV1())
	router.Mount("/api/steam", cfg.routesAPI())

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 6 * time.Second,
	}

	fmt.Printf("starting server on port %s\n", port)
	return server.ListenAndServe()
}

func getEnvOrFail(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("environment variable %s must be set", key))
	}
	return val
}

func serveIndex(w http.ResponseWriter, req *http.Request) {
	file, err := staticFiles.Open("static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
