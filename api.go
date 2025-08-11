package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
)

const maxChirpLength = 140

type Chirp struct {
    Body string `json:"body"`
}

type okResponse struct {
    CleanedBody string `json:"cleaned_body"`
}

// Test with:
// curl -d '{"body":"01234567890"}' -X POST http://localhost:8080/api/validate_chirp
func validateChirp(w http.ResponseWriter, r *http.Request) {
    chirp := Chirp {}

    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&chirp)
    if err != nil {
        writeErrorResponse(w, http.StatusInternalServerError, "Error decoding JSON", err)
        return
    }

    log.Printf("Chirp length: %d\n", len(chirp.Body))
    if len(chirp.Body) > maxChirpLength {
        writeErrorResponse(w, http.StatusBadRequest, "Chirp is too long", nil)
        return
    }

    cleaned := filterChirp(chirp.Body)
    log.Printf("Filtered chirp: %s\n", cleaned)

    response := okResponse {
        CleanedBody: cleaned,
    }

    writeResponse(w, http.StatusOK, response)
}

func filterChirp(message string) string {
    badWords := map[string]struct{} {
        "kerfuffle": {},
        "sharbert": {},
        "fornax": {},
    }

    words := strings.Split(message, " ")

    for i, word := range words {
        if _, ok := badWords[strings.ToLower(word)]; ok {
            words[i] = "****"
        }
    }

    return strings.Join(words, " ")
}
