package main

import (
    "log"
    "net/http"
)

func main() {
    const port = "8080"

    // https://pkg.go.dev/net/http#ServeMux
    mux := http.NewServeMux()

    // https://pkg.go.dev/net/http#Server
    server := http.Server {
        Addr: ":" + port,
        Handler: mux,
    }

    // https://pkg.go.dev/net/http#Server.ListenAndServe
    log.Printf("Serving on port: %s\n", port)
    log.Fatal(server.ListenAndServe())
}
