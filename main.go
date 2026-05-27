package main

import (
	"database/sql"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/PhanNam1501/chirpy/internal/database"
	"github.com/PhanNam1501/chirpy/internal/handler"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func middlewareMetricsInc(cfg *handler.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET not set")
	}
	apiKey := os.Getenv("POLKA_KEY")
	if apiKey == "" {
		panic("API KEY not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		panic(err)
	}

	dbQueries := database.New(db)
	cfg := &handler.Config{
		FileserverHits: atomic.Int32{},
		DB:             dbQueries,
		Platform:       platform,
		JWTSecret:      jwtSecret,
		ApiKey:         apiKey,
	}

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("GET /api/healthz", handler.Healthz)
	serveMux.HandleFunc("GET /admin/metrics", cfg.MetricsHandler)
	serveMux.HandleFunc("GET /api/chirps", cfg.GetAllChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", cfg.GetChirpByID)
	serveMux.HandleFunc("POST /api/validate_chirp", cfg.ValidateChirp)
	serveMux.HandleFunc("POST /admin/reset", cfg.Reset)
	serveMux.HandleFunc("POST /api/users", cfg.AddUser)
	serveMux.HandleFunc("POST /api/chirps", cfg.CreateChirp)
	serveMux.HandleFunc("POST /api/login", cfg.Login)
	serveMux.HandleFunc("POST /api/refresh", cfg.Refresh)
	serveMux.HandleFunc("POST /api/revoke", cfg.Revoke)
	serveMux.HandleFunc("POST /api/polka/webhooks", cfg.PolkaWebhook)
	serveMux.HandleFunc("PUT /api/users", cfg.UpdateUser)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.DeleteChirpByID)

	fileServerHandler := http.FileServer(http.Dir("."))
	serveMux.Handle("/app/", middlewareMetricsInc(cfg, http.StripPrefix("/app", fileServerHandler)))

	server := &http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	server.ListenAndServe()
}
