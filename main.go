package main

import (
    "log"
    "net/http"
)

func main() {
    const port = "8080"
    const root = "./www"

    // https://pkg.go.dev/net/http#ServeMux
    mux := http.NewServeMux()

    // https://pkg.go.dev/net/http#Server
    server := http.Server {
        Addr: ":" + port,
        Handler: mux,
    }

    // https://pkg.go.dev/net/http#FileServer
    fs := http.FileServer(http.Dir(root))

    // https://pkg.go.dev/net/http#ServeMux.Handle
    mux.Handle("/", fs)

    // https://pkg.go.dev/net/http#Server.ListenAndServe
    log.Printf("Serving files from %s on port: %s\n", root, port)
    log.Fatal(server.ListenAndServe())
}
