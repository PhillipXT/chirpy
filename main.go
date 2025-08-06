package main

import (
    "log"
    "net/http"
)

func main() {
    const port = "8080"
    const root = "./www"

    cfg := Config {}
    cfg.fileserverHits.Store(0)

    // https://pkg.go.dev/net/http#FileServer
    fs := http.FileServer(http.Dir(root))

    // https://pkg.go.dev/net/http#ServeMux
    mux := http.NewServeMux()

    // https://pkg.go.dev/net/http#ServeMux.Handle
    //mux.Handle("/", fs)
    handler := http.StripPrefix("/app", fs)
    mux.Handle("/app/", cfg.mwIncrementCounter(handler))
    mux.HandleFunc("/healthz", checkHealth)
    mux.HandleFunc("/metrics", cfg.checkMetrics)
    mux.HandleFunc("/reset", cfg.resetHitCounter)

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

