package api

import (
	"net/http"
	"sync/atomic"

	"github.com/mjossany/Chirpy/internal/database"
)

type Config struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
	Port           string
	TokenSecret    string
}

func NewRouter(cfg *Config, filepathRoot string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("/api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/login", cfg.handleLogin)
	mux.HandleFunc("POST /api/users", cfg.handleCreateUser)
	mux.HandleFunc("GET /api/chirps", cfg.handleGetAllChirps)
	mux.HandleFunc("POST /api/chirps", cfg.handleCreateChirp)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handleGetChirpById)
	mux.HandleFunc("/admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("/admin/reset", cfg.handlerReset)
	return mux
}
