package main

import (
    "fmt"
    "log"
    "net/http"
    "sync/atomic"

    "github.com/PhillipXT/chirpy/internal/database"
)

type Config struct {
    fileserverHits atomic.Int32
    db *database.Queries
    platform string
}

func (cfg *Config) mwIncrementCounter(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    })
}

func (cfg *Config) checkMetrics(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)

    response := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())

    w.Write([]byte(response))
}

func (cfg *Config) resetHitCounter(w http.ResponseWriter, r *http.Request) {

    log.Printf("Platform: %s\n", cfg.platform)

    if cfg.platform != "dev" {
        writeErrorResponse(w, http.StatusForbidden, "Action not allowed", nil)
        return
    }

    err := cfg.db.Reset(r.Context())
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error deleting users", nil)
        return
    }

    cfg.fileserverHits.Store(0)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hits reset to 0 and database reset to initial state"))
}
