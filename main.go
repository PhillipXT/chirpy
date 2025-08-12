package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"
    "sync/atomic"

    "github.com/joho/godotenv"
    "github.com/PhillipXT/chirpy/internal/database"

    _ "github.com/lib/pq"
)

func main() {
    const port = "8080"
    const root = "./www"

    godotenv.Load()

    dbURL := os.Getenv("DB_URL")
    log.Printf("Using %s\n", dbURL)

    platform := os.Getenv("PLATFORM")
    if platform == "" {
        log.Fatal("PLATFORM must be set")
    }

    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Println("Error opening database, exiting...")
    }

    dbQueries := database.New(db)
    log.Printf("Show this: %v\n", dbQueries)

    cfg := Config {
        fileserverHits: atomic.Int32{},
        db: dbQueries,
        platform: platform,
    }

    // https://pkg.go.dev/net/http#FileServer
    fs := http.FileServer(http.Dir(root))

    // https://pkg.go.dev/net/http#ServeMux
    mux := http.NewServeMux()

    // Web Endpoints =======================================================>

    // https://pkg.go.dev/net/http#ServeMux.Handle
    handler := http.StripPrefix("/app", fs)
    mux.Handle("/app/", cfg.mwIncrementCounter(handler))

    // Api Endpoints =======================================================>

    mux.HandleFunc("GET /api/healthz", checkHealth)

    mux.HandleFunc("POST /api/users", cfg.createUser)
    mux.HandleFunc("POST /api/validate_chirp", validateChirp)

    // Admin Endpoints =====================================================>

    mux.HandleFunc("GET /admin/metrics", cfg.checkMetrics)

    // =====================================================================>

    // curl -X POST http://localhost:8080/api/reset
    mux.HandleFunc("POST /admin/reset", cfg.resetHitCounter)

    // https://pkg.go.dev/net/http#Server
    // http.Server is a struct that defines the server configuration
    server := http.Server {
        Addr: ":" + port,
        Handler: mux,
    }

    // https://pkg.go.dev/net/http#Server.ListenAndServe
    log.Printf("Serving files from %s on port: %s\n", root, port)
    log.Fatal(server.ListenAndServe())
}

func checkHealth(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(http.StatusText(http.StatusOK)))
}

