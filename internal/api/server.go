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
}

func NewRouter(cfg *Config, filepathRoot string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("/api/healthz", handlerReadiness)
	mux.HandleFunc("/api/users", cfg.handleCreateUser)
	mux.HandleFunc("/api/chirps", cfg.handleCreateChirp)
	mux.HandleFunc("/admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("/admin/reset", cfg.handlerReset)
	return mux
}
