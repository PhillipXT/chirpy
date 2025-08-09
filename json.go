package main

import (
    "encoding/json"
    "net/http"
    "log"
)

type errorResponse struct {
    Error string
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {

    if err != nil {
        log.Println(err)
    }

    response := errorResponse {
        Error: message,
    }

    writeResponse(w, statusCode, response)
}

func writeResponse(w http.ResponseWriter, statusCode int, payload interface{}) {

    w.Header().Set("Content-Type", "application/json")

    data, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error in JSON: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(statusCode)
    w.Write(data)
}

