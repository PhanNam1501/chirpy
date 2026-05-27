package handler

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/PhanNam1501/chirpy/internal/database"
)

type Config struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
	JWTSecret      string
	ApiKey         string
}

func (cfg *Config) Reset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.FileserverHits.Store(0)
	if cfg.Platform == "dev" {
		cfg.DB.DeleteAllUsers(r.Context())
	}
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *Config) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.FileserverHits.Load()
	html := `<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited ` + fmt.Sprintf("%d", hits) + ` times!</p>
		</body>
		</html>`
	w.Write([]byte(html))
}
